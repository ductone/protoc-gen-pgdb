// Code generated by protoc-gen-pgdb 0.1.0 from models/zoo/v1/zoo.proto. DO NOT EDIT
package v1

import (
	"strings"

	"time"

	animals_v1 "github.com/ductone/protoc-gen-pgdb/example/models/animals/v1"

	"github.com/doug-martin/goqu/v9/exp"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	"google.golang.org/protobuf/proto"
)

type pgdbDescriptorShop struct{}

var (
	instancepgdbDescriptorShop pgdb_v1.Descriptor = &pgdbDescriptorShop{}
)

func (d *pgdbDescriptorShop) TableName() string {
	return "pb_shop_models_zoo_v1_ca2425f6"
}

func (d *pgdbDescriptorShop) IsPartitioned() bool {
	return false
}

func (d *pgdbDescriptorShop) Fields(opts ...pgdb_v1.DescriptorFieldOptionFunc) []*pgdb_v1.Column {
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

	rv = append(rv, &pgdb_v1.Column{
		Name:               df.ColumnName("fur"),
		Type:               "int4",
		Nullable:           df.Nullable(false),
		OverrideExpression: "",
		Default:            "0",
	})

	rv = append(rv, &pgdb_v1.Column{
		Name:               df.ColumnName("medium_oneof"),
		Type:               "int4",
		Nullable:           df.Nullable(false),
		OverrideExpression: "",
		Default:            "0",
	})

	rv = append(rv, ((*animals_v1.PaperBook)(nil)).DBReflect().Descriptor().Fields(df.Nested("50$")...)...)

	rv = append(rv, ((*animals_v1.EBook)(nil)).DBReflect().Descriptor().Fields(df.Nested("51$")...)...)

	rv = append(rv, ((*animals_v1.ScalarValue)(nil)).DBReflect().Descriptor().Fields(df.Nested("52$")...)...)

	rv = append(rv, ((*Shop_Manager)(nil)).DBReflect().Descriptor().Fields(df.Nested("5$")...)...)

	return rv
}

func (d *pgdbDescriptorShop) PKSKField() *pgdb_v1.Column {
	return &pgdb_v1.Column{
		Table: "pb_shop_models_zoo_v1_ca2425f6",
		Name:  "pb$pksk",
		Type:  "varchar",
	}
}

func (d *pgdbDescriptorShop) DataField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Table: "pb_shop_models_zoo_v1_ca2425f6", Name: "pb$pb_data", Type: "bytea"}
}

func (d *pgdbDescriptorShop) SearchField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Table: "pb_shop_models_zoo_v1_ca2425f6", Name: "pb$fts_data", Type: "tsvector"}
}

func (d *pgdbDescriptorShop) VersioningField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Table: "pb_shop_models_zoo_v1_ca2425f6", Name: "pb$created_at", Type: "timestamptz"}
}

func (d *pgdbDescriptorShop) TenantField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Table: "pb_shop_models_zoo_v1_ca2425f6", Name: "pb$tenant_id", Type: "varchar"}
}

func (d *pgdbDescriptorShop) IndexPrimaryKey(opts ...pgdb_v1.IndexOptionsFunc) *pgdb_v1.Index {
	io := pgdb_v1.NewIndexOptions(opts)
	_ = io

	return &pgdb_v1.Index{
		Name:               io.IndexName("pksk_shop_models_zoo_v1_a0f13ce7"),
		Method:             pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE,
		IsPrimary:          true,
		IsUnique:           true,
		IsDropped:          false,
		Columns:            []string{io.ColumnName("tenant_id"), io.ColumnName("pksk")},
		OverrideExpression: "",
	}

}

func (d *pgdbDescriptorShop) Indexes(opts ...pgdb_v1.IndexOptionsFunc) []*pgdb_v1.Index {
	io := pgdb_v1.NewIndexOptions(opts)
	_ = io
	rv := make([]*pgdb_v1.Index, 0)

	if !io.IsNested {

		rv = append(rv, &pgdb_v1.Index{
			Name:               io.IndexName("pksk_shop_models_zoo_v1_a0f13ce7"),
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
			Name:               io.IndexName("pksk_split_shop_models_zoo_v1_c45b3ad3"),
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
			Name:               io.IndexName("pksk_split2_shop_models_zoo_v1_3fd2424c"),
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
			Name:               io.IndexName("fts_data_shop_models_zoo_v1_1b685f12"),
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

func (d *pgdbDescriptorShop) Statistics(opts ...pgdb_v1.StatisticOptionsFunc) []*pgdb_v1.Statistic {
	io := pgdb_v1.NewStatisticOption(opts)
	_ = io
	rv := make([]*pgdb_v1.Statistic, 0)

	return rv
}

type ShopMediumType int32

var ShopMedium = struct {
	Paper    ShopMediumType
	Ebook    ShopMediumType
	Anything ShopMediumType
}{
	Paper:    50,
	Ebook:    51,
	Anything: 52,
}

type pgdbMessageShop struct {
	self *Shop
}

func (dbr *Shop) DBReflect() pgdb_v1.Message {
	return &pgdbMessageShop{
		self: dbr,
	}
}

func (m *pgdbMessageShop) Descriptor() pgdb_v1.Descriptor {
	return instancepgdbDescriptorShop
}

func (m *pgdbMessageShop) Record(opts ...pgdb_v1.RecordOptionsFunc) (exp.Record, error) {
	ro := pgdb_v1.NewRecordOptions(opts)
	_ = ro
	nullExp := exp.NewLiteralExpression("NULL")
	_ = nullExp

	var sb strings.Builder

	rv := exp.Record{}

	if !ro.IsNested {

		cfv0 := strings.ReplaceAll(string(m.self.TenantId), "\u0000", "")

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

		_, _ = sb.WriteString("models_zoo_v1_shop")

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

		cfv4tmp = append(cfv4tmp, m.self.GetPaper().DBReflect().SearchData()...)

		cfv4tmp = append(cfv4tmp, m.self.GetEbook().DBReflect().SearchData()...)

		cfv4tmp = append(cfv4tmp, m.self.GetAnything().DBReflect().SearchData()...)

		cfv4tmp = append(cfv4tmp, m.self.GetMgr().DBReflect().SearchData()...)

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

	v1 := strings.ReplaceAll(string(m.self.GetId()), "\u0000", "")

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

	v3tmp := m.self.GetPaper()
	v3opts := ro.Nested("50$")
	if v3tmp == nil {
		v3opts = append(v3opts, pgdb_v1.RecordOptionNulled(true))
	}

	v3, err := pgdb_v1.MarshalNestedRecord(v3tmp, v3opts...)
	if err != nil {
		return nil, err
	}

	for k, v := range v3 {
		if ro.Nulled {
			rv[k] = nullExp
		} else {
			rv[k] = v
		}
	}

	v4tmp := m.self.GetEbook()
	v4opts := ro.Nested("51$")
	if v4tmp == nil {
		v4opts = append(v4opts, pgdb_v1.RecordOptionNulled(true))
	}

	v4, err := pgdb_v1.MarshalNestedRecord(v4tmp, v4opts...)
	if err != nil {
		return nil, err
	}

	for k, v := range v4 {
		if ro.Nulled {
			rv[k] = nullExp
		} else {
			rv[k] = v
		}
	}

	v5tmp := m.self.GetAnything()
	v5opts := ro.Nested("52$")
	if v5tmp == nil {
		v5opts = append(v5opts, pgdb_v1.RecordOptionNulled(true))
	}

	v5, err := pgdb_v1.MarshalNestedRecord(v5tmp, v5opts...)
	if err != nil {
		return nil, err
	}

	for k, v := range v5 {
		if ro.Nulled {
			rv[k] = nullExp
		} else {
			rv[k] = v
		}
	}

	v6 := int32(m.self.GetFur())

	if ro.Nulled {
		rv[ro.ColumnName("fur")] = nullExp
	} else {
		rv[ro.ColumnName("fur")] = v6
	}

	v7tmp := m.self.GetMgr()
	v7opts := ro.Nested("5$")
	if v7tmp == nil {
		v7opts = append(v7opts, pgdb_v1.RecordOptionNulled(true))
	}

	v7, err := pgdb_v1.MarshalNestedRecord(v7tmp, v7opts...)
	if err != nil {
		return nil, err
	}

	for k, v := range v7 {
		if ro.Nulled {
			rv[k] = nullExp
		} else {
			rv[k] = v
		}
	}

	oneof1 := uint32(0)

	switch m.self.GetMedium().(type) {

	case *Shop_Paper:
		oneof1 = 50

	case *Shop_Ebook:
		oneof1 = 51

	case *Shop_Anything:
		oneof1 = 52

	}

	if ro.Nulled {
		rv[ro.ColumnName("medium_oneof")] = nullExp
	} else {
		rv[ro.ColumnName("medium_oneof")] = oneof1
	}

	return rv, nil
}

func (m *pgdbMessageShop) SearchData(opts ...pgdb_v1.RecordOptionsFunc) []*pgdb_v1.SearchContent {
	rv := []*pgdb_v1.SearchContent{

		{
			Type:   pgdb_v1.FieldOptions_FULL_TEXT_TYPE_EXACT,
			Weight: pgdb_v1.FieldOptions_FULL_TEXT_WEIGHT_UNSPECIFIED,
			Value:  m.self.GetId(),
		},
	}

	rv = append(rv, m.self.GetPaper().DBReflect().SearchData()...)

	rv = append(rv, m.self.GetEbook().DBReflect().SearchData()...)

	rv = append(rv, m.self.GetAnything().DBReflect().SearchData()...)

	rv = append(rv, m.self.GetMgr().DBReflect().SearchData()...)

	return rv
}

type ShopDB struct {
	tableName string
}

type ShopDBQueryBuilder struct {
	tableName string
}

type ShopDBQueryUnsafe struct {
	tableName string
}

type ShopDBColumns struct {
	tableName string
}

func (x *Shop) DB() *ShopDB {
	return &ShopDB{tableName: x.DBReflect().Descriptor().TableName()}
}

func (x *ShopDB) TableName() string {
	return x.tableName
}

func (x *ShopDB) Query() *ShopDBQueryBuilder {
	return &ShopDBQueryBuilder{tableName: x.tableName}
}

func (x *ShopDB) Columns() *ShopDBColumns {
	return &ShopDBColumns{tableName: x.tableName}
}

func (x *ShopDB) WithTable(t string) *ShopDB {
	return &ShopDB{tableName: t}
}

func (x *ShopDBQueryBuilder) WithTable(t string) *ShopDBQueryBuilder {
	return &ShopDBQueryBuilder{tableName: t}
}

func (x *ShopDBQueryBuilder) Unsafe() *ShopDBQueryUnsafe {
	return &ShopDBQueryUnsafe{tableName: x.tableName}
}

type ShopTenantIdSafeOperators struct {
	column    string
	tableName string
}

func (x *ShopTenantIdSafeOperators) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

func (x *ShopTenantIdSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Eq(v)
}

func (x *ShopTenantIdSafeOperators) Gt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gt(v)
}

func (x *ShopTenantIdSafeOperators) Gte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gte(v)
}

func (x *ShopTenantIdSafeOperators) Lt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Lt(v)
}

func (x *ShopTenantIdSafeOperators) Lte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Lte(v)
}

func (x *ShopTenantIdSafeOperators) In(v []string) exp.BooleanExpression {
	if len(v) == 0 {
		return exp.NewBooleanExpression(exp.EqOp, exp.NewLiteralExpression("FALSE"), true)
	}
	return exp.NewIdentifierExpression("", x.tableName, x.column).In(v)
}

func (x *ShopTenantIdSafeOperators) NotIn(v []string) exp.BooleanExpression {
	if len(v) == 0 {
		return exp.NewBooleanExpression(exp.EqOp, exp.NewLiteralExpression("TRUE"), true)
	}
	return exp.NewIdentifierExpression("", x.tableName, x.column).NotIn(v)
}

func (x *ShopTenantIdSafeOperators) IsNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).IsNull()
}

func (x *ShopTenantIdSafeOperators) IsNotEmpty() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gt("")
}

func (x *ShopTenantIdSafeOperators) IsNotNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).IsNotNull()
}

func (x *ShopTenantIdSafeOperators) Between(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Between(exp.NewRangeVal(start, end))
}

func (x *ShopTenantIdSafeOperators) NotBetween(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).NotBetween(exp.NewRangeVal(start, end))
}

func (x *ShopDBQueryBuilder) TenantId() *ShopTenantIdSafeOperators {
	return &ShopTenantIdSafeOperators{tableName: x.tableName, column: "pb$" + "tenant_id"}
}

type ShopPKSKSafeOperators struct {
	column    string
	tableName string
}

func (x *ShopPKSKSafeOperators) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

func (x *ShopPKSKSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Eq(v)
}

func (x *ShopPKSKSafeOperators) Gt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gt(v)
}

func (x *ShopPKSKSafeOperators) Gte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gte(v)
}

func (x *ShopPKSKSafeOperators) Lt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Lt(v)
}

func (x *ShopPKSKSafeOperators) Lte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Lte(v)
}

func (x *ShopPKSKSafeOperators) In(v []string) exp.BooleanExpression {
	if len(v) == 0 {
		return exp.NewBooleanExpression(exp.EqOp, exp.NewLiteralExpression("FALSE"), true)
	}
	return exp.NewIdentifierExpression("", x.tableName, x.column).In(v)
}

func (x *ShopPKSKSafeOperators) NotIn(v []string) exp.BooleanExpression {
	if len(v) == 0 {
		return exp.NewBooleanExpression(exp.EqOp, exp.NewLiteralExpression("TRUE"), true)
	}
	return exp.NewIdentifierExpression("", x.tableName, x.column).NotIn(v)
}

func (x *ShopPKSKSafeOperators) IsNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).IsNull()
}

func (x *ShopPKSKSafeOperators) IsNotEmpty() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gt("")
}

func (x *ShopPKSKSafeOperators) IsNotNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).IsNotNull()
}

func (x *ShopPKSKSafeOperators) Between(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Between(exp.NewRangeVal(start, end))
}

func (x *ShopPKSKSafeOperators) NotBetween(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).NotBetween(exp.NewRangeVal(start, end))
}

func (x *ShopDBQueryBuilder) PKSK() *ShopPKSKSafeOperators {
	return &ShopPKSKSafeOperators{tableName: x.tableName, column: "pb$" + "pksk"}
}

type ShopPKSafeOperators struct {
	column    string
	tableName string
}

func (x *ShopPKSafeOperators) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

func (x *ShopPKSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Eq(v)
}

func (x *ShopPKSafeOperators) Gt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gt(v)
}

func (x *ShopPKSafeOperators) Gte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gte(v)
}

func (x *ShopPKSafeOperators) Lt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Lt(v)
}

func (x *ShopPKSafeOperators) Lte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Lte(v)
}

func (x *ShopPKSafeOperators) In(v []string) exp.BooleanExpression {
	if len(v) == 0 {
		return exp.NewBooleanExpression(exp.EqOp, exp.NewLiteralExpression("FALSE"), true)
	}
	return exp.NewIdentifierExpression("", x.tableName, x.column).In(v)
}

func (x *ShopPKSafeOperators) NotIn(v []string) exp.BooleanExpression {
	if len(v) == 0 {
		return exp.NewBooleanExpression(exp.EqOp, exp.NewLiteralExpression("TRUE"), true)
	}
	return exp.NewIdentifierExpression("", x.tableName, x.column).NotIn(v)
}

func (x *ShopPKSafeOperators) IsNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).IsNull()
}

func (x *ShopPKSafeOperators) IsNotEmpty() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gt("")
}

func (x *ShopPKSafeOperators) IsNotNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).IsNotNull()
}

func (x *ShopPKSafeOperators) Between(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Between(exp.NewRangeVal(start, end))
}

func (x *ShopPKSafeOperators) NotBetween(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).NotBetween(exp.NewRangeVal(start, end))
}

func (x *ShopDBQueryBuilder) PK() *ShopPKSafeOperators {
	return &ShopPKSafeOperators{tableName: x.tableName, column: "pb$" + "pk"}
}

type ShopSKSafeOperators struct {
	column    string
	tableName string
}

func (x *ShopSKSafeOperators) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

func (x *ShopSKSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Eq(v)
}

func (x *ShopSKSafeOperators) Gt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gt(v)
}

func (x *ShopSKSafeOperators) Gte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gte(v)
}

func (x *ShopSKSafeOperators) Lt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Lt(v)
}

func (x *ShopSKSafeOperators) Lte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Lte(v)
}

func (x *ShopSKSafeOperators) In(v []string) exp.BooleanExpression {
	if len(v) == 0 {
		return exp.NewBooleanExpression(exp.EqOp, exp.NewLiteralExpression("FALSE"), true)
	}
	return exp.NewIdentifierExpression("", x.tableName, x.column).In(v)
}

func (x *ShopSKSafeOperators) NotIn(v []string) exp.BooleanExpression {
	if len(v) == 0 {
		return exp.NewBooleanExpression(exp.EqOp, exp.NewLiteralExpression("TRUE"), true)
	}
	return exp.NewIdentifierExpression("", x.tableName, x.column).NotIn(v)
}

func (x *ShopSKSafeOperators) IsNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).IsNull()
}

func (x *ShopSKSafeOperators) IsNotEmpty() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Gt("")
}

func (x *ShopSKSafeOperators) IsNotNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).IsNotNull()
}

func (x *ShopSKSafeOperators) Between(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Between(exp.NewRangeVal(start, end))
}

func (x *ShopSKSafeOperators) NotBetween(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).NotBetween(exp.NewRangeVal(start, end))
}

func (x *ShopDBQueryBuilder) SK() *ShopSKSafeOperators {
	return &ShopSKSafeOperators{tableName: x.tableName, column: "pb$" + "sk"}
}

type ShopFTSDataSafeOperators struct {
	column    string
	tableName string
}

func (x *ShopFTSDataSafeOperators) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

func (x *ShopFTSDataSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column).Eq(v)
}

func (x *ShopDBQueryBuilder) FTSData() *ShopFTSDataSafeOperators {
	return &ShopFTSDataSafeOperators{tableName: x.tableName, column: "pb$" + "fts_data"}
}

type ShopTenantIdQueryType struct {
	column    string
	tableName string
}

func (x *ShopDBQueryUnsafe) TenantId() *ShopTenantIdQueryType {
	return &ShopTenantIdQueryType{tableName: x.tableName, column: "pb$" + "tenant_id"}
}

func (x *ShopTenantIdQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

type ShopPKSKQueryType struct {
	column    string
	tableName string
}

func (x *ShopDBQueryUnsafe) PKSK() *ShopPKSKQueryType {
	return &ShopPKSKQueryType{tableName: x.tableName, column: "pb$" + "pksk"}
}

func (x *ShopPKSKQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

type ShopPKQueryType struct {
	column    string
	tableName string
}

func (x *ShopDBQueryUnsafe) PK() *ShopPKQueryType {
	return &ShopPKQueryType{tableName: x.tableName, column: "pb$" + "pk"}
}

func (x *ShopPKQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

type ShopSKQueryType struct {
	column    string
	tableName string
}

func (x *ShopDBQueryUnsafe) SK() *ShopSKQueryType {
	return &ShopSKQueryType{tableName: x.tableName, column: "pb$" + "sk"}
}

func (x *ShopSKQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

type ShopFTSDataQueryType struct {
	column    string
	tableName string
}

func (x *ShopDBQueryUnsafe) FTSData() *ShopFTSDataQueryType {
	return &ShopFTSDataQueryType{tableName: x.tableName, column: "pb$" + "fts_data"}
}

func (x *ShopFTSDataQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

type ShopPBDataQueryType struct {
	column    string
	tableName string
}

func (x *ShopDBQueryUnsafe) PBData() *ShopPBDataQueryType {
	return &ShopPBDataQueryType{tableName: x.tableName, column: "pb$" + "pb_data"}
}

func (x *ShopPBDataQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

type ShopIdQueryType struct {
	column    string
	tableName string
}

func (x *ShopDBQueryUnsafe) Id() *ShopIdQueryType {
	return &ShopIdQueryType{tableName: x.tableName, column: "pb$" + "id"}
}

func (x *ShopIdQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

type ShopCreatedAtQueryType struct {
	column    string
	tableName string
}

func (x *ShopDBQueryUnsafe) CreatedAt() *ShopCreatedAtQueryType {
	return &ShopCreatedAtQueryType{tableName: x.tableName, column: "pb$" + "created_at"}
}

func (x *ShopCreatedAtQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

type ShopFurQueryType struct {
	column    string
	tableName string
}

func (x *ShopDBQueryUnsafe) Fur() *ShopFurQueryType {
	return &ShopFurQueryType{tableName: x.tableName, column: "pb$" + "fur"}
}

func (x *ShopFurQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

type ShopMediumQueryType struct {
	column    string
	tableName string
}

func (x *ShopDBQueryUnsafe) Medium() *ShopMediumQueryType {
	return &ShopMediumQueryType{tableName: x.tableName, column: "pb$" + "medium_oneof"}
}

func (x *ShopMediumQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

func (x *ShopDBColumns) WithTable(t string) *ShopDBColumns {
	return &ShopDBColumns{tableName: t}
}

func (x *ShopDBColumns) TenantId() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "tenant_id")
}

func (x *ShopDBColumns) PKSK() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "pksk")
}

func (x *ShopDBColumns) PK() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "pk")
}

func (x *ShopDBColumns) SK() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "sk")
}

func (x *ShopDBColumns) FTSData() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "fts_data")
}

func (x *ShopDBColumns) PBData() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "pb_data")
}

func (x *ShopDBColumns) Id() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "id")
}

func (x *ShopDBColumns) CreatedAt() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "created_at")
}

func (x *ShopDBColumns) Fur() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "fur")
}

func (x *ShopDBColumns) Medium() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "medium_oneof")
}

type pgdbDescriptorShop_Manager struct{}

var (
	instancepgdbDescriptorShop_Manager pgdb_v1.Descriptor = &pgdbDescriptorShop_Manager{}
)

func (d *pgdbDescriptorShop_Manager) TableName() string {
	return "pb_manager_models_zoo_v1_6ccf2214"
}

func (d *pgdbDescriptorShop_Manager) IsPartitioned() bool {
	return false
}

func (d *pgdbDescriptorShop_Manager) Fields(opts ...pgdb_v1.DescriptorFieldOptionFunc) []*pgdb_v1.Column {
	df := pgdb_v1.NewDescriptorFieldOption(opts)
	_ = df

	rv := make([]*pgdb_v1.Column, 0)

	rv = append(rv, &pgdb_v1.Column{
		Name:               df.ColumnName("id"),
		Type:               "int4",
		Nullable:           df.Nullable(false),
		OverrideExpression: "",
		Default:            "0",
	})

	return rv
}

func (d *pgdbDescriptorShop_Manager) PKSKField() *pgdb_v1.Column {
	return &pgdb_v1.Column{
		Table: "pb_manager_models_zoo_v1_6ccf2214",
		Name:  "pb$pksk",
		Type:  "varchar",
	}
}

func (d *pgdbDescriptorShop_Manager) DataField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Table: "pb_manager_models_zoo_v1_6ccf2214", Name: "pb$pb_data", Type: "bytea"}
}

func (d *pgdbDescriptorShop_Manager) SearchField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Table: "pb_manager_models_zoo_v1_6ccf2214", Name: "pb$fts_data", Type: "tsvector"}
}

func (d *pgdbDescriptorShop_Manager) VersioningField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Table: "pb_manager_models_zoo_v1_6ccf2214", Name: "pb$", Type: "timestamptz"}
}

func (d *pgdbDescriptorShop_Manager) TenantField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Table: "pb_manager_models_zoo_v1_6ccf2214", Name: "pb$tenant_id", Type: "varchar"}
}

func (d *pgdbDescriptorShop_Manager) IndexPrimaryKey(opts ...pgdb_v1.IndexOptionsFunc) *pgdb_v1.Index {
	io := pgdb_v1.NewIndexOptions(opts)
	_ = io

	return nil

}

func (d *pgdbDescriptorShop_Manager) Indexes(opts ...pgdb_v1.IndexOptionsFunc) []*pgdb_v1.Index {
	io := pgdb_v1.NewIndexOptions(opts)
	_ = io
	rv := make([]*pgdb_v1.Index, 0)

	return rv
}

func (d *pgdbDescriptorShop_Manager) Statistics(opts ...pgdb_v1.StatisticOptionsFunc) []*pgdb_v1.Statistic {
	io := pgdb_v1.NewStatisticOption(opts)
	_ = io
	rv := make([]*pgdb_v1.Statistic, 0)

	return rv
}

type pgdbMessageShop_Manager struct {
	self *Shop_Manager
}

func (dbr *Shop_Manager) DBReflect() pgdb_v1.Message {
	return &pgdbMessageShop_Manager{
		self: dbr,
	}
}

func (m *pgdbMessageShop_Manager) Descriptor() pgdb_v1.Descriptor {
	return instancepgdbDescriptorShop_Manager
}

func (m *pgdbMessageShop_Manager) Record(opts ...pgdb_v1.RecordOptionsFunc) (exp.Record, error) {
	ro := pgdb_v1.NewRecordOptions(opts)
	_ = ro
	nullExp := exp.NewLiteralExpression("NULL")
	_ = nullExp

	rv := exp.Record{}

	v1 := int32(m.self.GetId())

	if ro.Nulled {
		rv[ro.ColumnName("id")] = nullExp
	} else {
		rv[ro.ColumnName("id")] = v1
	}

	return rv, nil
}

func (m *pgdbMessageShop_Manager) SearchData(opts ...pgdb_v1.RecordOptionsFunc) []*pgdb_v1.SearchContent {
	rv := []*pgdb_v1.SearchContent{}

	return rv
}

type Shop_ManagerDB struct {
	tableName string
}

type Shop_ManagerDBQueryBuilder struct {
	tableName string
}

type Shop_ManagerDBQueryUnsafe struct {
	tableName string
}

type Shop_ManagerDBColumns struct {
	tableName string
}

func (x *Shop_Manager) DB() *Shop_ManagerDB {
	return &Shop_ManagerDB{tableName: x.DBReflect().Descriptor().TableName()}
}

func (x *Shop_ManagerDB) TableName() string {
	return x.tableName
}

func (x *Shop_ManagerDB) Query() *Shop_ManagerDBQueryBuilder {
	return &Shop_ManagerDBQueryBuilder{tableName: x.tableName}
}

func (x *Shop_ManagerDB) Columns() *Shop_ManagerDBColumns {
	return &Shop_ManagerDBColumns{tableName: x.tableName}
}

func (x *Shop_ManagerDB) WithTable(t string) *Shop_ManagerDB {
	return &Shop_ManagerDB{tableName: t}
}

func (x *Shop_ManagerDBQueryBuilder) WithTable(t string) *Shop_ManagerDBQueryBuilder {
	return &Shop_ManagerDBQueryBuilder{tableName: t}
}

func (x *Shop_ManagerDBQueryBuilder) Unsafe() *Shop_ManagerDBQueryUnsafe {
	return &Shop_ManagerDBQueryUnsafe{tableName: x.tableName}
}

type Shop_ManagerIdQueryType struct {
	column    string
	tableName string
}

func (x *Shop_ManagerDBQueryUnsafe) Id() *Shop_ManagerIdQueryType {
	return &Shop_ManagerIdQueryType{tableName: x.tableName, column: "pb$" + "id"}
}

func (x *Shop_ManagerIdQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.column)
}

func (x *Shop_ManagerDBColumns) WithTable(t string) *Shop_ManagerDBColumns {
	return &Shop_ManagerDBColumns{tableName: t}
}

func (x *Shop_ManagerDBColumns) Id() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "id")
}
