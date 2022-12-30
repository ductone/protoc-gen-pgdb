package pgdb

import (
	"fmt"
	"io"
	"sort"

	"github.com/ductone/protoc-gen-pgdb/internal/slice"
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
	Strconv              bool
	Time                 bool
	PgType               bool

	ctx          pgsgo.Context
	input        pgs.File
	ProtoImports []ImportAlias
}

type ImportAlias struct {
	Import string
	Dummy  string
}

func (ix *importTracker) AddProtoEntity(entity pgs.Field) {
	elist := []pgs.Entity{entity}
	elist = append(elist, slice.Convert(entity.Imports(), func(file pgs.File) pgs.Entity {
		return file
	})...)

	tmp := make(map[ImportAlias]struct{})
	for _, v := range ix.ProtoImports {
		tmp[v] = struct{}{}
	}

	for _, entity := range elist {
		// TOD(pquerna): is there a better method to do this?
		if ix.ctx.ImportPath(entity) == ix.ctx.ImportPath(ix.input) {
			// os.Stderr.WriteString(fmt.Sprintf("[%p] [skip] AddProtoEntity: %s == %s\n", ix, ix.ctx.ImportPath(entity), ix.ctx.ImportPath(ix.input)))
			// os.Stderr.WriteString(fmt.Sprintf("[%p] [skip]     for %s -> %s\n", ix, entity.FullyQualifiedName(), string(ix.ctx.ImportPath(entity))))

			continue
		}
		// os.Stderr.WriteString(fmt.Sprintf("[add] AddProtoEntity: %s != %s\n", ix.ctx.ImportPath(entity), ix.ctx.ImportPath(ix.input)))
		key := fmt.Sprintf(`%s "%s"`, ix.ctx.PackageName(entity), string(ix.ctx.ImportPath(entity)))
		ia := ImportAlias{
			Import: key,
			//			Dummy:  string(ix.ctx.Type(entity)),
		}
		tmp[ia] = struct{}{}
	}

	ix.ProtoImports = make([]ImportAlias, 0, len(tmp))
	for k := range tmp {
		ix.ProtoImports = append(ix.ProtoImports, k)
	}

	sort.Slice(ix.ProtoImports, func(i, j int) bool {
		return ix.ProtoImports[i].Import > ix.ProtoImports[j].Import
	})
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
