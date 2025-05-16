package v1

import (
	"testing"

	"github.com/doug-martin/goqu/v9/exp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ DBReflectMessage = (*mockDBReflect)(nil)

type mockDBReflect struct {
	msg *mockMessage
}

func (m mockDBReflect) DBReflect() Message {
	return m.msg
}

var _ Message = (*mockMessage)(nil)

type mockMessage struct {
	desc   *mockInsertDescriptor
	record exp.Record
}

func (m *mockMessage) Descriptor() Descriptor {
	return m.desc
}

func (m *mockMessage) Record(opts ...RecordOptionsFunc) (exp.Record, error) {
	return m.record, nil
}

func (m *mockMessage) SearchData(opts ...RecordOptionsFunc) []*SearchContent {
	return nil
}

type mockInsertDescriptor struct {
	tableName string
}

func (m *mockInsertDescriptor) TableName() string {
	return m.tableName
}

func (m *mockInsertDescriptor) Fields(opts ...DescriptorFieldOptionFunc) []*Column {
	df := NewDescriptorFieldOption(opts)
	_ = df
	rv := make([]*Column, 0)
	return rv
}

func (m *mockInsertDescriptor) PKSKField() *Column {
	return &Column{
		Table: m.tableName,
		Name:  "pb$pksk",
		Type:  "varchar",
	}
}

func (m *mockInsertDescriptor) PKSKV2Field() *Column {
	return &Column{
		Table:     m.tableName,
		Name:      "pb$pkskv2",
		Type:      "varchar",
		Nullable:  true,
		Collation: "C",
	}
}

func (m *mockInsertDescriptor) DataField() *Column {
	return &Column{
		Table: m.tableName,
		Name:  "pb$pb_data",
		Type:  "bytea",
	}
}

func (m *mockInsertDescriptor) SearchField() *Column {
	return &Column{
		Table: m.tableName,
		Name:  "pb$fts_data",
		Type:  "tsvector",
	}
}

func (m *mockInsertDescriptor) VersioningField() *Column {
	return &Column{
		Table: m.tableName,
		Name:  "pb$",
		Type:  "timestamptz",
	}
}

func (m *mockInsertDescriptor) TenantField() *Column {
	return &Column{
		Table: m.tableName,
		Name:  "pb$tenant_id",
		Type:  "varchar",
	}
}

func (m *mockInsertDescriptor) IsPartitioned() bool {
	return false
}

func (m *mockInsertDescriptor) IsPartitionedByCreatedAt() bool {
	return false
}

func (m *mockInsertDescriptor) GetPartitionedByKsuidFieldName() string {
	return ""
}

func (m *mockInsertDescriptor) Indexes(opts ...IndexOptionsFunc) []*Index {
	io := NewIndexOptions(opts)
	_ = io
	rv := make([]*Index, 0)

	if !io.IsNested {
		rv = append(rv, &Index{
			Name:               io.IndexName("pksk_" + m.tableName),
			Method:             MessageOptions_Index_INDEX_METHOD_BTREE,
			IsPrimary:          true,
			IsUnique:           true,
			IsDropped:          false,
			Columns:            []string{io.ColumnName("tenant_id"), io.ColumnName("pksk")},
			OverrideExpression: "",
		})
	}

	return rv
}

func (m *mockInsertDescriptor) IndexPrimaryKey(opts ...IndexOptionsFunc) *Index {
	io := NewIndexOptions(opts)
	_ = io
	return &Index{
		Name:               io.IndexName("pksk_" + m.tableName),
		Method:             MessageOptions_Index_INDEX_METHOD_BTREE,
		IsPrimary:          true,
		IsUnique:           true,
		IsDropped:          false,
		Columns:            []string{io.ColumnName("tenant_id"), io.ColumnName("pksk")},
		OverrideExpression: "",
	}
}

func (m *mockInsertDescriptor) Statistics(opts ...StatisticOptionsFunc) []*Statistic {
	io := NewStatisticOption(opts)
	_ = io
	rv := make([]*Statistic, 0)

	return rv
}

func (m *mockInsertDescriptor) GetPartitionDateRange() MessageOptions_PartitionedByDateRange {
	return MessageOptions_PARTITIONED_BY_DATE_RANGE_UNSPECIFIED
}

func TestBackfillPKSKV2(t *testing.T) {
	t.Parallel()

	dbr := &mockDBReflect{
		msg: &mockMessage{
			desc: &mockInsertDescriptor{
				tableName: "test_table",
			},
			record: exp.Record{
				"pb$pk":        "partition",
				"pb$sk":        "sort",
				"pb$pksk":      "partition|sort",
				"pb$pb_data":   []byte(`{"key": "value"}`),
				"pb$tenant_id": "tenant",
				"pb$":          "not a timestamp but whatever",
			},
		},
	}

	sql, params, err := BackfillPKSKV2(dbr)
	require.NoError(t, err)
	assert.Equal(t,
		`INSERT INTO "test_table" ("pb$", "pb$pb_data", "pb$pk", "pb$pksk", "pb$pkskv2", "pb$sk", "pb$tenant_id") VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT ON CONSTRAINT "pbidx_pksk_test_table" DO UPDATE SET "pb$pkskv2"="excluded"."pb$pksk" WHERE ("test_table"."pb$pkskv2" IS NULL)`, //nolint:revive // this is just a long lint
		sql,
	)
	require.Len(t, params, 7)
}
