package pgdb

import (
	"io"

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
