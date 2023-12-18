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

	// TODO(scott) makes helper functions and move to pgdb_v1

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
	require.Equal(t, smsg.DBReflect().Descriptor().TableName(), tableName, "Table name did not match proto")

	// Test sub-tables for partitions
	fakeTenantIds := []string{"t1", "t2", "t3"}
	protoTableName := smsg.DBReflect().Descriptor().TableName()
	// Create sub-tables
	for _, tenantId := range fakeTenantIds {
		_, err = pg.DB.Exec(ctx, fmt.Sprintf("CREATE TABLE %s_%s PARTITION OF %s FOR VALUES IN ('%s')", protoTableName, tenantId, protoTableName, tenantId))
		require.NoErrorf(t, err, "TestSchemaFoodPasta: failed to create partitioned table: '\n%s\n'", ct)
	}

	// Verify sub-partition tables
	sqlSubTables := `SELECT
			parent.relname      AS parent,
			child.relname       AS child
		FROM pg_inherits
			JOIN pg_class parent            ON pg_inherits.inhparent = parent.oid
			JOIN pg_class child             ON pg_inherits.inhrelid   = child.oid
			JOIN pg_namespace nmsp_parent   ON nmsp_parent.oid  = parent.relnamespace
			JOIN pg_namespace nmsp_child    ON nmsp_child.oid   = child.relnamespace
		WHERE parent.relname='%s';`

	rows, err = pg.DB.Query(ctx, fmt.Sprintf(sqlSubTables, protoTableName))
	require.NoErrorf(t, err, "TestSchemaFoodPasta: failed to count partitioned tables query: '\n%s\n'", sqlSubTables)
	defer rows.Close()

	var parentTable string
	var childTable string
	rowCount := 0
	selectedSubTableNames := make([]string, 0, len(fakeTenantIds))

	for rows.Next() {
		err = rows.Scan(&parentTable, &childTable)
		require.NoError(t, err)
		fmt.Printf("%s: %s\n", parentTable, childTable)
		require.Equal(t, protoTableName, parentTable, "Parent table name did not match proto")
		selectedSubTableNames = append(selectedSubTableNames, childTable)
		rowCount += 1
	}

	require.NoError(t, rows.Err())
	require.Equal(t, len(fakeTenantIds), rowCount, "Should have one sub-partition table per fake tenant")
	for _, tenantId := range fakeTenantIds {
		require.Contains(t, selectedSubTableNames, fmt.Sprintf("%s_%s", protoTableName, tenantId), "Should have a sub-partition table for tenant %s", tenantId)
	}

}
