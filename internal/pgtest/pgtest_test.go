package pgtest

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPostgreSQL(t *testing.T) {
	ctx := context.Background()
	t.Parallel()

	assert := require.New(t)

	pg, err := Start()
	assert.NoError(err)
	assert.NotNil(pg)

	_, err = pg.DB.Exec(ctx, "CREATE TABLE test (val text)")
	assert.NoError(err)

	err = pg.Stop()
	assert.NoError(err)
}

func TestPostgreSQLWithConfig(t *testing.T) {
	ctx := context.Background()
	t.Parallel()

	assert := require.New(t)
	pg, err := New().From("/usr/bin/").Start()
	assert.NoError(err)
	assert.NotNil(pg)

	_, err = pg.DB.Exec(ctx, "CREATE TABLE test (val text)")
	assert.NoError(err)

	err = pg.Stop()
	assert.NoError(err)
}

func TestPersistent(t *testing.T) {
	ctx := context.Background()
	t.Parallel()

	assert := require.New(t)

	dir, err := os.MkdirTemp("", "pgtest")
	assert.NoError(err)
	defer os.RemoveAll(dir)

	pg, err := StartPersistent(dir)
	assert.NoError(err)
	assert.NotNil(pg)

	_, err = pg.DB.Exec(ctx, "CREATE TABLE test (val text)")
	assert.NoError(err)

	_, err = pg.DB.Exec(ctx, "INSERT INTO test VALUES ('foo')")
	assert.NoError(err)

	err = pg.Stop()
	assert.NoError(err)

	// Open it again
	pg, err = StartPersistent(dir)
	assert.NoError(err)
	assert.NotNil(pg)

	var val string
	err = pg.DB.QueryRow(ctx, "SELECT val FROM test").Scan(&val)
	assert.NoError(err)
	assert.Equal(val, "foo")

	err = pg.Stop()
	assert.NoError(err)
}
