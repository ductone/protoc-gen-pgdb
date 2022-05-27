package v1

import (
	"github.com/doug-martin/goqu/v9/exp"
)

type PGDBMessage interface {
	PGDBReflect() Message
}

type Message interface {
	Descriptor() Descriptor

	Record(opts ...QueryOptions) exp.Record
	ConflictUpdateExpression(opts ...QueryOptions) exp.ConflictUpdateExpression
}
