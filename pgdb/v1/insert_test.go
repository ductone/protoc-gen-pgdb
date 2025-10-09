package v1

import (
	"testing"

	"github.com/doug-martin/goqu/v9/exp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ DBReflectMessage = (*mockDBReflect)(nil)

type mockDBReflect struct {
	tenantID  string
	pk        string
	sk        string
	updatedAt string
}

func (m *mockDBReflect) DBReflect(dialect Dialect) Message {
	return &mockMessage{
		dialect:   dialect,
		tenantID:  m.tenantID,
		pk:        m.pk,
		sk:        m.sk,
		updatedAt: m.updatedAt,
	}
}

var _ Message = (*mockMessage)(nil)

type mockMessage struct {
	dialect   Dialect
	tenantID  string
	pk        string
	sk        string
	updatedAt string
}

func (m *mockMessage) Descriptor() Descriptor {
	return &mockDescriptor{
		tableName: "table_name",
	}
}

func (m *mockMessage) Record(opts ...RecordOptionsFunc) (exp.Record, error) {
	tenantID := m.tenantID
	if tenantID == "" {
		tenantID = "tenant_id"
	}
	pk := m.pk
	if pk == "" {
		pk = "pk"
	}
	sk := m.sk
	if sk == "" {
		sk = "sk"
	}
	updatedAt := m.updatedAt
	if updatedAt == "" {
		updatedAt = "1900-01-01T00:00:00Z"
	}
	return exp.Record{
		"pb$tenant_id":  tenantID,
		"pb$pk":         pk,
		"pb$sk":         sk,
		"pb$updated_at": updatedAt,
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
		Dialect      Dialect
		ExpectedSQL  string
		ExpectedArgs []any
	}{
		{
			Name:        "empty dialect",
			ExpectedSQL: `INSERT INTO "table_name" ("pb$pk", "pb$sk", "pb$tenant_id", "pb$updated_at") VALUES ($1, $2, $3, $4) ON CONFLICT ON CONSTRAINT "pbidx_table_name" DO UPDATE SET "pb$pk"="excluded"."pb$pk","pb$sk"="excluded"."pb$sk","pb$updated_at"="excluded"."pb$updated_at" WHERE ("table_name"."pb$updated_at" <= "excluded"."pb$updated_at")`, //nolint:revive // query
			ExpectedArgs: []any{
				"pk",
				"sk",
				"tenant_id",
				"1900-01-01T00:00:00Z",
			},
		},
		{
			Name:        "unspecified dialect",
			Dialect:     DialectUnspecified,
			ExpectedSQL: `INSERT INTO "table_name" ("pb$pk", "pb$sk", "pb$tenant_id", "pb$updated_at") VALUES ($1, $2, $3, $4) ON CONFLICT ON CONSTRAINT "pbidx_table_name" DO UPDATE SET "pb$pk"="excluded"."pb$pk","pb$sk"="excluded"."pb$sk","pb$updated_at"="excluded"."pb$updated_at" WHERE ("table_name"."pb$updated_at" <= "excluded"."pb$updated_at")`, //nolint:revive // query
			ExpectedArgs: []any{
				"pk",
				"sk",
				"tenant_id",
				"1900-01-01T00:00:00Z",
			},
		},
		{
			Name:        "v13",
			Dialect:     DialectV13,
			ExpectedSQL: `INSERT INTO "table_name" ("pb$pk", "pb$sk", "pb$tenant_id", "pb$updated_at") VALUES ($1, $2, $3, $4) ON CONFLICT ON CONSTRAINT "pbidx_table_name" DO UPDATE SET "pb$pk"="excluded"."pb$pk","pb$sk"="excluded"."pb$sk","pb$updated_at"="excluded"."pb$updated_at" WHERE ("table_name"."pb$updated_at" <= "excluded"."pb$updated_at")`, //nolint:revive // query
			ExpectedArgs: []any{
				"pk",
				"sk",
				"tenant_id",
				"1900-01-01T00:00:00Z",
			},
		},
		{
			Name:        "v17",
			Dialect:     DialectV17,
			ExpectedSQL: `INSERT INTO "table_name" ("pb$pk", "pb$pksk", "pb$sk", "pb$tenant_id", "pb$updated_at") VALUES ($1, $2, $3, $4, $5) ON CONFLICT ON CONSTRAINT "pbidx_table_name" DO UPDATE SET "pb$pk"="excluded"."pb$pk","pb$sk"="excluded"."pb$sk","pb$updated_at"="excluded"."pb$updated_at" WHERE ("table_name"."pb$updated_at" <= "excluded"."pb$updated_at")`, //nolint:revive // query
			ExpectedArgs: []any{
				"pk",
				"pk|sk",
				"sk",
				"tenant_id",
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
				actualSQL, actualArgs, err = Insert(&mockDBReflect{}, tc.Dialect)
			})
			require.NoError(t, err)
			assert.Equal(t, tc.ExpectedSQL, actualSQL)
			assert.Equal(t, tc.ExpectedArgs, actualArgs)
		})
	}
}

func TestInsertMany(t *testing.T) {
	t.Parallel()
	tt := []struct {
		Name         string
		Dialect      Dialect
		Messages     []*mockDBReflect
		ExpectedSQL  string
		ExpectedArgs []any
	}{
		{
			Name:    "empty dialect - three rows",
			Dialect: DialectUnspecified,
			Messages: []*mockDBReflect{
				{tenantID: "t1", pk: "pk1", sk: "sk1", updatedAt: "2024-01-01T00:00:00Z"},
				{tenantID: "t1", pk: "pk2", sk: "sk2", updatedAt: "2024-01-02T00:00:00Z"},
				{tenantID: "t1", pk: "pk3", sk: "sk3", updatedAt: "2024-01-03T00:00:00Z"},
			},
			ExpectedSQL: `INSERT INTO "table_name" ("pb$pk", "pb$sk", "pb$tenant_id", "pb$updated_at") VALUES ($1, $2, $3, $4), ($5, $6, $7, $8), ($9, $10, $11, $12) ON CONFLICT ON CONSTRAINT "pbidx_table_name" DO UPDATE SET "pb$pk"="excluded"."pb$pk","pb$sk"="excluded"."pb$sk","pb$updated_at"="excluded"."pb$updated_at" WHERE ("table_name"."pb$updated_at" <= "excluded"."pb$updated_at")`, //nolint:revive // query
			ExpectedArgs: []any{
				"pk1", "sk1", "t1", "2024-01-01T00:00:00Z",
				"pk2", "sk2", "t1", "2024-01-02T00:00:00Z",
				"pk3", "sk3", "t1", "2024-01-03T00:00:00Z",
			},
		},
		{
			Name:    "v13 - three rows",
			Dialect: DialectV13,
			Messages: []*mockDBReflect{
				{tenantID: "t1", pk: "pk1", sk: "sk1", updatedAt: "2024-01-01T00:00:00Z"},
				{tenantID: "t1", pk: "pk2", sk: "sk2", updatedAt: "2024-01-02T00:00:00Z"},
				{tenantID: "t1", pk: "pk3", sk: "sk3", updatedAt: "2024-01-03T00:00:00Z"},
			},
			ExpectedSQL: `INSERT INTO "table_name" ("pb$pk", "pb$sk", "pb$tenant_id", "pb$updated_at") VALUES ($1, $2, $3, $4), ($5, $6, $7, $8), ($9, $10, $11, $12) ON CONFLICT ON CONSTRAINT "pbidx_table_name" DO UPDATE SET "pb$pk"="excluded"."pb$pk","pb$sk"="excluded"."pb$sk","pb$updated_at"="excluded"."pb$updated_at" WHERE ("table_name"."pb$updated_at" <= "excluded"."pb$updated_at")`, //nolint:revive // query
			ExpectedArgs: []any{
				"pk1", "sk1", "t1", "2024-01-01T00:00:00Z",
				"pk2", "sk2", "t1", "2024-01-02T00:00:00Z",
				"pk3", "sk3", "t1", "2024-01-03T00:00:00Z",
			},
		},
		{
			Name:    "v17 - three rows",
			Dialect: DialectV17,
			Messages: []*mockDBReflect{
				{tenantID: "t1", pk: "pk1", sk: "sk1", updatedAt: "2024-01-01T00:00:00Z"},
				{tenantID: "t1", pk: "pk2", sk: "sk2", updatedAt: "2024-01-02T00:00:00Z"},
				{tenantID: "t1", pk: "pk3", sk: "sk3", updatedAt: "2024-01-03T00:00:00Z"},
			},
			ExpectedSQL: `INSERT INTO "table_name" ("pb$pk", "pb$pksk", "pb$sk", "pb$tenant_id", "pb$updated_at") VALUES ($1, $2, $3, $4, $5), ($6, $7, $8, $9, $10), ($11, $12, $13, $14, $15) ON CONFLICT ON CONSTRAINT "pbidx_table_name" DO UPDATE SET "pb$pk"="excluded"."pb$pk","pb$sk"="excluded"."pb$sk","pb$updated_at"="excluded"."pb$updated_at" WHERE ("table_name"."pb$updated_at" <= "excluded"."pb$updated_at")`, //nolint:revive // query
			ExpectedArgs: []any{
				"pk1", "pk1|sk1", "sk1", "t1", "2024-01-01T00:00:00Z",
				"pk2", "pk2|sk2", "sk2", "t1", "2024-01-02T00:00:00Z",
				"pk3", "pk3|sk3", "sk3", "t1", "2024-01-03T00:00:00Z",
			},
		},
		{
			Name:    "v13 - two rows",
			Dialect: DialectV13,
			Messages: []*mockDBReflect{
				{tenantID: "t1", pk: "pk1", sk: "sk1", updatedAt: "2024-01-01T00:00:00Z"},
				{tenantID: "t1", pk: "pk2", sk: "sk2", updatedAt: "2024-01-02T00:00:00Z"},
			},
			ExpectedSQL: `INSERT INTO "table_name" ("pb$pk", "pb$sk", "pb$tenant_id", "pb$updated_at") VALUES ($1, $2, $3, $4), ($5, $6, $7, $8) ON CONFLICT ON CONSTRAINT "pbidx_table_name" DO UPDATE SET "pb$pk"="excluded"."pb$pk","pb$sk"="excluded"."pb$sk","pb$updated_at"="excluded"."pb$updated_at" WHERE ("table_name"."pb$updated_at" <= "excluded"."pb$updated_at")`, //nolint:revive // query
			ExpectedArgs: []any{
				"pk1", "sk1", "t1", "2024-01-01T00:00:00Z",
				"pk2", "sk2", "t1", "2024-01-02T00:00:00Z",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			var actualSQL string
			var actualArgs []any
			var err error

			// Convert to DBReflectMessage slice
			msgs := make([]DBReflectMessage, len(tc.Messages))
			for i, m := range tc.Messages {
				msgs[i] = m
			}

			require.NotPanics(t, func() {
				actualSQL, actualArgs, err = InsertMany(tc.Dialect, msgs...)
			})
			require.NoError(t, err)
			assert.Equal(t, tc.ExpectedSQL, actualSQL)
			assert.Equal(t, tc.ExpectedArgs, actualArgs)
		})
	}
}

func TestInsertMany_EmptySlice(t *testing.T) {
	t.Parallel()
	_, _, err := InsertMany[DBReflectMessage](DialectV13)
	require.Error(t, err)
	require.Contains(t, err.Error(), "at least one message")
}
