package pgdb

import (
	"io"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type qbContext struct {
	ReceiverType string
	DbType       string
	QueryType    string
	QueryFields  []string
	UnsafeType   string
	UnsafeFields []string
	ColumnType   string
	ColumnFields []*fieldContext
}

func (module *Module) renderQueryBuilder(ctx pgsgo.Context, w io.Writer, in pgs.File, m pgs.Message, ix *importTracker) error {
	c := module.getQueryBuilder(ctx, m, ix)
	return templates["query_builder.tmpl"].Execute(w, c)
}

func (module *Module) getQueryBuilder(ctx pgsgo.Context, m pgs.Message, ix *importTracker) *qbContext {
	nsgFields := module.getMessageFields(ctx, m, ix, "m.self")
	return &qbContext{
		ReceiverType: ctx.Name(m).String(),
		DbType:       getDbType(ctx, m),
		QueryType:    getQueryType(ctx, m),
		UnsafeType:   getUnsafeType(ctx, m),
		ColumnType:   getColumnType(ctx, m),
		ColumnFields: nsgFields,
	}
}

func getDbType(ctx pgsgo.Context, m pgs.Message) string {
	return ctx.Name(m).String() + "DB"
}

func getQueryType(ctx pgsgo.Context, m pgs.Message) string {
	return ctx.Name(m).String() + "DBQueryBuilder"
}

func getUnsafeType(ctx pgsgo.Context, m pgs.Message) string {
	return ctx.Name(m).String() + "DBQueryUnsafe"
}

func getColumnType(ctx pgsgo.Context, m pgs.Message) string {
	return ctx.Name(m).String() + "DBColumns"
}
