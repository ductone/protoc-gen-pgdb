// Spawns a PostgreSQL server with a single database configured. Ideal for unit
// tests where you want a clean instance each time. Then clean up afterwards.
//
// Requires PostgreSQL to be installed on your system (but it doesn't have to be running).
package pgtest

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PG struct {
	dir string
	cmd *exec.Cmd
	DB  *pgxpool.Pool

	persistent bool

	stderr io.ReadCloser
	stdout io.ReadCloser
}

// Start a new PostgreSQL database, on temporary storage.
//
// This database has fsync disabled for performance, so it might run faster
// than your production database. This makes it less reliable in case of system
// crashes, but we don't care about that anyway during unit testing.
//
// Use the DB field to access the database connection
func Start() (*PG, error) {
	return start(New())
}

// Starts a new PostgreSQL database
//
// Will listen on a unix socket and initialize the database in the given
// folder, if needed. Data isn't removed when calling Stop(), so this database
// can be used multiple times. Allows using PostgreSQL as an embedded databases
// (such as SQLite). Not for production usage!
func StartPersistent(folder string) (*PG, error) {
	return start(New().DataDir(folder).Persistent())
}

// start Starts a new PostgreSQL database
//
// Will listen on a unix socket and initialize the database in the given
// folder (config.Dir), if needed.
// Data isn't removed when calling Stop() if config.Persistent == true,
// so this database
// can be used multiple times. Allows using PostgreSQL as an embedded databases
// (such as SQLite). Not for production usage!
func start(config *PGConfig) (*PG, error) {
	backgroundContext := context.Background()
	ctx, done := context.WithTimeout(backgroundContext, time.Second*15)
	defer done()

	// Handle dropping permissions when running as root
	me, err := user.Current()
	if err != nil {
		return nil, err
	}
	isRoot := me.Username == "root"

	pgUID := int(0)
	pgGID := int(0)
	if isRoot {
		pgUser, err := user.Lookup("postgres")
		if err != nil {
			return nil, fmt.Errorf("could not find postgres user, which is required when running as root: %s", err)
		}

		uid, err := strconv.ParseInt(pgUser.Uid, 10, 64)
		if err != nil {
			return nil, err
		}
		pgUID = int(uid)

		gid, err := strconv.ParseInt(pgUser.Gid, 10, 64)
		if err != nil {
			return nil, err
		}
		pgGID = int(gid)
	}

	// Prepare data directory
	dir := config.Dir
	if config.Dir == "" {
		d, err := os.MkdirTemp("", "pgtest")
		if err != nil {
			return nil, err
		}
		dir = d
	}

	dataDir := filepath.Join(dir, "data")
	sockDir := filepath.Join(dir, "sock")

	err = os.MkdirAll(dataDir, 0711)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(sockDir, 0711)
	if err != nil {
		return nil, err
	}

	if isRoot {
		err = os.Chmod(dir, 0711)
		if err != nil {
			return nil, err
		}

		err = os.Chown(dataDir, pgUID, pgGID)
		if err != nil {
			return nil, err
		}

		err = os.Chown(sockDir, pgUID, pgGID)
		if err != nil {
			return nil, err
		}
	}

	// Find executables root path
	binPath, err := findBinPath(config.BinDir)
	if err != nil {
		return nil, err
	}

	// Initialize PostgreSQL data directory
	_, err = os.Stat(filepath.Join(dataDir, "postgresql.conf"))
	if os.IsNotExist(err) {
		init := prepareCommand(ctx, isRoot, filepath.Join(binPath, "initdb"),
			"-D", dataDir,
			"--no-sync",
		)
		out, err := init.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize DB: %w -> %s", err, string(out))
		}
	}

	// Start PostgreSQL
	cmd := prepareCommand(backgroundContext, isRoot, filepath.Join(binPath, "postgres"),
		"-D", dataDir, // Data directory
		"-k", sockDir, // Location for the UNIX socket
		"-h", "", // Disable TCP listening
		"-F", // No fsync, just go fast
	)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stderr.Close()
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, abort("Failed to start PostgreSQL", cmd, stderr, stdout, err)
	}

	// Connect to DB
	dsn := makeDSN(sockDir, "postgres", isRoot)
	// Prepare test database
	err = retry(ctx, func() error {
		// when debugging, you might want to look at this loop!
		//		spew.Dump("attempting to connect ", dsn)
		db, err := pgxpool.Connect(ctx, dsn)
		if err != nil {
			return err
		}
		var exists bool
		err = db.QueryRow(ctx, "SELECT true FROM pg_database WHERE datname = 'test'").Scan(&exists)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		if exists {
			return nil
		}

		_, err = db.Exec(ctx, "CREATE DATABASE test")
		db.Close()
		return err
	}, 50, 50*time.Millisecond)
	if err != nil {
		return nil, abort("Failed to initialize DB", cmd, stderr, stdout, err)
	}

	// Connect to it properly
	dsn = makeDSN(sockDir, "test", isRoot)
	db, err := pgxpool.Connect(backgroundContext, dsn)
	if err != nil {
		return nil, abort("Failed to connect to test DB", cmd, stderr, stdout, err)
	}

	pg := &PG{
		cmd: cmd,
		dir: dir,

		DB: db,

		persistent: config.IsPersistent,

		stderr: stderr,
		stdout: stdout,
	}

	return pg, nil
}

// Stop the database and remove storage files.
func (p *PG) Stop() error {
	if p == nil {
		return nil
	}

	if !p.persistent {
		defer func() {
			// Always try to remove it
			os.RemoveAll(p.dir)
		}()
	}

	err := p.cmd.Process.Signal(os.Interrupt)
	if err != nil {
		return err
	}
	done := make(chan struct{})
	go func() {
		select {
		case <-time.After(time.Second * 2):
			_ = p.cmd.Process.Signal(os.Kill)
		case <-done:
			return
		}
	}()
	err = p.cmd.Wait()
	close(done)
	if err != nil {
		return err
	}

	if p.stderr != nil {
		p.stderr.Close()
	}

	if p.stdout != nil {
		p.stdout.Close()
	}

	return nil
}

// Needed because Ubuntu doesn't put initdb in $PATH
// binDir a path to a directory that contains postgresql binaries
func findBinPath(binDir string) (string, error) {
	// In $PATH (e.g. Fedora) great!
	if binDir == "" {
		p, err := exec.LookPath("initdb")
		if err == nil {
			return path.Dir(p), nil
		}
	}

	// Look for a PostgreSQL in one of the folders Ubuntu uses
	folders := []string{
		binDir,
		"/usr/lib/postgresql/",
	}
	for _, folder := range folders {
		f, err := os.Stat(folder)
		if os.IsNotExist(err) {
			continue
		}
		if !f.IsDir() {
			continue
		}

		files, err := os.ReadDir(folder)
		if err != nil {
			return "", err
		}
		for _, fi := range files {
			if !fi.IsDir() && fi.Name() == "initdb" {
				return filepath.Join(folder), nil
			}

			if !fi.IsDir() {
				continue
			}

			binPath := filepath.Join(folder, fi.Name(), "bin")
			_, err := os.Stat(filepath.Join(binPath, "initdb"))
			if err == nil {
				return binPath, nil
			}
		}
	}

	return "", fmt.Errorf("did not find PostgreSQL executables installed")
}

func makeDSN(sockDir, dbname string, isRoot bool) string {
	user := ""
	if isRoot {
		user = "user=postgres"
	}
	return fmt.Sprintf("host=%s dbname=%s %s", sockDir, dbname, user)
}

func retry(ctx context.Context, fn func() error, attempts int, interval time.Duration) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		err := fn()
		if err == nil {
			return nil
		}

		attempts -= 1
		if attempts <= 0 {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
		}
	}
}

func prepareCommand(ctx context.Context, isRoot bool, command string, args ...string) *exec.Cmd {
	if !isRoot {
		return exec.CommandContext(ctx, command, args...)
	}

	for i, a := range args {
		if a == "" {
			args[i] = "''"
		}
	}

	return exec.CommandContext(ctx, "su",
		"-",
		"postgres",
		"-c",
		strings.Join(append([]string{command}, args...), " "),
	)
}

func abort(msg string, cmd *exec.Cmd, stderr, stdout io.ReadCloser, err error) error {
	cmd.Process.Signal(os.Interrupt)
	cmd.Wait()

	serr, _ := io.ReadAll(stderr)
	sout, _ := io.ReadAll(stdout)
	stderr.Close()
	stdout.Close()
	return fmt.Errorf("%s: %s\nOUT: %s\nERR: %s", msg, err, string(sout), string(serr))
}