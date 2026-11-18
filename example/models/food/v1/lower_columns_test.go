package v1

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ductone/protoc-gen-pgdb/internal/pgtest"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
)

func findIndex(t *testing.T, indexes []*pgdb_v1.Index, key string) *pgdb_v1.Index {
	t.Helper()
	for _, idx := range indexes {
		if strings.Contains(idx.Name, key) {
			return idx
		}
	}
	require.Failf(t, "index not found", "no index name contains %q", key)
	return nil
}

// TestLowerColumnsIndex covers the lower_columns index option: the generated
// expression wraps only the listed columns in lower(), keeps Columns populated,
// composes with partial_deleted_at_is_null, honours nested column prefixes, and
// renders DDL that Postgres accepts.
func TestLowerColumnsIndex(t *testing.T) {
	for _, dialect := range []pgdb_v1.Dialect{pgdb_v1.DialectV13, pgdb_v1.DialectV17} {
		desc := (&PastaIngredient{}).DBReflect(dialect).Descriptor()
		indexes := desc.Indexes()

		ci := findIndex(t, indexes, "pasta_id_ci_pasta_ingredient")
		require.Equal(t, []string{"pb$tenant_id", "pb$pasta_id"}, ci.Columns)
		require.Equal(t, "pb$tenant_id, lower(pb$pasta_id)", ci.OverrideExpression)
		require.Empty(t, ci.WherePredicate)

		alive := findIndex(t, indexes, "pasta_id_ci_alive_only")
		require.Equal(t, []string{"pb$tenant_id", "pb$pasta_id"}, alive.Columns)
		require.Equal(t, "pb$tenant_id, lower(pb$pasta_id)", alive.OverrideExpression)
		require.Equal(t, "pb$deleted_at IS NULL", alive.WherePredicate)

		// Nested: the lowered column picks up the nested prefix, tenant_id is
		// retargeted to the top-level column.
		nested := desc.Indexes(pgdb_v1.NewIndexOptions(nil).Nested("11$")...)
		nci := findIndex(t, nested, "pasta_id_ci_pasta_ingredient")
		require.Equal(t, "pb$tenant_id, lower(pb$11$pasta_id)", nci.OverrideExpression)
	}

	schema, err := pgdb_v1.IndexSchema(&PastaIngredient{}, pgdb_v1.DialectV13)
	require.NoError(t, err)
	var ciSQL, aliveSQL string
	for _, line := range schema {
		switch {
		case strings.Contains(line, "pasta_id_ci_alive_only"):
			aliveSQL = line
		case strings.Contains(line, "pasta_id_ci"):
			ciSQL = line
		}
	}
	require.Contains(t, ciSQL, "USING\n  BTREE\n(\npb$tenant_id, lower(pb$pasta_id)\n)\n")
	require.NotContains(t, ciSQL, "WHERE")
	require.Contains(t, aliveSQL, "USING\n  BTREE\n(\npb$tenant_id, lower(pb$pasta_id)\n)\nWHERE pb$deleted_at IS NULL\n")
}

func TestLowerColumnsIndexPostgres(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin; CREATE EXTENSION IF NOT EXISTS vector;")
	require.NoError(t, err)

	msg := &PastaIngredient{}
	schema, err := pgdb_v1.CreateSchema(msg, pgdb_v1.DialectV13)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "failed to execute sql: '\n%s\n'", line)
	}

	tableName := msg.DBReflect(pgdb_v1.DialectV13).Descriptor().TableName()
	rows, err := pg.DB.Query(ctx,
		`SELECT indexdef FROM pg_indexes WHERE tablename = $1 AND indexdef LIKE '%lower(%'`, tableName)
	require.NoError(t, err)
	defs := []string{}
	for rows.Next() {
		var def string
		require.NoError(t, rows.Scan(&def))
		defs = append(defs, def)
	}
	require.NoError(t, rows.Err())
	require.Len(t, defs, 3, "expected all lower_columns indexes: %v", defs)
	pastaDefs := 0
	for _, def := range defs {
		if strings.Contains(def, `lower("pb$pasta_id")`) {
			pastaDefs++
		} else {
			require.Contains(t, def, `lower("pb$label")`)
		}
	}
	require.Equal(t, 2, pastaDefs)
}

// TestLowerColumnsSafeOperators checks that a column covered only by a
// lower(column) index gets no plain-column safe operators, while a column that
// another plain index also covers keeps them.
func TestLowerColumnsSafeOperators(t *testing.T) {
	builder := reflect.TypeOf(&PastaIngredientDBQueryBuilder{})
	unsafe := reflect.TypeOf(&PastaIngredientDBQueryUnsafe{})

	_, hasSafeLabel := builder.MethodByName("Label")
	require.False(t, hasSafeLabel, "label is only covered by a lower() index")
	_, hasUnsafeLabel := unsafe.MethodByName("Label")
	require.True(t, hasUnsafeLabel, "label must stay reachable via Unsafe()")

	pastaID, ok := builder.MethodByName("PastaId")
	require.True(t, ok, "pasta_id is also covered by plain indexes")
	ops := pastaID.Type.Out(0)
	for _, op := range []string{"Eq", "Gt", "Gte", "Lt", "Lte", "IsNotEmpty"} {
		_, ok := ops.MethodByName(op)
		require.Truef(t, ok, "PastaId safe operators should keep %s", op)
	}

	desc := (&PastaIngredient{}).DBReflect(pgdb_v1.DialectV13).Descriptor()
	label := findIndex(t, desc.Indexes(), "label_ci")
	require.Equal(t, []string{"pb$tenant_id", "pb$label"}, label.Columns)
	require.Equal(t, "pb$tenant_id, lower(pb$label)", label.OverrideExpression)
}
