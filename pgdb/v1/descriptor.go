package v1

import (
	"github.com/doug-martin/goqu/v9/exp"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ColumnSourceKind indicates the origin/nature of a column.
type ColumnSourceKind int32

//nolint:revive // Using underscore naming to match proto enum style.
const (
	// ColumnSourceKind_PROTO_FIELD is a direct mapping to a single proto field.
	// ProtoFieldPath and ProtoPath are populated.
	ColumnSourceKind_PROTO_FIELD ColumnSourceKind = 0

	// ColumnSourceKind_TENANT is the tenant_id column. Maps to a proto field
	// (specified by tenant_id_field option or defaulting to "tenant_id").
	// ProtoFieldPath points to that field.
	ColumnSourceKind_TENANT ColumnSourceKind = 1

	// ColumnSourceKind_PRIMARY_KEY is for computed key columns (pk, sk, pksk, pkskv2).
	// Derived from dynamo key configuration. ProtoFieldPath is nil.
	ColumnSourceKind_PRIMARY_KEY ColumnSourceKind = 2

	// ColumnSourceKind_DATA is the pb_data column containing serialized protobuf.
	ColumnSourceKind_DATA ColumnSourceKind = 3

	// ColumnSourceKind_SEARCH is the fts_data tsvector column for full-text search.
	ColumnSourceKind_SEARCH ColumnSourceKind = 4

	// ColumnSourceKind_ONEOF is the oneof discriminator column (*_oneof).
	// OneofName is populated.
	ColumnSourceKind_ONEOF ColumnSourceKind = 5

	// ColumnSourceKind_VECTOR is an expanded vector column (field_N).
	// ProtoFieldPath points to the source vector field.
	ColumnSourceKind_VECTOR ColumnSourceKind = 6
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
	IsPartitionedByCreatedAt() bool
	GetPartitionedByKsuidFieldName() string

	Indexes(opts ...IndexOptionsFunc) []*Index
	IndexPrimaryKey(opts ...IndexOptionsFunc) *Index

	Statistics(opts ...StatisticOptionsFunc) []*Statistic

	GetPartitionDateRange() MessageOptions_PartitionedByDateRange

	GetStorageParameters() *MessageOptions_StorageParameters
}

type DescriptorFieldOption struct {
	Prefix        string
	IsNested      bool
	IsPartitioned bool

	// ProtoPathPrefix is the dot-delimited proto path prefix for nested messages.
	// e.g., "zoo_shop" means nested field "fur" becomes "zoo_shop.fur"
	ProtoPathPrefix string

	// ProtoFieldPathPrefix is the slice of proto field numbers from root.
	// e.g., [11] means nested field 4 becomes [11, 4]
	ProtoFieldPathPrefix []int32

	// TableOverride overrides the table name for nested columns.
	// When set, nested columns use this table instead of their own message's table.
	TableOverride string
}

func (dfo DescriptorFieldOption) Nullable(defaultValue bool) bool {
	if dfo.IsNested {
		return true
	}
	return defaultValue
}

// ExtendProtoPath combines the prefix with a field's ProtoPath.
func (dfo DescriptorFieldOption) ExtendProtoPath(fieldPath string) string {
	if dfo.ProtoPathPrefix == "" {
		return fieldPath
	}
	if fieldPath == "" {
		return dfo.ProtoPathPrefix
	}
	return dfo.ProtoPathPrefix + "." + fieldPath
}

// ExtendProtoFieldPath combines the prefix with a field's ProtoFieldPath.
func (dfo DescriptorFieldOption) ExtendProtoFieldPath(fieldPath []int32) []int32 {
	if len(dfo.ProtoFieldPathPrefix) == 0 {
		return fieldPath
	}
	if len(fieldPath) == 0 {
		return dfo.ProtoFieldPathPrefix
	}
	result := make([]int32, len(dfo.ProtoFieldPathPrefix)+len(fieldPath))
	copy(result, dfo.ProtoFieldPathPrefix)
	copy(result[len(dfo.ProtoFieldPathPrefix):], fieldPath)
	return result
}

// TableName returns the table override if set, otherwise the default table name.
func (dfo DescriptorFieldOption) TableName(defaultTable string) string {
	if dfo.TableOverride != "" {
		return dfo.TableOverride
	}
	return defaultTable
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

// descriptorFieldProtoPath sets the proto path prefix for nested messages.
func descriptorFieldProtoPath(prefix string) DescriptorFieldOptionFunc {
	return func(option *DescriptorFieldOption) {
		option.ProtoPathPrefix = prefix
	}
}

// descriptorFieldProtoFieldPath sets the proto field path prefix for nested messages.
func descriptorFieldProtoFieldPath(prefix []int32) DescriptorFieldOptionFunc {
	return func(option *DescriptorFieldOption) {
		option.ProtoFieldPathPrefix = prefix
	}
}

// descriptorFieldTableOverride sets the table name override for nested messages.
func descriptorFieldTableOverride(table string) DescriptorFieldOptionFunc {
	return func(option *DescriptorFieldOption) {
		option.TableOverride = table
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

// Nested returns options for a nested message.
// The prefix is the column name prefix (e.g., "11$").
func (r *DescriptorFieldOption) Nested(prefix string) []DescriptorFieldOptionFunc {
	return []DescriptorFieldOptionFunc{
		DescriptorFieldPrefix(r.Prefix + prefix),
		DescriptorFieldNested(true),
	}
}

// NestedWithPath returns options for a nested message with proto path info.
// The prefix is the column name prefix (e.g., "11$"), fieldNum is the proto field number,
// fieldName is the proto field name for the path, and parentTable is the parent message's
// table name (nested columns belong to parent's table).
func (r *DescriptorFieldOption) NestedWithPath(prefix string, fieldNum int32, fieldName string, parentTable string) []DescriptorFieldOptionFunc {
	newProtoPath := r.ProtoPathPrefix
	if newProtoPath == "" {
		newProtoPath = fieldName
	} else {
		newProtoPath = newProtoPath + "." + fieldName
	}

	newProtoFieldPath := make([]int32, len(r.ProtoFieldPathPrefix)+1)
	copy(newProtoFieldPath, r.ProtoFieldPathPrefix)
	newProtoFieldPath[len(r.ProtoFieldPathPrefix)] = fieldNum

	// Propagate table override - if already set, use it; otherwise use parentTable
	tableOverride := r.TableOverride
	if tableOverride == "" {
		tableOverride = parentTable
	}

	return []DescriptorFieldOptionFunc{
		DescriptorFieldPrefix(r.Prefix + prefix),
		DescriptorFieldNested(true),
		descriptorFieldProtoPath(newProtoPath),
		descriptorFieldProtoFieldPath(newProtoFieldPath),
		descriptorFieldTableOverride(tableOverride),
	}
}

type Column struct {
	// === SQL metadata ===

	// Table is the SQL table name this column belongs to.
	Table string

	// Name is the SQL column name (includes prefix like "pb$").
	Name string

	// Type is the PostgreSQL data type (e.g., "text", "int4", "_int4").
	Type string

	// Nullable indicates if the column allows NULL values.
	Nullable bool

	// OverrideExpression is used for generated columns (e.g., GENERATED ALWAYS AS ...).
	OverrideExpression string

	// Default is the default value expression.
	Default string

	// Collation is the column collation (e.g., "C").
	Collation string

	// === Proto field metadata ===

	// SourceKind indicates the origin of this column.
	SourceKind ColumnSourceKind

	// ProtoFieldPath is the path of proto field numbers from root to this column.
	// nil for synthetic columns (PRIMARY_KEY, DATA, SEARCH, ONEOF without backing field).
	ProtoFieldPath []int32

	// ProtoPath is the dot-delimited path of proto field names.
	// Empty for synthetic columns.
	// Examples: "id", "zoo_shop.fur", "zoo_shop.anything.sfixed_64"
	ProtoPath string

	// ProtoKind is the protobuf kind of the leaf field.
	// For synthetic columns, this is the effective kind (e.g., StringKind for pk/sk).
	ProtoKind protoreflect.Kind

	// ProtoTypeName is the fully qualified proto type name:
	//   - For enums: ".models.animals.v1.FurType"
	//   - For messages: ".google.protobuf.Struct"
	//   - For scalars: ""
	ProtoTypeName string

	// IsRepeated is true for repeated/array proto fields.
	IsRepeated bool

	// OneofName is set for columns related to a proto oneof.
	// For discriminator columns: the oneof name.
	// For member columns: the oneof they belong to.
	OneofName string
}

// IsVirtual returns true for columns not directly mapped to a single proto field.
//
// Deprecated: Use SourceKind instead for more precise classification.
func (c *Column) IsVirtual() bool {
	return c.SourceKind != ColumnSourceKind_PROTO_FIELD
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

type Statistic struct {
	Name string
	// this indivates the Stats by this name MAY be dropped, once in IsDropped state, other fields may be empty.
	IsDropped bool
	Kinds     []MessageOptions_Stat_StatsKind
	Columns   []string
}
