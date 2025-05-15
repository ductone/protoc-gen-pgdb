package pgdb

import (
	"io"
	"sort"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
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
	NetIP                bool

	ctx   pgsgo.Context
	input pgs.File

	typeMapper map[pgs.Name]pgs.FilePath
}

type ImportAlias struct {
	Import string
	Alias  string
}

func (ix *importTracker) ProtoImports() []ImportAlias {
	rv := make([]ImportAlias, 0, len(ix.typeMapper))
	for k, v := range ix.typeMapper {
		rv = append(rv, ImportAlias{
			Import: v.String(),
			Alias:  k.String(),
		})
	}
	sort.Slice(rv, func(i, j int) bool {
		return rv[i].Alias > rv[j].Alias
	})
	return rv
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
