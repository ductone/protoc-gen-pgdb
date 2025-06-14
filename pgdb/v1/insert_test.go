package v1

import (
	"testing"

	"github.com/doug-martin/goqu/v9/exp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ DBReflectMessage = (*mockDBReflect)(nil)

type mockDBReflect struct{}

func (m *mockDBReflect) DBReflect() Message {
	return m.DBReflectWithDialect(DialectUnspecified)
}

func (m *mockDBReflect) DBReflectWithDialect(dialect Dialect) Message {
	return &mockMessage{dialect: dialect}
}

var _ Message = (*mockMessage)(nil)

type mockMessage struct {
	dialect Dialect
}

func (m *mockMessage) Descriptor() Descriptor {
	return &mockDescriptor{
		tableName: "table_name",
	}
}

func (m *mockMessage) Record(opts ...RecordOptionsFunc) (exp.Record, error) {
	return exp.Record{
		"pb$tenant_id":  "tenant_id",
		"pb$pk":         "pk",
		"pb$sk":         "sk",
		"pb$updated_at": "1900-01-01T00:00:00Z",
	}, nil
}

func (m *mockMessage) SearchData(opts ...RecordOptionsFunc) []*SearchContent {
	return nil
}

func (m *mockMessage) Dialect() Dialect {
	return m.dialect
}

func TestInsert(t *testing.T) {
	t.Parallel()
	tt := []struct {
		Name         string
		DialectOpts  []DialectOpt
		ExpectedSQL  string
		ExpectedArgs []any
	}{
		{
			Name:        "empty opts",
			DialectOpts: []DialectOpt{},
			ExpectedSQL: `INSERT INTO "table_name" ("pb$pk", "pb$sk", "pb$tenant_id", "pb$updated_at") VALUES ($1, $2, $3, $4) ON CONFLICT ON CONSTRAINT "pbidx_table_name" DO UPDATE SET "pb$pk"="excluded"."pb$pk","pb$sk"="excluded"."pb$sk","pb$updated_at"="excluded"."pb$updated_at" WHERE ("table_name"."pb$updated_at" <= $5::timestamptz)`, //nolint:revive // query
			ExpectedArgs: []any{
				"pk",
				"sk",
				"tenant_id",
				"1900-01-01T00:00:00Z",
				"1900-01-01T00:00:00Z",
			},
		},
		{
			Name:        "unspecified dialect",
			DialectOpts: []DialectOpt{DialectUnspecified},
			ExpectedSQL: `INSERT INTO "table_name" ("pb$pk", "pb$sk", "pb$tenant_id", "pb$updated_at") VALUES ($1, $2, $3, $4) ON CONFLICT ON CONSTRAINT "pbidx_table_name" DO UPDATE SET "pb$pk"="excluded"."pb$pk","pb$sk"="excluded"."pb$sk","pb$updated_at"="excluded"."pb$updated_at" WHERE ("table_name"."pb$updated_at" <= $5::timestamptz)`, //nolint:revive // query
			ExpectedArgs: []any{
				"pk",
				"sk",
				"tenant_id",
				"1900-01-01T00:00:00Z",
				"1900-01-01T00:00:00Z",
			},
		},
		{
			Name:        "v13",
			DialectOpts: []DialectOpt{DialectV13},
			ExpectedSQL: `INSERT INTO "table_name" ("pb$pk", "pb$sk", "pb$tenant_id", "pb$updated_at") VALUES ($1, $2, $3, $4) ON CONFLICT ON CONSTRAINT "pbidx_table_name" DO UPDATE SET "pb$pk"="excluded"."pb$pk","pb$sk"="excluded"."pb$sk","pb$updated_at"="excluded"."pb$updated_at" WHERE ("table_name"."pb$updated_at" <= $5::timestamptz)`, //nolint:revive // query
			ExpectedArgs: []any{
				"pk",
				"sk",
				"tenant_id",
				"1900-01-01T00:00:00Z",
				"1900-01-01T00:00:00Z",
			},
		},
		{
			Name:        "v17",
			DialectOpts: []DialectOpt{DialectV17},
			ExpectedSQL: `INSERT INTO "table_name" ("pb$pk", "pb$pksk", "pb$sk", "pb$tenant_id", "pb$updated_at") VALUES ($1, $2, $3, $4, $5) ON CONFLICT ON CONSTRAINT "pbidx_table_name" DO UPDATE SET "pb$pk"="excluded"."pb$pk","pb$sk"="excluded"."pb$sk","pb$updated_at"="excluded"."pb$updated_at" WHERE ("table_name"."pb$updated_at" <= $6::timestamptz)`, //nolint:revive // query
			ExpectedArgs: []any{
				"pk",
				"pk|sk",
				"sk",
				"tenant_id",
				"1900-01-01T00:00:00Z",
				"1900-01-01T00:00:00Z",
			},
		},
		{
			Name:        "multiple opts passed",
			DialectOpts: []DialectOpt{DialectV13, DialectV17},
			ExpectedSQL: `INSERT INTO "table_name" ("pb$pk", "pb$sk", "pb$tenant_id", "pb$updated_at") VALUES ($1, $2, $3, $4) ON CONFLICT ON CONSTRAINT "pbidx_table_name" DO UPDATE SET "pb$pk"="excluded"."pb$pk","pb$sk"="excluded"."pb$sk","pb$updated_at"="excluded"."pb$updated_at" WHERE ("table_name"."pb$updated_at" <= $5::timestamptz)`, //nolint:revive // query
			ExpectedArgs: []any{
				"pk",
				"sk",
				"tenant_id",
				"1900-01-01T00:00:00Z",
				"1900-01-01T00:00:00Z",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			var actualSQL string
			var actualArgs []any
			var err error
			require.NotPanics(t, func() {
				actualSQL, actualArgs, err = Insert(&mockDBReflect{}, tc.DialectOpts...)
			})
			require.NoError(t, err)
			assert.Equal(t, tc.ExpectedSQL, actualSQL)
			assert.Equal(t, tc.ExpectedArgs, actualArgs)
		})
	}
}
