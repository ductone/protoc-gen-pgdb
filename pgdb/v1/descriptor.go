package v1

import (
	"github.com/jackc/pgtype"
)

// Descriptor is the same for all instances of a Message
type Descriptor interface {
	TableName() string
	Fields() []*Field
	DataField() *Field
	SearchField() *Field
	Indexes() []*Index
}

type Field struct {
	DataType *pgtype.DataType
	Name     string
}

type Index struct {
	Name string
}
