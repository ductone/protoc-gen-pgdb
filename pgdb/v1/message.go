package v1

import (
	"github.com/doug-martin/goqu/v9/exp"
)

type DBReflectMessage interface {
	DBReflect() Message
	DBReflectWithDialect(dialect Dialect) Message
}

type Message interface {
	Descriptor() Descriptor

	Record(opts ...RecordOptionsFunc) (exp.Record, error)
	SearchData(opts ...RecordOptionsFunc) []*SearchContent
	Dialect() Dialect
}

type ColumnExpression interface {
	Identifier() exp.IdentifierExpression
}
