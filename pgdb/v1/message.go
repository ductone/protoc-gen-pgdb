package v1

import (
	"github.com/doug-martin/goqu/v9/exp"
)

type DBReflectMessage interface {
	DBReflect() Message
}

type Message interface {
	Descriptor() Descriptor

	Record(opts ...QueryOptions) (exp.Record, error)
	ConflictUpdateExpression(opts ...QueryOptions) (exp.ConflictUpdateExpression, error)
}
