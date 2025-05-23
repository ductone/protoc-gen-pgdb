package pgdb

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

func New() pgs.Module {
	return &Module{ModuleBase: &pgs.ModuleBase{}}
}

const (
	moduleName = "pgdb"
	version    = "0.1.0"
)

type Module struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
}

var _ pgs.Module = (*Module)(nil)

func (m *Module) InitContext(ctx pgs.BuildContext) {
	m.ModuleBase.InitContext(ctx)
	m.ctx = pgsgo.InitContext(ctx.Parameters())
}

func (m *Module) Name() string {
	return moduleName
}

func (m *Module) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	for _, f := range targets {
		msgs := f.AllMessages()
		if n := len(msgs); n == 0 {
			m.Debugf("No messagess in %v, skipping", f.Name())
			continue
		}
		m.processFile(m.ctx, f)
	}
	return m.Artifacts()
}

func (m *Module) processFile(ctx pgsgo.Context, f pgs.File) {
	out := bytes.Buffer{}
	err := m.applyTemplate(ctx, &out, f)
	if err != nil {
		m.Logf("couldn't apply template: %s", err)
		m.Fail("code generation failed")
	} else {
		generatedFileName := m.ctx.OutputPath(f).SetExt(fmt.Sprintf(".%s.go", moduleName)).String()
		if ok, _ := strconv.ParseBool(os.Getenv("PGDB_DEBUG_FILE_RAW")); ok {
			spew.Fdump(os.Stderr, out.String())
			_, _ = fmt.Fprintf(os.Stderr, "\n%s\n", out.String())
		}
		m.AddGeneratorFile(generatedFileName, out.String())
	}
}

func (module *Module) applyTemplate(ctx pgsgo.Context, outputBuffer *bytes.Buffer, in pgs.File) error {
	ix := &importTracker{
		ctx:        ctx,
		input:      in,
		typeMapper: make(map[pgs.Name]pgs.FilePath),
	}
	buf := &bytes.Buffer{}

	for _, m := range in.AllMessages() {
		fext := pgdb_v1.MessageOptions{}
		_, err := m.Extension(pgdb_v1.E_Msg, &fext)
		if err != nil {
			return fmt.Errorf("pgdb: applyTemplate: failed to extract Message extension from '%s': %w", m.FullyQualifiedName(), err)
		}

		if fext.GetDisabled() {
			continue
		}

		err = module.renderDescriptor(ctx, buf, in, m, ix)
		if err != nil {
			return err
		}

		err = module.renderMessage(ctx, buf, in, m, ix)
		if err != nil {
			return err
		}

		err = module.renderQueryBuilder(ctx, buf, in, m, ix)
		if err != nil {
			return err
		}
	}

	err := module.renderHeader(ctx, outputBuffer, in, ix)
	if err != nil {
		return err
	}
	_, err = io.Copy(outputBuffer, buf)
	if err != nil {
		return err
	}

	if ok, _ := strconv.ParseBool(os.Getenv("PGDB_DUMP_FILE")); ok {
		tdr := os.TempDir()
		_ = os.WriteFile(filepath.Join(tdr, "t.go"), outputBuffer.Bytes(), 0600)
		_, _ = os.Stderr.WriteString(filepath.Join(tdr, "t.go") + "\n")
	}
	return nil
}
