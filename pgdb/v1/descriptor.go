package v1

// Descriptor is the same for all instances of a Message.
type Descriptor interface {
	TableName() string
	Fields() []*Column
	DataField() *Column
	SearchField() *Column
	Indexes() []*Index
}

type Column struct {
	Name               string
	Type               string
	Nullable           bool
	OverrideExpression string
}

type Index struct {
	Name string
	// only one index can be set to IsPrimary
	IsPrimary bool
	IsUnique  bool
	// this indivates the Index by this name MAY be dropped, once in IsDropped state, other fields may be empty.
	IsDropped bool
	Method    MessageOptions_Index_IndexMethod
	Columns   []string
	// OverrideExpression if set, this string is used to render indexes contents, instead of the Columns list.
	OverrideExpression string
	WherePredicate     string
}
