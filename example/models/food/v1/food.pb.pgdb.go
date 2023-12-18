// Code generated by protoc-gen-pgdb 0.1.0 from models/food/v1/food.proto. DO NOT EDIT
package v1

import (
	"strings"

	"time"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"

	"github.com/doug-martin/goqu/v9/exp"
	"google.golang.org/protobuf/proto"
)

type pgdbDescriptorPasta struct{}

var (
	instancepgdbDescriptorPasta pgdb_v1.Descriptor = &pgdbDescriptorPasta{}
)

func (d *pgdbDescriptorPasta) TableName() string {
	return "pb_pasta_models_food_v1_29fd1107"
}

func (d *pgdbDescriptorPasta) IsPartitioned() bool {
	return true
}

func (d *pgdbDescriptorPasta) Fields(opts ...pgdb_v1.DescriptorFieldOptionFunc) []*pgdb_v1.Column {
	df := pgdb_v1.NewDescriptorFieldOption(opts)
	_ = df

	rv := make([]*pgdb_v1.Column, 0)

	if !df.IsNested {

		rv = append(rv, &pgdb_v1.Column{
			Name:               df.ColumnName("tenant_id"),
			Type:               "varchar",
			Nullable:           df.Nullable(false),
			OverrideExpression: "",
			Default:            "",
		})

	}

	if !df.IsNested {

		rv = append(rv, &pgdb_v1.Column{
			Name:               df.ColumnName("pksk"),
			Type:               "varchar",
			Nullable:           df.Nullable(false),
			OverrideExpression: "varchar GENERATED ALWAYS AS (pb$pk || '|' || pb$sk) STORED",
			Default:            "",
		})

	}

	if !df.IsNested {

		rv = append(rv, &pgdb_v1.Column{
			Name:               df.ColumnName("pk"),
			Type:               "varchar",
			Nullable:           df.Nullable(false),
			OverrideExpression: "",
			Default:            "",
		})

	}

	if !df.IsNested {

		rv = append(rv, &pgdb_v1.Column{
			Name:               df.ColumnName("sk"),
			Type:               "varchar",
			Nullable:           df.Nullable(false),
			OverrideExpression: "",
			Default:            "",
		})

	}

	if !df.IsNested {

		rv = append(rv, &pgdb_v1.Column{
			Name:               df.ColumnName("fts_data"),
			Type:               "tsvector",
			Nullable:           df.Nullable(true),
			OverrideExpression: "",
			Default:            "",
		})

	}

	if !df.IsNested {

		rv = append(rv, &pgdb_v1.Column{
			Name:               df.ColumnName("pb_data"),
			Type:               "bytea",
			Nullable:           df.Nullable(false),
			OverrideExpression: "",
			Default:            "",
		})

	}

	rv = append(rv, &pgdb_v1.Column{
		Name:               df.ColumnName("id"),
		Type:               "text",
		Nullable:           df.Nullable(false),
		OverrideExpression: "",
		Default:            "''",
	})

	rv = append(rv, &pgdb_v1.Column{
		Name:               df.ColumnName("created_at"),
		Type:               "timestamptz",
		Nullable:           df.Nullable(true),
		OverrideExpression: "",
		Default:            "",
	})

	return rv
}

func (d *pgdbDescriptorPasta) PKSKField() *pgdb_v1.Column {
	return &pgdb_v1.Column{
		Table: "pb_pasta_models_food_v1_29fd1107",
		Name:  "pb$pksk",
		Type:  "varchar",
	}
}

func (d *pgdbDescriptorPasta) DataField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Table: "pb_pasta_models_food_v1_29fd1107", Name: "pb$pb_data", Type: "bytea"}
}

func (d *pgdbDescriptorPasta) SearchField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Table: "pb_pasta_models_food_v1_29fd1107", Name: "pb$fts_data", Type: "tsvector"}
}

func (d *pgdbDescriptorPasta) VersioningField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Table: "pb_pasta_models_food_v1_29fd1107", Name: "pb$created_at", Type: "timestamptz"}
}

func (d *pgdbDescriptorPasta) TenantField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Table: "pb_pasta_models_food_v1_29fd1107", Name: "pb$tenant_id", Type: "varchar"}
}

func (d *pgdbDescriptorPasta) IndexPrimaryKey(opts ...pgdb_v1.IndexOptionsFunc) *pgdb_v1.Index {
	io := pgdb_v1.NewIndexOptions(opts)
	_ = io

	return &pgdb_v1.Index{
		Name:               io.IndexName("pksk_pasta_models_food_v1_441e44c9"),
		Method:             pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE,
		IsPrimary:          true,
		IsUnique:           true,
		IsDropped:          false,
		Columns:            []string{io.ColumnName("tenant_id"), io.ColumnName("pksk")},
		OverrideExpression: "",
	}

}

func (d *pgdbDescriptorPasta) Indexes(opts ...pgdb_v1.IndexOptionsFunc) []*pgdb_v1.Index {
	io := pgdb_v1.NewIndexOptions(opts)
	_ = io
	rv := make([]*pgdb_v1.Index, 0)

	if !io.IsNested {

		rv = append(rv, &pgdb_v1.Index{
			Name:               io.IndexName("pksk_pasta_models_food_v1_441e44c9"),
			Method:             pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE,
			IsPrimary:          true,
			IsUnique:           true,
			IsDropped:          false,
			Columns:            []string{io.ColumnName("tenant_id"), io.ColumnName("pksk")},
			OverrideExpression: "",
		})

	}

	if !io.IsNested {

		rv = append(rv, &pgdb_v1.Index{
			Name:               io.IndexName("pksk_split_pasta_models_food_v1_4c4cc274"),
			Method:             pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE,
			IsPrimary:          false,
			IsUnique:           false,
			IsDropped:          true,
			Columns:            []string{io.ColumnName("tenant_id"), io.ColumnName("pk"), io.ColumnName("sk")},
			OverrideExpression: "",
		})

	}

	if !io.IsNested {

		rv = append(rv, &pgdb_v1.Index{
			Name:               io.IndexName("pksk_split2_pasta_models_food_v1_65c526aa"),
			Method:             pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE,
			IsPrimary:          false,
			IsUnique:           false,
			IsDropped:          false,
			Columns:            []string{io.ColumnName("tenant_id"), io.ColumnName("pk"), io.ColumnName("sk")},
			OverrideExpression: "",
		})

	}

	if !io.IsNested {

		rv = append(rv, &pgdb_v1.Index{
			Name:               io.IndexName("fts_data_pasta_models_food_v1_77400ba0"),
			Method:             pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE_GIN,
			IsPrimary:          false,
			IsUnique:           false,
			IsDropped:          false,
			Columns:            []string{io.ColumnName("tenant_id"), io.ColumnName("fts_data")},
			OverrideExpression: "",
		})

	}

	return rv
}

type pgdbMessagePasta struct {
	self *Pasta
}

func (dbr *Pasta) DBReflect() pgdb_v1.Message {
	return &pgdbMessagePasta{
		self: dbr,
	}
}

func (m *pgdbMessagePasta) Descriptor() pgdb_v1.Descriptor {
	return instancepgdbDescriptorPasta
}

func (m *pgdbMessagePasta) Record(opts ...pgdb_v1.RecordOptionsFunc) (exp.Record, error) {
	ro := pgdb_v1.NewRecordOptions(opts)
	_ = ro
	nullExp := exp.NewLiteralExpression("NULL")
	_ = nullExp

	var sb strings.Builder

	rv := exp.Record{}

	if !ro.IsNested {

		cfv0 := string(m.self.TenantId)

		if ro.Nulled {
			rv[ro.ColumnName("tenant_id")] = nullExp
		} else {
			rv[ro.ColumnName("tenant_id")] = cfv0
		}

	}

	if !ro.IsNested {

	}

	if !ro.IsNested {

		sb.Reset()

		_, _ = sb.WriteString("models_food_v1_pasta")

		_, _ = sb.WriteString(":")

		_, _ = sb.WriteString(m.self.TenantId)

		_, _ = sb.WriteString(":")

		_, _ = sb.WriteString(m.self.Id)

		cfv2 := sb.String()

		if ro.Nulled {
			rv[ro.ColumnName("pk")] = nullExp
		} else {
			rv[ro.ColumnName("pk")] = cfv2
		}

	}

	if !ro.IsNested {

		sb.Reset()

		_, _ = sb.WriteString("example")

		cfv3 := sb.String()

		if ro.Nulled {
			rv[ro.ColumnName("sk")] = nullExp
		} else {
			rv[ro.ColumnName("sk")] = cfv3
		}

	}

	if !ro.IsNested {

		cfv4tmp := []*pgdb_v1.SearchContent{

			{
				Type:   pgdb_v1.FieldOptions_FULL_TEXT_TYPE_EXACT,
				Weight: pgdb_v1.FieldOptions_FULL_TEXT_WEIGHT_UNSPECIFIED,
				Value:  m.self.GetId(),
			},
		}

		cfv4 := pgdb_v1.FullTextSearchVectors(cfv4tmp)

		if ro.Nulled {
			rv[ro.ColumnName("fts_data")] = nullExp
		} else {
			rv[ro.ColumnName("fts_data")] = cfv4
		}

	}

	if !ro.IsNested {

		cfv5, err := proto.Marshal(m.self)
		if err != nil {
			return nil, err
		}

		if ro.Nulled {
			rv[ro.ColumnName("pb_data")] = nullExp
		} else {
			rv[ro.ColumnName("pb_data")] = cfv5
		}

	}

	v1 := string(m.self.GetId())

	if ro.Nulled {
		rv[ro.ColumnName("id")] = nullExp
	} else {
		rv[ro.ColumnName("id")] = v1
	}

	var v2 *time.Time
	if m.self.GetCreatedAt().IsValid() {
		v2tmp := m.self.GetCreatedAt().AsTime()
		v2 = &v2tmp
	}

	if ro.Nulled {
		rv[ro.ColumnName("created_at")] = nullExp
	} else {
		rv[ro.ColumnName("created_at")] = v2
	}

	return rv, nil
}

func (m *pgdbMessagePasta) SearchData(opts ...pgdb_v1.RecordOptionsFunc) []*pgdb_v1.SearchContent {
	rv := []*pgdb_v1.SearchContent{

		{
			Type:   pgdb_v1.FieldOptions_FULL_TEXT_TYPE_EXACT,
			Weight: pgdb_v1.FieldOptions_FULL_TEXT_WEIGHT_UNSPECIFIED,
			Value:  m.self.GetId(),
		},
	}

	return rv
}

type PastaDB struct {
	tableName string
}

type PastaDBQueryBuilder struct {
	tableName string
}

type PastaDBQueryUnsafe struct {
	tableName string
}

type PastaDBColumns struct {
	tableName string
}

func (x *Pasta) DB() *PastaDB {
	return &PastaDB{tableName: x.DBReflect().Descriptor().TableName()}
}

func (x *PastaDB) TableName() string {
	return x.tableName
}

func (x *PastaDB) Query() *PastaDBQueryBuilder {
	return &PastaDBQueryBuilder{tableName: x.tableName}
}

func (x *PastaDB) Columns() *PastaDBColumns {
	return &PastaDBColumns{tableName: x.tableName}
}

func (x *PastaDB) WithTable(t string) *PastaDB {
	return &PastaDB{tableName: t}
}

func (x *PastaDBQueryBuilder) WithTable(t string) *PastaDBQueryBuilder {
	return &PastaDBQueryBuilder{tableName: t}
}

func (x *PastaDBQueryBuilder) Unsafe() *PastaDBQueryUnsafe {
	return &PastaDBQueryUnsafe{tableName: x.tableName}
}

type PastaTenantIdSafeOperators struct {
	column    string
	tableName string
}

func (x *PastaTenantIdSafeOperators) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

func (x *PastaTenantIdSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Eq(v)
}

func (x *PastaTenantIdSafeOperators) Gt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gt(v)
}

func (x *PastaTenantIdSafeOperators) Gte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gte(v)
}

func (x *PastaTenantIdSafeOperators) Lt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Lt(v)
}

func (x *PastaTenantIdSafeOperators) Lte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Lte(v)
}

func (x *PastaTenantIdSafeOperators) In(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).In(v)
}

func (x *PastaTenantIdSafeOperators) NotIn(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).NotIn(v)
}

func (x *PastaTenantIdSafeOperators) IsNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).IsNull()
}

func (x *PastaTenantIdSafeOperators) IsNotEmpty() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gt("")
}

func (x *PastaTenantIdSafeOperators) IsNotNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).IsNotNull()
}

func (x *PastaTenantIdSafeOperators) Between(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Between(exp.NewRangeVal(start, end))
}

func (x *PastaTenantIdSafeOperators) NotBetween(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).NotBetween(exp.NewRangeVal(start, end))
}

func (x *PastaDBQueryBuilder) TenantId() *PastaTenantIdSafeOperators {
	return &PastaTenantIdSafeOperators{tableName: x.tableName, column: "pb$" + "tenant_id"}
}

type PastaPKSKSafeOperators struct {
	column    string
	tableName string
}

func (x *PastaPKSKSafeOperators) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

func (x *PastaPKSKSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Eq(v)
}

func (x *PastaPKSKSafeOperators) Gt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gt(v)
}

func (x *PastaPKSKSafeOperators) Gte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gte(v)
}

func (x *PastaPKSKSafeOperators) Lt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Lt(v)
}

func (x *PastaPKSKSafeOperators) Lte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Lte(v)
}

func (x *PastaPKSKSafeOperators) In(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).In(v)
}

func (x *PastaPKSKSafeOperators) NotIn(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).NotIn(v)
}

func (x *PastaPKSKSafeOperators) IsNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).IsNull()
}

func (x *PastaPKSKSafeOperators) IsNotEmpty() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gt("")
}

func (x *PastaPKSKSafeOperators) IsNotNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).IsNotNull()
}

func (x *PastaPKSKSafeOperators) Between(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Between(exp.NewRangeVal(start, end))
}

func (x *PastaPKSKSafeOperators) NotBetween(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).NotBetween(exp.NewRangeVal(start, end))
}

func (x *PastaDBQueryBuilder) PKSK() *PastaPKSKSafeOperators {
	return &PastaPKSKSafeOperators{tableName: x.tableName, column: "pb$" + "pksk"}
}

type PastaPKSafeOperators struct {
	column    string
	tableName string
}

func (x *PastaPKSafeOperators) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

func (x *PastaPKSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Eq(v)
}

func (x *PastaPKSafeOperators) Gt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gt(v)
}

func (x *PastaPKSafeOperators) Gte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gte(v)
}

func (x *PastaPKSafeOperators) Lt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Lt(v)
}

func (x *PastaPKSafeOperators) Lte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Lte(v)
}

func (x *PastaPKSafeOperators) In(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).In(v)
}

func (x *PastaPKSafeOperators) NotIn(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).NotIn(v)
}

func (x *PastaPKSafeOperators) IsNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).IsNull()
}

func (x *PastaPKSafeOperators) IsNotEmpty() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gt("")
}

func (x *PastaPKSafeOperators) IsNotNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).IsNotNull()
}

func (x *PastaPKSafeOperators) Between(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Between(exp.NewRangeVal(start, end))
}

func (x *PastaPKSafeOperators) NotBetween(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).NotBetween(exp.NewRangeVal(start, end))
}

func (x *PastaDBQueryBuilder) PK() *PastaPKSafeOperators {
	return &PastaPKSafeOperators{tableName: x.tableName, column: "pb$" + "pk"}
}

type PastaSKSafeOperators struct {
	column    string
	tableName string
}

func (x *PastaSKSafeOperators) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

func (x *PastaSKSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Eq(v)
}

func (x *PastaSKSafeOperators) Gt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gt(v)
}

func (x *PastaSKSafeOperators) Gte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gte(v)
}

func (x *PastaSKSafeOperators) Lt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Lt(v)
}

func (x *PastaSKSafeOperators) Lte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Lte(v)
}

func (x *PastaSKSafeOperators) In(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).In(v)
}

func (x *PastaSKSafeOperators) NotIn(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).NotIn(v)
}

func (x *PastaSKSafeOperators) IsNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).IsNull()
}

func (x *PastaSKSafeOperators) IsNotEmpty() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gt("")
}

func (x *PastaSKSafeOperators) IsNotNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).IsNotNull()
}

func (x *PastaSKSafeOperators) Between(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Between(exp.NewRangeVal(start, end))
}

func (x *PastaSKSafeOperators) NotBetween(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).NotBetween(exp.NewRangeVal(start, end))
}

func (x *PastaDBQueryBuilder) SK() *PastaSKSafeOperators {
	return &PastaSKSafeOperators{tableName: x.tableName, column: "pb$" + "sk"}
}

type PastaFTSDataSafeOperators struct {
	column    string
	tableName string
}

func (x *PastaFTSDataSafeOperators) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

func (x *PastaFTSDataSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Eq(v)
}

func (x *PastaDBQueryBuilder) FTSData() *PastaFTSDataSafeOperators {
	return &PastaFTSDataSafeOperators{tableName: x.tableName, column: "pb$" + "fts_data"}
}

type PastaTenantIdQueryType struct {
	column    string
	tableName string
}

func (x *PastaDBQueryUnsafe) TenantId() *PastaTenantIdQueryType {
	return &PastaTenantIdQueryType{tableName: x.tableName, column: "pb$" + "tenant_id"}
}

func (x *PastaTenantIdQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

type PastaPKSKQueryType struct {
	column    string
	tableName string
}

func (x *PastaDBQueryUnsafe) PKSK() *PastaPKSKQueryType {
	return &PastaPKSKQueryType{tableName: x.tableName, column: "pb$" + "pksk"}
}

func (x *PastaPKSKQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

type PastaPKQueryType struct {
	column    string
	tableName string
}

func (x *PastaDBQueryUnsafe) PK() *PastaPKQueryType {
	return &PastaPKQueryType{tableName: x.tableName, column: "pb$" + "pk"}
}

func (x *PastaPKQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

type PastaSKQueryType struct {
	column    string
	tableName string
}

func (x *PastaDBQueryUnsafe) SK() *PastaSKQueryType {
	return &PastaSKQueryType{tableName: x.tableName, column: "pb$" + "sk"}
}

func (x *PastaSKQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

type PastaFTSDataQueryType struct {
	column    string
	tableName string
}

func (x *PastaDBQueryUnsafe) FTSData() *PastaFTSDataQueryType {
	return &PastaFTSDataQueryType{tableName: x.tableName, column: "pb$" + "fts_data"}
}

func (x *PastaFTSDataQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

type PastaPBDataQueryType struct {
	column    string
	tableName string
}

func (x *PastaDBQueryUnsafe) PBData() *PastaPBDataQueryType {
	return &PastaPBDataQueryType{tableName: x.tableName, column: "pb$" + "pb_data"}
}

func (x *PastaPBDataQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

type PastaIdQueryType struct {
	column    string
	tableName string
}

func (x *PastaDBQueryUnsafe) Id() *PastaIdQueryType {
	return &PastaIdQueryType{tableName: x.tableName, column: "pb$" + "id"}
}

func (x *PastaIdQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

type PastaCreatedAtQueryType struct {
	column    string
	tableName string
}

func (x *PastaDBQueryUnsafe) CreatedAt() *PastaCreatedAtQueryType {
	return &PastaCreatedAtQueryType{tableName: x.tableName, column: "pb$" + "created_at"}
}

func (x *PastaCreatedAtQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

func (x *PastaDBColumns) WithTable(t string) *PastaDBColumns {
	return &PastaDBColumns{tableName: t}
}

func (x *PastaDBColumns) TenantId() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "tenant_id")
}

func (x *PastaDBColumns) PKSK() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "pksk")
}

func (x *PastaDBColumns) PK() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "pk")
}

func (x *PastaDBColumns) SK() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "sk")
}

func (x *PastaDBColumns) FTSData() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "fts_data")
}

func (x *PastaDBColumns) PBData() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "pb_data")
}

func (x *PastaDBColumns) Id() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "id")
}

func (x *PastaDBColumns) CreatedAt() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "created_at")
}