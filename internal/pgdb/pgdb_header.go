package pgdb

import (
	"fmt"
	"io"
	"sort"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type importTracker struct {
	PGDBV1 bool
	XPQ    bool

	Bytes                bool
	Fmt                  bool
	GoquExp              bool
	JSON                 bool
	ProtobufEncodingJSON bool
	ProtobufProto        bool
	Strings              bool
	Time                 bool
	PgType               bool

	ctx          pgsgo.Context
	input        pgs.File
	ProtoImports []string
}

func (ix *importTracker) AddProtoEntity(entity pgs.Entity) {
	// TOD(pquerna): is there a better method to do this?
	if ix.ctx.ImportPath(entity) == ix.ctx.ImportPath(ix.input) {
		// os.Stderr.WriteString(fmt.Sprintf("[skip] AddProtoEntity: %s == %s\n", ix.ctx.ImportPath(entity), ix.ctx.ImportPath(ix.input)))
		return
	}
	// os.Stderr.WriteString(fmt.Sprintf("[add] AddProtoEntity: %s != %s\n", ix.ctx.ImportPath(entity), ix.ctx.ImportPath(ix.input)))
	tmp := make(map[string]struct{})
	for _, v := range ix.ProtoImports {
		tmp[v] = struct{}{}
	}
	key := fmt.Sprintf(`%s "%s"`, ix.ctx.PackageName(entity), string(ix.ctx.ImportPath(entity)))
	tmp[key] = struct{}{}
	ix.ProtoImports = make([]string, 0, len(tmp))
	for k, _ := range tmp {
		ix.ProtoImports = append(ix.ProtoImports, k)
	}
	sort.Strings(ix.ProtoImports)
}

type headerTemplateContext struct {
	Version     string
	PackageName string
	SourceFile  string
	Imports     *importTracker
}

func (module *Module) renderHeader(ctx pgsgo.Context, w io.Writer, in pgs.File, ix *importTracker) error {
	c := &headerTemplateContext{
		Version:     version,
		SourceFile:  in.Name().String(),
		PackageName: ctx.PackageName(in).String(),
		Imports:     ix,
	}

	return templates["header.tmpl"].Execute(w, c)
}
