package pgdb

import (
	"io"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

type dropTemplateContext struct {
	ReceiverType string
	TableName    string
}

// renderDrop generates the DROP TABLE / TRUNCATE TABLE helpers for a message
// that has the `drop_enabled` message option set. When a message is
// drop_enabled, this is the only code generated for it.
func (module *Module) renderDrop(ctx pgsgo.Context, w io.Writer, in pgs.File, m pgs.Message, ix *importTracker) error {
	ix.PGDBV1 = true

	tableName, err := getTableName(m)
	if err != nil {
		return err
	}

	c := &dropTemplateContext{
		ReceiverType: ctx.Name(m).String(),
		TableName:    tableName,
	}

	return templates["drop.tmpl"].Execute(w, c)
}
