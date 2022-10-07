package v1

import (
	"github.com/doug-martin/goqu/v9/exp"
)

type DBReflectMessage interface {
	DBReflect() Message
}

type Message interface {
	Descriptor() Descriptor

	Record(opts ...RecordOptionsFunc) (exp.Record, error)
	ConflictUpdateExpression(opts ...RecordOptionsFunc) (exp.ConflictUpdateExpression, error)
}
