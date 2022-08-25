package pgdb

import (
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
	ColumnFields []string
}

func (module *Module) getQueryBuilder(ctx pgsgo.Context, msg pgs.Message, ix *importTracker) *qbContext {
	return &qbContext{}
}
