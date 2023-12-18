package v1

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/ductone/protoc-gen-pgdb/internal/pgtest"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	"github.com/stretchr/testify/require"
)

func TestSchemaFoodPasta(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	smsg := &Pasta{}
	schema, err := pgdb_v1.CreateSchema(smsg)
	require.NoError(t, err)
	for _, line := range schema {
		fmt.Printf("%s \n", line)
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "TestSchemaFoodPasta: failed to execute sql: '\n%s\n'", line)
		// os.Stderr.WriteString(line)
		// os.Stderr.WriteString("\n------\n")
	}
	ct := schema[0]
	require.Contains(t, ct, "CREATE TABLE")
	require.Equal(t, 2,
		strings.Count(ct, "$pksk"),
		"Create table should contain only one pksk field + index: %s", ct,
	)

	require.Contains(t, ct, "PARTITION BY LIST")

	require.Equal(t, 1,
		strings.Count(ct, "fts_data"),
		"Create table should contain only one fts_data field: %s", ct,
	)
	_, err = pg.DB.Exec(ctx, "DROP TABLE "+smsg.DBReflect().Descriptor().TableName())
	require.NoErrorf(t, err, "TestSchemaFoodPasta: failed to drop")

	schema, err = pgdb_v1.Migrations(ctx, pg.DB, smsg)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "TestSchemaFoodPasta: failed to execute sql: '\n%s\n'", line)
		// os.Stderr.WriteString(line)
		// os.Stderr.WriteString("\n------\n")
	}
	ct = schema[0]
	require.Contains(t, ct, "CREATE TABLE")

	// Verify number of master partition tables
	partTablesQuery := `SELECT count(t.tablename), t.tablename
		FROM pg_tables t
		LEFT JOIN pg_partitioned_table p ON p.partrelid = (SELECT oid FROM pg_class WHERE relname = t.tablename)
		WHERE t.schemaname NOT IN ('pg_catalog', 'information_schema') AND p.partrelid IS NOT NULL
		GROUP BY t.tablename;`

	rows, err := pg.DB.Query(ctx, partTablesQuery)
	require.NoErrorf(t, err, "TestSchemaFoodPasta: failed to count partitioned tables query: '\n%s\n'", partTablesQuery)
	defer rows.Close()
	var partTableCount int
	var tableName string

	for rows.Next() {
		err = rows.Scan(&partTableCount, &tableName)
		require.NoError(t, err)
		fmt.Printf("%s %d\n", tableName, partTableCount)
	}

	require.NoError(t, rows.Err())
	require.Equal(t, 1, partTableCount, "Should have one master partition table")

}
