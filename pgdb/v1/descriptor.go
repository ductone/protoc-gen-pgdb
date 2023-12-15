package v1

import (
	"github.com/doug-martin/goqu/v9/exp"
)

// Descriptor is the same for all instances of a Message.
type Descriptor interface {
	TableName() string

	Fields(opts ...DescriptorFieldOptionFunc) []*Column

	PKSKField() *Column
	DataField() *Column
	SearchField() *Column
	VersioningField() *Column
	TenantField() *Column
	IsPartitioned() bool

	Indexes(opts ...IndexOptionsFunc) []*Index
	IndexPrimaryKey(opts ...IndexOptionsFunc) *Index
}

type DescriptorFieldOption struct {
	Prefix        string
	IsNested      bool
	IsPartitioned bool
}

func (dfo DescriptorFieldOption) Nullable(defaultValue bool) bool {
	if dfo.IsNested {
		return true
	}
	return defaultValue
}

type DescriptorFieldOptionFunc func(option *DescriptorFieldOption)

func DescriptorFieldPrefix(prefix string) DescriptorFieldOptionFunc {
	return func(option *DescriptorFieldOption) {
		option.Prefix = prefix
	}
}

func DescriptorFieldNested(b bool) DescriptorFieldOptionFunc {
	return func(option *DescriptorFieldOption) {
		option.IsNested = b
	}
}

func NewDescriptorFieldOption(opts []DescriptorFieldOptionFunc) *DescriptorFieldOption {
	option := &DescriptorFieldOption{
		Prefix: "pb$",
	}
	for _, opt := range opts {
		opt(option)
	}
	return option
}

func (r *DescriptorFieldOption) ColumnName(in string) string {
	return r.Prefix + in
}

func (r *DescriptorFieldOption) Nested(prefix string) []DescriptorFieldOptionFunc {
	return []DescriptorFieldOptionFunc{
		DescriptorFieldPrefix(r.Prefix + prefix),
		DescriptorFieldNested(true),
	}
}

type Column struct {
	Table              string
	Name               string
	Type               string
	Nullable           bool
	OverrideExpression string
	Default            string
}

func (x *Column) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.Table, x.Name)
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
