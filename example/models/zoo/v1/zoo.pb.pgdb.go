// Code generated by protoc-gen-pgdb 0.1.0 from models/zoo/v1/zoo.proto. DO NOT EDIT
package v1

import (
	"encoding/json"
	"strings"

	"time"

	v1 "github.com/ductone/protoc-gen-pgdb/example/models/animals/v1"

	"github.com/doug-martin/goqu/v9/exp"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	"github.com/ductone/protoc-gen-pgdb/pgdb/v1/xpq"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type pgdbDescriptorShop struct{}

var (
	instancepgdbDescriptorShop pgdb_v1.Descriptor = &pgdbDescriptorShop{}
)

func (d *pgdbDescriptorShop) TableName() string {
	return "pb_shop_models_zoo_v1_ca2425f6"
}

func (d *pgdbDescriptorShop) Fields(opts ...pgdb_v1.DescriptorFieldOptionFunc) []*pgdb_v1.Column {
	df := pgdb_v1.NewDescriptorFieldOption(opts)
	_ = df
	rv := []*pgdb_v1.Column{
		{
			Name:               df.ColumnName("tenant_id"),
			Type:               "varchar",
			Nullable:           false,
			OverrideExpression: "",
		}, {
			Name:               df.ColumnName("pksk"),
			Type:               "varchar",
			Nullable:           false,
			OverrideExpression: "varchar GENERATED ALWAYS AS (pb$pk || '|' || pb$sk) STORED",
		}, {
			Name:               df.ColumnName("pk"),
			Type:               "varchar",
			Nullable:           false,
			OverrideExpression: "",
		}, {
			Name:               df.ColumnName("sk"),
			Type:               "varchar",
			Nullable:           false,
			OverrideExpression: "",
		}, {
			Name:               df.ColumnName("fts_data"),
			Type:               "tsvector",
			Nullable:           true,
			OverrideExpression: "",
		}, {
			Name:               df.ColumnName("pb_data"),
			Type:               "bytea",
			Nullable:           false,
			OverrideExpression: "",
		}, {
			Name:               df.ColumnName("id"),
			Type:               "text",
			Nullable:           false,
			OverrideExpression: "",
		}, {
			Name:               df.ColumnName("created_at"),
			Type:               "timestamptz",
			Nullable:           true,
			OverrideExpression: "",
		}, {
			Name:               df.ColumnName("fur"),
			Type:               "int4",
			Nullable:           false,
			OverrideExpression: "",
		}, {
			Name:               df.ColumnName("medium_oneof"),
			Type:               "int4",
			Nullable:           false,
			OverrideExpression: "",
		},
	}

	rv = append(rv, ((*v1.PaperBook)(nil)).DBReflect().Descriptor().Fields(df.Nested("50$")...)...)

	rv = append(rv, ((*v1.EBook)(nil)).DBReflect().Descriptor().Fields(df.Nested("51$")...)...)

	rv = append(rv, ((*Shop_Manager)(nil)).DBReflect().Descriptor().Fields(df.Nested("5$")...)...)

	return rv
}

func (d *pgdbDescriptorShop) DataField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Name: "pb_data", Type: "bytea"}
}

func (d *pgdbDescriptorShop) SearchField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Name: "fts_data", Type: "tsvector"}
}

func (d *pgdbDescriptorShop) VersioningField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Name: "pb$created_at", Type: "timestamptz"}
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
	rv := []*pgdb_v1.Index{
		{
			Name:               io.IndexName("pksk_shop_models_zoo_v1_a0f13ce7"),
			Method:             pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE,
			IsPrimary:          true,
			IsUnique:           true,
			IsDropped:          false,
			Columns:            []string{io.ColumnName("tenant_id"), io.ColumnName("pksk")},
			OverrideExpression: "",
		}, {
			Name:               io.IndexName("pksk_shop_models_zoo_v1_a0f13ce7"),
			Method:             pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE,
			IsPrimary:          false,
			IsUnique:           true,
			IsDropped:          false,
			Columns:            []string{io.ColumnName("tenant_id"), io.ColumnName("pk"), io.ColumnName("sk")},
			OverrideExpression: "",
		}, {
			Name:               io.IndexName("fts_data_shop_models_zoo_v1_1b685f12"),
			Method:             pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE_GIN,
			IsPrimary:          false,
			IsUnique:           false,
			IsDropped:          false,
			Columns:            []string{io.ColumnName("tenant_id"), io.ColumnName("fts_data")},
			OverrideExpression: "",
		},
	}

	rv = append(rv, ((*v1.PaperBook)(nil)).DBReflect().Descriptor().Indexes(io.Nested("50$")...)...)

	rv = append(rv, ((*v1.EBook)(nil)).DBReflect().Descriptor().Indexes(io.Nested("51$")...)...)

	rv = append(rv, ((*Shop_Manager)(nil)).DBReflect().Descriptor().Indexes(io.Nested("5$")...)...)

	return rv
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

	var sb strings.Builder

	rv := exp.Record{}

	cfv0 := string(m.self.TenantId)

	rv[ro.ColumnName("tenant_id")] = cfv0

	sb.Reset()

	_, _ = sb.WriteString("models_zoo_v1_shop")

	_, _ = sb.WriteString(":")

	_, _ = sb.WriteString(m.self.TenantId)

	_, _ = sb.WriteString(":")

	_, _ = sb.WriteString(m.self.Id)

	cfv2 := sb.String()

	rv[ro.ColumnName("pk")] = cfv2

	sb.Reset()

	_, _ = sb.WriteString("example")

	cfv3 := sb.String()

	rv[ro.ColumnName("sk")] = cfv3

	cfv4tmp := []*pgdb_v1.SearchContent{

		{
			Type:   pgdb_v1.FieldOptions_FULL_TEXT_TYPE_EXACT,
			Weight: pgdb_v1.FieldOptions_FULL_TEXT_WEIGHT_UNSPECIFIED,
			Value:  m.self.Id,
		},
	}

	cfv4tmp = append(cfv4tmp, m.self.GetPaper().DBReflect().SearchData()...)

	cfv4tmp = append(cfv4tmp, m.self.GetEbook().DBReflect().SearchData()...)

	cfv4tmp = append(cfv4tmp, m.self.GetMgr().DBReflect().SearchData()...)

	cfv4 := pgdb_v1.FullTextSearchVectors(cfv4tmp)

	rv[ro.ColumnName("fts_data")] = cfv4

	cfv5, err := proto.Marshal(m.self)
	if err != nil {
		return nil, err
	}

	rv[ro.ColumnName("pb_data")] = cfv5

	v1 := string(m.self.GetId())

	rv[ro.ColumnName("id")] = v1

	var v2 *time.Time
	if m.self.GetCreatedAt().IsValid() {
		v2tmp := m.self.GetCreatedAt().AsTime()
		v2 = &v2tmp
	}

	rv[ro.ColumnName("created_at")] = v2

	v3, err := pgdb_v1.MarshalNestedRecord(m.self.GetPaper(), ro.Nested("50$")...)
	if err != nil {
		return nil, err
	}

	for k, v := range v3 {
		rv[k] = v
	}

	v4, err := pgdb_v1.MarshalNestedRecord(m.self.GetEbook(), ro.Nested("51$")...)
	if err != nil {
		return nil, err
	}

	for k, v := range v4 {
		rv[k] = v
	}

	v5 := int32(m.self.GetFur())

	rv[ro.ColumnName("fur")] = v5

	v6, err := pgdb_v1.MarshalNestedRecord(m.self.GetMgr(), ro.Nested("5$")...)
	if err != nil {
		return nil, err
	}

	for k, v := range v6 {
		rv[k] = v
	}

	oneof1 := uint32(0)

	switch m.self.GetMedium().(type) {

	case *Shop_Paper:
		oneof1 = 50

	case *Shop_Ebook:
		oneof1 = 51

	}

	rv[ro.ColumnName("medium_oneof")] = oneof1

	return rv, nil
}

func (m *pgdbMessageShop) SearchData(opts ...pgdb_v1.RecordOptionsFunc) []*pgdb_v1.SearchContent {
	rv := []*pgdb_v1.SearchContent{

		{
			Type:   pgdb_v1.FieldOptions_FULL_TEXT_TYPE_EXACT,
			Weight: pgdb_v1.FieldOptions_FULL_TEXT_WEIGHT_UNSPECIFIED,
			Value:  m.self.Id,
		},
	}

	rv = append(rv, m.self.GetPaper().DBReflect().SearchData()...)

	rv = append(rv, m.self.GetEbook().DBReflect().SearchData()...)

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
	prefix    string
	tableName string
}

func (x *ShopTenantIdSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").Eq(v)
}

func (x *ShopTenantIdSafeOperators) Neq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").Neq(v)
}

func (x *ShopTenantIdSafeOperators) Gt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").Gt(v)
}

func (x *ShopTenantIdSafeOperators) Gte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").Gte(v)
}

func (x *ShopTenantIdSafeOperators) Lt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").Lt(v)
}

func (x *ShopTenantIdSafeOperators) Lte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").Lte(v)
}

func (x *ShopTenantIdSafeOperators) In(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").In(v)
}

func (x *ShopTenantIdSafeOperators) NotIn(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").NotIn(v)
}

func (x *ShopTenantIdSafeOperators) IsNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").IsNull()
}

func (x *ShopTenantIdSafeOperators) IsNotNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").IsNotNull()
}

func (x *ShopTenantIdSafeOperators) Between(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").Between(exp.NewRangeVal(start, end))
}

func (x *ShopTenantIdSafeOperators) NotBetween(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").NotBetween(exp.NewRangeVal(start, end))
}

func (x *ShopDBQueryBuilder) TenantId() *ShopTenantIdSafeOperators {
	return &ShopTenantIdSafeOperators{tableName: x.tableName, prefix: "pb$"}
}

type ShopPKSKSafeOperators struct {
	prefix    string
	tableName string
}

func (x *ShopPKSKSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").Eq(v)
}

func (x *ShopPKSKSafeOperators) Neq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").Neq(v)
}

func (x *ShopPKSKSafeOperators) Gt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").Gt(v)
}

func (x *ShopPKSKSafeOperators) Gte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").Gte(v)
}

func (x *ShopPKSKSafeOperators) Lt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").Lt(v)
}

func (x *ShopPKSKSafeOperators) Lte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").Lte(v)
}

func (x *ShopPKSKSafeOperators) In(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").In(v)
}

func (x *ShopPKSKSafeOperators) NotIn(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").NotIn(v)
}

func (x *ShopPKSKSafeOperators) IsNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").IsNull()
}

func (x *ShopPKSKSafeOperators) IsNotNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").IsNotNull()
}

func (x *ShopPKSKSafeOperators) Between(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").Between(exp.NewRangeVal(start, end))
}

func (x *ShopPKSKSafeOperators) NotBetween(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").NotBetween(exp.NewRangeVal(start, end))
}

func (x *ShopDBQueryBuilder) PKSK() *ShopPKSKSafeOperators {
	return &ShopPKSKSafeOperators{tableName: x.tableName, prefix: "pb$"}
}

type ShopPKSafeOperators struct {
	prefix    string
	tableName string
}

func (x *ShopPKSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").Eq(v)
}

func (x *ShopPKSafeOperators) Neq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").Neq(v)
}

func (x *ShopPKSafeOperators) Gt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").Gt(v)
}

func (x *ShopPKSafeOperators) Gte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").Gte(v)
}

func (x *ShopPKSafeOperators) Lt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").Lt(v)
}

func (x *ShopPKSafeOperators) Lte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").Lte(v)
}

func (x *ShopPKSafeOperators) In(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").In(v)
}

func (x *ShopPKSafeOperators) NotIn(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").NotIn(v)
}

func (x *ShopPKSafeOperators) IsNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").IsNull()
}

func (x *ShopPKSafeOperators) IsNotNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").IsNotNull()
}

func (x *ShopPKSafeOperators) Between(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").Between(exp.NewRangeVal(start, end))
}

func (x *ShopPKSafeOperators) NotBetween(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").NotBetween(exp.NewRangeVal(start, end))
}

func (x *ShopDBQueryBuilder) PK() *ShopPKSafeOperators {
	return &ShopPKSafeOperators{tableName: x.tableName, prefix: "pb$"}
}

type ShopSKSafeOperators struct {
	prefix    string
	tableName string
}

func (x *ShopSKSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").Eq(v)
}

func (x *ShopSKSafeOperators) Neq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").Neq(v)
}

func (x *ShopSKSafeOperators) Gt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").Gt(v)
}

func (x *ShopSKSafeOperators) Gte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").Gte(v)
}

func (x *ShopSKSafeOperators) Lt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").Lt(v)
}

func (x *ShopSKSafeOperators) Lte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").Lte(v)
}

func (x *ShopSKSafeOperators) In(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").In(v)
}

func (x *ShopSKSafeOperators) NotIn(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").NotIn(v)
}

func (x *ShopSKSafeOperators) IsNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").IsNull()
}

func (x *ShopSKSafeOperators) IsNotNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").IsNotNull()
}

func (x *ShopSKSafeOperators) Between(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").Between(exp.NewRangeVal(start, end))
}

func (x *ShopSKSafeOperators) NotBetween(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").NotBetween(exp.NewRangeVal(start, end))
}

func (x *ShopDBQueryBuilder) SK() *ShopSKSafeOperators {
	return &ShopSKSafeOperators{tableName: x.tableName, prefix: "pb$"}
}

type ShopFTSDataSafeOperators struct {
	prefix    string
	tableName string
}

func (x *ShopFTSDataSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"fts_data").Eq(v)
}

func (x *ShopFTSDataSafeOperators) ObjectContains(obj interface{}) (exp.Expression, error) {
	var err error
	var data []byte

	pm, ok := obj.(proto.Message)
	if ok {
		data, err = protojson.Marshal(pm)
	} else {
		data, err = json.Marshal(obj)
	}
	if err != nil {
		return nil, err
	}

	idExp := exp.NewIdentifierExpression("", x.tableName, x.prefix+"fts_data")
	return exp.NewLiteralExpression("(? @> ?::jsonb)", idExp, string(data)), nil
}

func (x *ShopFTSDataSafeOperators) ObjectPathExists(path string) exp.Expression {
	idExp := exp.NewIdentifierExpression("", x.tableName, x.prefix+"fts_data")
	return exp.NewLiteralExpression("(? ? ?)", idExp, exp.NewLiteralExpression("@?"), path)
}

func (x *ShopFTSDataSafeOperators) ObjectPath(path string) exp.Expression {
	idExp := exp.NewIdentifierExpression("", x.tableName, x.prefix+"fts_data")
	return exp.NewLiteralExpression("? @@ ?", idExp, path)
}

func (x *ShopFTSDataSafeOperators) ObjectKeyExists(key string) exp.Expression {
	idExp := exp.NewIdentifierExpression("", x.tableName, x.prefix+"fts_data")
	return exp.NewLiteralExpression("? \\? ?", idExp, key)
}

func (x *ShopFTSDataSafeOperators) ObjectAnyKeyExists(keys ...string) exp.Expression {
	idExp := exp.NewIdentifierExpression("", x.tableName, x.prefix+"fts_data")
	return exp.NewLiteralExpression("(? ? ?)", idExp, exp.NewLiteralExpression("?|"), xpq.StringArray(keys))
}

func (x *ShopFTSDataSafeOperators) ObjectAllKeyExists(keys ...string) exp.Expression {
	idExp := exp.NewIdentifierExpression("", x.tableName, x.prefix+"fts_data")
	return exp.NewLiteralExpression("(? ? ?)", idExp, exp.NewLiteralExpression("?&"), xpq.StringArray(keys))
}

func (x *ShopDBQueryBuilder) FTSData() *ShopFTSDataSafeOperators {
	return &ShopFTSDataSafeOperators{tableName: x.tableName, prefix: "pb$"}
}

func (x *ShopDBQueryUnsafe) TenantId() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, "tenant_id")
}

func (x *ShopDBQueryUnsafe) PKSK() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, "pksk")
}

func (x *ShopDBQueryUnsafe) PK() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, "pk")
}

func (x *ShopDBQueryUnsafe) SK() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, "sk")
}

func (x *ShopDBQueryUnsafe) FTSData() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, "fts_data")
}

func (x *ShopDBQueryUnsafe) PBData() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, "pb_data")
}

func (x *ShopDBQueryUnsafe) Id() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, "id")
}

func (x *ShopDBQueryUnsafe) CreatedAt() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, "created_at")
}

func (x *ShopDBQueryUnsafe) Fur() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, "fur")
}

func (x *ShopDBQueryUnsafe) Medium() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, "medium_oneof")
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

func (d *pgdbDescriptorShop_Manager) Fields(opts ...pgdb_v1.DescriptorFieldOptionFunc) []*pgdb_v1.Column {
	df := pgdb_v1.NewDescriptorFieldOption(opts)
	_ = df
	rv := []*pgdb_v1.Column{
		{
			Name:               df.ColumnName("id"),
			Type:               "int4",
			Nullable:           false,
			OverrideExpression: "",
		},
	}

	return rv
}

func (d *pgdbDescriptorShop_Manager) DataField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Name: "pb_data", Type: "bytea"}
}

func (d *pgdbDescriptorShop_Manager) SearchField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Name: "fts_data", Type: "tsvector"}
}

func (d *pgdbDescriptorShop_Manager) VersioningField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Name: "pb$", Type: "timestamptz"}
}

func (d *pgdbDescriptorShop_Manager) IndexPrimaryKey(opts ...pgdb_v1.IndexOptionsFunc) *pgdb_v1.Index {
	io := pgdb_v1.NewIndexOptions(opts)
	_ = io

	return nil

}

func (d *pgdbDescriptorShop_Manager) Indexes(opts ...pgdb_v1.IndexOptionsFunc) []*pgdb_v1.Index {
	io := pgdb_v1.NewIndexOptions(opts)
	_ = io
	rv := []*pgdb_v1.Index{}

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

	rv := exp.Record{}

	v1 := int32(m.self.GetId())

	rv[ro.ColumnName("id")] = v1

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

func (x *Shop_ManagerDBQueryUnsafe) Id() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, "id")
}

func (x *Shop_ManagerDBColumns) WithTable(t string) *Shop_ManagerDBColumns {
	return &Shop_ManagerDBColumns{tableName: t}
}

func (x *Shop_ManagerDBColumns) Id() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "id")
}
