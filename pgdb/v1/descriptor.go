package v1

// Descriptor is the same for all instances of a Message
type Descriptor interface {
	TableName() string
	Fields() []*Field
	DataField() *Field
	SearchField() *Field
	Indexes() []*Index
}

type Field struct {
	Name string
	Type string
}

type Index struct {
	Name string
}
