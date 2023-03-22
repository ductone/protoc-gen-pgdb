// Code generated by protoc-gen-pgdb 0.1.0 from models/city/v1/city.proto. DO NOT EDIT
package v1

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	zoo_v1 "github.com/ductone/protoc-gen-pgdb/example/models/zoo/v1"

	animals_v1 "github.com/ductone/protoc-gen-pgdb/example/models/animals/v1"

	"github.com/doug-martin/goqu/v9/exp"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	"github.com/ductone/protoc-gen-pgdb/pgdb/v1/xpq"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type pgdbDescriptorAttractions struct{}

var (
	instancepgdbDescriptorAttractions pgdb_v1.Descriptor = &pgdbDescriptorAttractions{}
)

func (d *pgdbDescriptorAttractions) TableName() string {
	return "pb_attractions_models_city_v1_e136cbfc"
}

func (d *pgdbDescriptorAttractions) Fields(opts ...pgdb_v1.DescriptorFieldOptionFunc) []*pgdb_v1.Column {
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
		Name:               df.ColumnName("numid"),
		Type:               "int4",
		Nullable:           df.Nullable(false),
		OverrideExpression: "",
		Default:            "0",
	})

	rv = append(rv, &pgdb_v1.Column{
		Name:               df.ColumnName("created_at"),
		Type:               "timestamptz",
		Nullable:           df.Nullable(true),
		OverrideExpression: "",
		Default:            "",
	})

	rv = append(rv, &pgdb_v1.Column{
		Name:               df.ColumnName("what_oneof"),
		Type:               "int4",
		Nullable:           df.Nullable(false),
		OverrideExpression: "",
		Default:            "",
	})

	rv = append(rv, ((*animals_v1.Pet)(nil)).DBReflect().Descriptor().Fields(df.Nested("10$")...)...)

	rv = append(rv, ((*zoo_v1.Shop)(nil)).DBReflect().Descriptor().Fields(df.Nested("11$")...)...)

	return rv
}

func (d *pgdbDescriptorAttractions) DataField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Table: "pb_attractions_models_city_v1_e136cbfc", Name: "pb_data", Type: "bytea"}
}

func (d *pgdbDescriptorAttractions) SearchField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Table: "pb_attractions_models_city_v1_e136cbfc", Name: "fts_data", Type: "tsvector"}
}

func (d *pgdbDescriptorAttractions) VersioningField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Table: "pb_attractions_models_city_v1_e136cbfc", Name: "pb$created_at", Type: "timestamptz"}
}

func (d *pgdbDescriptorAttractions) TenantField() *pgdb_v1.Column {
	return &pgdb_v1.Column{Table: "pb_attractions_models_city_v1_e136cbfc", Name: "pb$tenant_id", Type: "varchar"}
}

func (d *pgdbDescriptorAttractions) IndexPrimaryKey(opts ...pgdb_v1.IndexOptionsFunc) *pgdb_v1.Index {
	io := pgdb_v1.NewIndexOptions(opts)
	_ = io

	return &pgdb_v1.Index{
		Name:               io.IndexName("pksk_attractions_models_city_v1_1330fc81"),
		Method:             pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE,
		IsPrimary:          true,
		IsUnique:           true,
		IsDropped:          false,
		Columns:            []string{io.ColumnName("tenant_id"), io.ColumnName("pksk")},
		OverrideExpression: "",
	}

}

func (d *pgdbDescriptorAttractions) Indexes(opts ...pgdb_v1.IndexOptionsFunc) []*pgdb_v1.Index {
	io := pgdb_v1.NewIndexOptions(opts)
	_ = io
	rv := make([]*pgdb_v1.Index, 0)

	if !io.IsNested {

		rv = append(rv, &pgdb_v1.Index{
			Name:               io.IndexName("pksk_attractions_models_city_v1_1330fc81"),
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
			Name:               io.IndexName("pksk_attractions_models_city_v1_1330fc81"),
			Method:             pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE,
			IsPrimary:          false,
			IsUnique:           true,
			IsDropped:          false,
			Columns:            []string{io.ColumnName("tenant_id"), io.ColumnName("pk"), io.ColumnName("sk")},
			OverrideExpression: "",
		})

	}

	if !io.IsNested {

		rv = append(rv, &pgdb_v1.Index{
			Name:               io.IndexName("fts_data_attractions_models_city_v1_9239a529"),
			Method:             pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE_GIN,
			IsPrimary:          false,
			IsUnique:           false,
			IsDropped:          false,
			Columns:            []string{io.ColumnName("tenant_id"), io.ColumnName("fts_data")},
			OverrideExpression: "",
		})

	}

	rv = append(rv, ((*animals_v1.Pet)(nil)).DBReflect().Descriptor().Indexes(io.Nested("10$")...)...)

	rv = append(rv, ((*zoo_v1.Shop)(nil)).DBReflect().Descriptor().Indexes(io.Nested("11$")...)...)

	return rv
}

type pgdbMessageAttractions struct {
	self *Attractions
}

func (dbr *Attractions) DBReflect() pgdb_v1.Message {
	return &pgdbMessageAttractions{
		self: dbr,
	}
}

func (m *pgdbMessageAttractions) Descriptor() pgdb_v1.Descriptor {
	return instancepgdbDescriptorAttractions
}

func (m *pgdbMessageAttractions) Record(opts ...pgdb_v1.RecordOptionsFunc) (exp.Record, error) {
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

		_, _ = sb.WriteString("models_city_v1_attractions")

		_, _ = sb.WriteString(":")

		_, _ = sb.WriteString(m.self.TenantId)

		cfv2 := sb.String()

		if ro.Nulled {
			rv[ro.ColumnName("pk")] = nullExp
		} else {
			rv[ro.ColumnName("pk")] = cfv2
		}

	}

	if !ro.IsNested {

		sb.Reset()

		_, _ = sb.WriteString(m.self.Id)

		_, _ = sb.WriteString(":")

		_, _ = sb.WriteString(strconv.FormatInt(int64(m.self.Numid), 10))

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

		cfv4tmp = append(cfv4tmp, m.self.GetPet().DBReflect().SearchData()...)

		cfv4tmp = append(cfv4tmp, m.self.GetZooShop().DBReflect().SearchData()...)

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

	v2 := int32(m.self.GetNumid())

	if ro.Nulled {
		rv[ro.ColumnName("numid")] = nullExp
	} else {
		rv[ro.ColumnName("numid")] = v2
	}

	var v3 *time.Time
	if m.self.GetCreatedAt().IsValid() {
		v3tmp := m.self.GetCreatedAt().AsTime()
		v3 = &v3tmp
	}

	if ro.Nulled {
		rv[ro.ColumnName("created_at")] = nullExp
	} else {
		rv[ro.ColumnName("created_at")] = v3
	}

	v4tmp := m.self.GetPet()
	v4opts := ro.Nested("10$")
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

	v5tmp := m.self.GetZooShop()
	v5opts := ro.Nested("11$")
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

	oneof1 := uint32(0)

	switch m.self.GetWhat().(type) {

	case *Attractions_Pet:
		oneof1 = 10

	case *Attractions_ZooShop:
		oneof1 = 11

	}

	if ro.Nulled {
		rv[ro.ColumnName("what_oneof")] = nullExp
	} else {
		rv[ro.ColumnName("what_oneof")] = oneof1
	}

	return rv, nil
}

func (m *pgdbMessageAttractions) SearchData(opts ...pgdb_v1.RecordOptionsFunc) []*pgdb_v1.SearchContent {
	rv := []*pgdb_v1.SearchContent{

		{
			Type:   pgdb_v1.FieldOptions_FULL_TEXT_TYPE_EXACT,
			Weight: pgdb_v1.FieldOptions_FULL_TEXT_WEIGHT_UNSPECIFIED,
			Value:  m.self.GetId(),
		},
	}

	rv = append(rv, m.self.GetPet().DBReflect().SearchData()...)

	rv = append(rv, m.self.GetZooShop().DBReflect().SearchData()...)

	return rv
}

type AttractionsDB struct {
	tableName string
}

type AttractionsDBQueryBuilder struct {
	tableName string
}

type AttractionsDBQueryUnsafe struct {
	tableName string
}

type AttractionsDBColumns struct {
	tableName string
}

func (x *Attractions) DB() *AttractionsDB {
	return &AttractionsDB{tableName: x.DBReflect().Descriptor().TableName()}
}

func (x *AttractionsDB) TableName() string {
	return x.tableName
}

func (x *AttractionsDB) Query() *AttractionsDBQueryBuilder {
	return &AttractionsDBQueryBuilder{tableName: x.tableName}
}

func (x *AttractionsDB) Columns() *AttractionsDBColumns {
	return &AttractionsDBColumns{tableName: x.tableName}
}

func (x *AttractionsDB) WithTable(t string) *AttractionsDB {
	return &AttractionsDB{tableName: t}
}

func (x *AttractionsDBQueryBuilder) WithTable(t string) *AttractionsDBQueryBuilder {
	return &AttractionsDBQueryBuilder{tableName: t}
}

func (x *AttractionsDBQueryBuilder) Unsafe() *AttractionsDBQueryUnsafe {
	return &AttractionsDBQueryUnsafe{tableName: x.tableName}
}

type AttractionsTenantIdSafeOperators struct {
	prefix    string
	tableName string
}

func (x *AttractionsTenantIdSafeOperators) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id")
}

func (x *AttractionsTenantIdSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").Eq(v)
}

func (x *AttractionsTenantIdSafeOperators) Neq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").Neq(v)
}

func (x *AttractionsTenantIdSafeOperators) Gt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").Gt(v)
}

func (x *AttractionsTenantIdSafeOperators) Gte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").Gte(v)
}

func (x *AttractionsTenantIdSafeOperators) Lt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").Lt(v)
}

func (x *AttractionsTenantIdSafeOperators) Lte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").Lte(v)
}

func (x *AttractionsTenantIdSafeOperators) In(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").In(v)
}

func (x *AttractionsTenantIdSafeOperators) NotIn(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").NotIn(v)
}

func (x *AttractionsTenantIdSafeOperators) IsNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").IsNull()
}

func (x *AttractionsTenantIdSafeOperators) IsNotNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").IsNotNull()
}

func (x *AttractionsTenantIdSafeOperators) Between(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").Between(exp.NewRangeVal(start, end))
}

func (x *AttractionsTenantIdSafeOperators) NotBetween(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id").NotBetween(exp.NewRangeVal(start, end))
}

func (x *AttractionsDBQueryBuilder) TenantId() *AttractionsTenantIdSafeOperators {
	return &AttractionsTenantIdSafeOperators{tableName: x.tableName, prefix: "pb$"}
}

type AttractionsPKSKSafeOperators struct {
	prefix    string
	tableName string
}

func (x *AttractionsPKSKSafeOperators) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk")
}

func (x *AttractionsPKSKSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").Eq(v)
}

func (x *AttractionsPKSKSafeOperators) Neq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").Neq(v)
}

func (x *AttractionsPKSKSafeOperators) Gt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").Gt(v)
}

func (x *AttractionsPKSKSafeOperators) Gte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").Gte(v)
}

func (x *AttractionsPKSKSafeOperators) Lt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").Lt(v)
}

func (x *AttractionsPKSKSafeOperators) Lte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").Lte(v)
}

func (x *AttractionsPKSKSafeOperators) In(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").In(v)
}

func (x *AttractionsPKSKSafeOperators) NotIn(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").NotIn(v)
}

func (x *AttractionsPKSKSafeOperators) IsNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").IsNull()
}

func (x *AttractionsPKSKSafeOperators) IsNotNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").IsNotNull()
}

func (x *AttractionsPKSKSafeOperators) Between(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").Between(exp.NewRangeVal(start, end))
}

func (x *AttractionsPKSKSafeOperators) NotBetween(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk").NotBetween(exp.NewRangeVal(start, end))
}

func (x *AttractionsDBQueryBuilder) PKSK() *AttractionsPKSKSafeOperators {
	return &AttractionsPKSKSafeOperators{tableName: x.tableName, prefix: "pb$"}
}

type AttractionsPKSafeOperators struct {
	prefix    string
	tableName string
}

func (x *AttractionsPKSafeOperators) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk")
}

func (x *AttractionsPKSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").Eq(v)
}

func (x *AttractionsPKSafeOperators) Neq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").Neq(v)
}

func (x *AttractionsPKSafeOperators) Gt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").Gt(v)
}

func (x *AttractionsPKSafeOperators) Gte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").Gte(v)
}

func (x *AttractionsPKSafeOperators) Lt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").Lt(v)
}

func (x *AttractionsPKSafeOperators) Lte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").Lte(v)
}

func (x *AttractionsPKSafeOperators) In(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").In(v)
}

func (x *AttractionsPKSafeOperators) NotIn(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").NotIn(v)
}

func (x *AttractionsPKSafeOperators) IsNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").IsNull()
}

func (x *AttractionsPKSafeOperators) IsNotNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").IsNotNull()
}

func (x *AttractionsPKSafeOperators) Between(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").Between(exp.NewRangeVal(start, end))
}

func (x *AttractionsPKSafeOperators) NotBetween(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk").NotBetween(exp.NewRangeVal(start, end))
}

func (x *AttractionsDBQueryBuilder) PK() *AttractionsPKSafeOperators {
	return &AttractionsPKSafeOperators{tableName: x.tableName, prefix: "pb$"}
}

type AttractionsSKSafeOperators struct {
	prefix    string
	tableName string
}

func (x *AttractionsSKSafeOperators) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk")
}

func (x *AttractionsSKSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").Eq(v)
}

func (x *AttractionsSKSafeOperators) Neq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").Neq(v)
}

func (x *AttractionsSKSafeOperators) Gt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").Gt(v)
}

func (x *AttractionsSKSafeOperators) Gte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").Gte(v)
}

func (x *AttractionsSKSafeOperators) Lt(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").Lt(v)
}

func (x *AttractionsSKSafeOperators) Lte(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").Lte(v)
}

func (x *AttractionsSKSafeOperators) In(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").In(v)
}

func (x *AttractionsSKSafeOperators) NotIn(v []string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").NotIn(v)
}

func (x *AttractionsSKSafeOperators) IsNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").IsNull()
}

func (x *AttractionsSKSafeOperators) IsNotNull() exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").IsNotNull()
}

func (x *AttractionsSKSafeOperators) Between(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").Between(exp.NewRangeVal(start, end))
}

func (x *AttractionsSKSafeOperators) NotBetween(start string, end string) exp.RangeExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk").NotBetween(exp.NewRangeVal(start, end))
}

func (x *AttractionsDBQueryBuilder) SK() *AttractionsSKSafeOperators {
	return &AttractionsSKSafeOperators{tableName: x.tableName, prefix: "pb$"}
}

type AttractionsFTSDataSafeOperators struct {
	prefix    string
	tableName string
}

func (x *AttractionsFTSDataSafeOperators) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"fts_data")
}

func (x *AttractionsFTSDataSafeOperators) Eq(v string) exp.BooleanExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"fts_data").Eq(v)
}

func (x *AttractionsFTSDataSafeOperators) ObjectContains(obj interface{}) (exp.Expression, error) {
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

func (x *AttractionsFTSDataSafeOperators) ObjectPathExists(path string) exp.Expression {
	idExp := exp.NewIdentifierExpression("", x.tableName, x.prefix+"fts_data")
	return exp.NewLiteralExpression("(? ? ?)", idExp, exp.NewLiteralExpression("@?"), path)
}

func (x *AttractionsFTSDataSafeOperators) ObjectPath(path string) exp.Expression {
	idExp := exp.NewIdentifierExpression("", x.tableName, x.prefix+"fts_data")
	return exp.NewLiteralExpression("? @@ ?", idExp, path)
}

func (x *AttractionsFTSDataSafeOperators) ObjectKeyExists(key string) exp.Expression {
	idExp := exp.NewIdentifierExpression("", x.tableName, x.prefix+"fts_data")
	return exp.NewLiteralExpression("? \\? ?", idExp, key)
}

func (x *AttractionsFTSDataSafeOperators) ObjectAnyKeyExists(keys ...string) exp.Expression {
	idExp := exp.NewIdentifierExpression("", x.tableName, x.prefix+"fts_data")
	return exp.NewLiteralExpression("(? ? ?)", idExp, exp.NewLiteralExpression("?|"), xpq.StringArray(keys))
}

func (x *AttractionsFTSDataSafeOperators) ObjectAllKeyExists(keys ...string) exp.Expression {
	idExp := exp.NewIdentifierExpression("", x.tableName, x.prefix+"fts_data")
	return exp.NewLiteralExpression("(? ? ?)", idExp, exp.NewLiteralExpression("?&"), xpq.StringArray(keys))
}

func (x *AttractionsDBQueryBuilder) FTSData() *AttractionsFTSDataSafeOperators {
	return &AttractionsFTSDataSafeOperators{tableName: x.tableName, prefix: "pb$"}
}

type AttractionsTenantIdQueryType struct {
	prefix    string
	tableName string
}

func (x *AttractionsDBQueryUnsafe) TenantId() *AttractionsTenantIdQueryType {
	return &AttractionsTenantIdQueryType{tableName: x.tableName, prefix: "pb$"}
}

func (x *AttractionsTenantIdQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"tenant_id")
}

type AttractionsPKSKQueryType struct {
	prefix    string
	tableName string
}

func (x *AttractionsDBQueryUnsafe) PKSK() *AttractionsPKSKQueryType {
	return &AttractionsPKSKQueryType{tableName: x.tableName, prefix: "pb$"}
}

func (x *AttractionsPKSKQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pksk")
}

type AttractionsPKQueryType struct {
	prefix    string
	tableName string
}

func (x *AttractionsDBQueryUnsafe) PK() *AttractionsPKQueryType {
	return &AttractionsPKQueryType{tableName: x.tableName, prefix: "pb$"}
}

func (x *AttractionsPKQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pk")
}

type AttractionsSKQueryType struct {
	prefix    string
	tableName string
}

func (x *AttractionsDBQueryUnsafe) SK() *AttractionsSKQueryType {
	return &AttractionsSKQueryType{tableName: x.tableName, prefix: "pb$"}
}

func (x *AttractionsSKQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"sk")
}

type AttractionsFTSDataQueryType struct {
	prefix    string
	tableName string
}

func (x *AttractionsDBQueryUnsafe) FTSData() *AttractionsFTSDataQueryType {
	return &AttractionsFTSDataQueryType{tableName: x.tableName, prefix: "pb$"}
}

func (x *AttractionsFTSDataQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"fts_data")
}

type AttractionsPBDataQueryType struct {
	prefix    string
	tableName string
}

func (x *AttractionsDBQueryUnsafe) PBData() *AttractionsPBDataQueryType {
	return &AttractionsPBDataQueryType{tableName: x.tableName, prefix: "pb$"}
}

func (x *AttractionsPBDataQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"pb_data")
}

type AttractionsIdQueryType struct {
	prefix    string
	tableName string
}

func (x *AttractionsDBQueryUnsafe) Id() *AttractionsIdQueryType {
	return &AttractionsIdQueryType{tableName: x.tableName, prefix: "pb$"}
}

func (x *AttractionsIdQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"id")
}

type AttractionsNumidQueryType struct {
	prefix    string
	tableName string
}

func (x *AttractionsDBQueryUnsafe) Numid() *AttractionsNumidQueryType {
	return &AttractionsNumidQueryType{tableName: x.tableName, prefix: "pb$"}
}

func (x *AttractionsNumidQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"numid")
}

type AttractionsCreatedAtQueryType struct {
	prefix    string
	tableName string
}

func (x *AttractionsDBQueryUnsafe) CreatedAt() *AttractionsCreatedAtQueryType {
	return &AttractionsCreatedAtQueryType{tableName: x.tableName, prefix: "pb$"}
}

func (x *AttractionsCreatedAtQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"created_at")
}

type AttractionsWhatQueryType struct {
	prefix    string
	tableName string
}

func (x *AttractionsDBQueryUnsafe) What() *AttractionsWhatQueryType {
	return &AttractionsWhatQueryType{tableName: x.tableName, prefix: "pb$"}
}

func (x *AttractionsWhatQueryType) Identifier() exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", x.tableName, x.prefix+"what_oneof")
}

func (x *AttractionsDBColumns) WithTable(t string) *AttractionsDBColumns {
	return &AttractionsDBColumns{tableName: t}
}

func (x *AttractionsDBColumns) TenantId() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "tenant_id")
}

func (x *AttractionsDBColumns) PKSK() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "pksk")
}

func (x *AttractionsDBColumns) PK() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "pk")
}

func (x *AttractionsDBColumns) SK() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "sk")
}

func (x *AttractionsDBColumns) FTSData() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "fts_data")
}

func (x *AttractionsDBColumns) PBData() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "pb_data")
}

func (x *AttractionsDBColumns) Id() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "id")
}

func (x *AttractionsDBColumns) Numid() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "numid")
}

func (x *AttractionsDBColumns) CreatedAt() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "created_at")
}

func (x *AttractionsDBColumns) What() exp.Expression {
	return exp.NewIdentifierExpression("", x.tableName, "what_oneof")
}
