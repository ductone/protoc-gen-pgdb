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

	msgs := []pgdb_v1.DBReflectMessage{
		&Pasta{},
		&SauceIngredient{},
		&PastaIngredient{},
	}

	for _, smsg := range msgs {
		schema, err := pgdb_v1.CreateSchema(smsg)
		require.NoError(t, err)
		for _, line := range schema {
			//fmt.Printf("%s \n", line)
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

		if smsg.DBReflect().Descriptor().IsPartitioned() {
			require.Contains(t, ct, "PARTITION BY LIST")
		} else {
			require.NotContains(t, ct, "PARTITION BY LIST")
		}

		require.Equal(t, 1,
			strings.Count(ct, "fts_data"),
			"Create table should contain only one fts_data field: %s", ct,
		)

		_, err = pg.DB.Exec(ctx, `ALTER TABLE `+smsg.DBReflect().Descriptor().TableName()+` DROP COLUMN "pb$id"`)
		require.NoErrorf(t, err, "TestSchemaFoodPasta: failed to drop col id")

		schema, err = pgdb_v1.Migrations(ctx, pg.DB, smsg)
		require.NoError(t, err)
		for _, line := range schema {
			//fmt.Println(line)
			_, err := pg.DB.Exec(ctx, line)
			require.NoErrorf(t, err, "TestSchemaFoodPasta: failed to execute sql: '\n%s\n'", line)
			// os.Stderr.WriteString(line)
			// os.Stderr.WriteString("\n------\n")
		}
		if len(schema) > 0 {
			ct = schema[0]
			require.Contains(t, ct, "ALTER TABLE")
		}

		fakeTenantIds := []string{"t1", "t2", "t3"}
		protoTableName := smsg.DBReflect().Descriptor().TableName()

		if smsg.DBReflect().Descriptor().IsPartitioned() {
			// TODO(scott) makes helper functions and move to pgdb_v1
			verifyMasterPartition(t, pg, protoTableName, fakeTenantIds)
			// Test sub-tables for partitions
			// Create sub-tables
			testCreatePartitionTables(t, pg, smsg, fakeTenantIds)
			verifySubTables(t, pg, protoTableName, fakeTenantIds)
			// Insert data into master table
			testInsertAndVerify(t, pg, protoTableName, fakeTenantIds, smsg)

		}
	}

}

func testCreatePartitionTables(t *testing.T, pg *pgtest.PG, msg pgdb_v1.DBReflectMessage, fakeTenantIds []string) {
	ctx := context.Background()
	// Create sub-tables
	tenantIter := TestTenantIterator(ctx, fakeTenantIds)
	// Don't really need tenantId in update func but good for logging purposes.
	pgdb_v1.TenantPartitionsUpdate(ctx, msg, tenantIter, func(ctx context.Context, tenantId string, schema string) error {
		_, err := pg.DB.Exec(ctx, schema)
		require.NoError(t, err)
		return nil
	})
}

func verifyMasterPartition(t *testing.T, pg *pgtest.PG, tableName string, fakeTenantIds []string) {
	ctx := context.Background()
	fmt.Println(tableName)
	// Verify number of master partition tables
	partTablesQuery := `SELECT count(t.tablename), t.tablename
		FROM pg_tables t
		LEFT JOIN pg_partitioned_table p ON p.partrelid = (SELECT oid FROM pg_class WHERE relname = t.tablename)
		WHERE t.schemaname NOT IN ('pg_catalog', 'information_schema') AND p.partrelid IS NOT NULL AND t.tablename = $1
		GROUP BY t.tablename;`

	rows, err := pg.DB.Query(ctx, partTablesQuery, tableName)
	require.NoErrorf(t, err, "TestSchemaFoodPasta: failed to count partitioned tables query: '\n%s\n'", partTablesQuery)
	defer rows.Close()
	var partTableCount int
	var queryTableName string

	for rows.Next() {
		err = rows.Scan(&partTableCount, &queryTableName)
		require.NoError(t, err)
	}

	require.NoError(t, rows.Err())
	require.Equal(t, 1, partTableCount, "Should have one master partition table")
	require.Equal(t, tableName, queryTableName, "Should have one master partition table")

}

func verifySubTables(t *testing.T, pg *pgtest.PG, tableName string, fakeTenantIds []string) {
	ctx := context.Background()
	// Verify number of sub tables
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

	rows, err := pg.DB.Query(ctx, fmt.Sprintf(sqlSubTables, tableName))
	require.NoErrorf(t, err, "TestSchemaFoodPasta: failed to count partitioned tables query: '\n%s\n'", sqlSubTables)
	defer rows.Close()

	var parentTable string
	var childTable string
	rowCount := 0
	selectedSubTableNames := make([]string, 0, len(fakeTenantIds))

	for rows.Next() {
		err = rows.Scan(&parentTable, &childTable)
		fmt.Printf("parent: %s, child: %s\n", parentTable, childTable)
		require.NoError(t, err)
		require.Equal(t, tableName, parentTable, "Parent table name did not match proto")
		selectedSubTableNames = append(selectedSubTableNames, childTable)
		rowCount += 1
	}

	require.NoError(t, rows.Err())
	require.Equal(t, len(fakeTenantIds), rowCount, "Should have one sub-partition table per fake tenant")
}

func testInsertAndVerify(t *testing.T, pg *pgtest.PG, tableName string, fakeTenantIds []string, msg pgdb_v1.DBReflectMessage) {
	ctx := context.Background()
	// Insert data into master table
	// Verify data in master table
	// Verify data in sub tables
	var obj1 pgdb_v1.DBReflectMessage
	var obj2 pgdb_v1.DBReflectMessage
	var obj3 pgdb_v1.DBReflectMessage
	switch msg.(type) {
	case *Pasta:
		obj1 = &Pasta{
			TenantId: "t1",
			Id:       "p1",
		}
		obj2 = &Pasta{
			TenantId: "t2",
			Id:       "p2",
		}
		obj3 = &Pasta{
			TenantId: "t3",
			Id:       "p3",
		}
	case *SauceIngredient:
		obj1 = &SauceIngredient{
			TenantId: "t1",
			Id:       "p1",
		}
		obj2 = &SauceIngredient{
			TenantId: "t2",
			Id:       "p2",
		}
		obj3 = &SauceIngredient{
			TenantId: "t3",
			Id:       "p3",
		}
	case *PastaIngredient:
		obj1 = &PastaIngredient{
			TenantId:     "t1",
			IngredientId: "p1",
		}
		obj2 = &PastaIngredient{
			TenantId:     "t2",
			IngredientId: "p2",
		}
		obj3 = &PastaIngredient{
			TenantId:     "t3",
			IngredientId: "p3",
		}
	}
	sql, args, err := pgdb_v1.Insert(obj1)
	require.NoError(t, err)
	_, err = pg.DB.Exec(ctx, sql, args...)
	require.NoError(t, err)

	sql, args, err = pgdb_v1.Insert(obj2)
	require.NoError(t, err)
	_, err = pg.DB.Exec(ctx, sql, args...)
	require.NoError(t, err)

	sql, args, err = pgdb_v1.Insert(obj3)
	require.NoError(t, err)
	_, err = pg.DB.Exec(ctx, sql, args...)
	require.NoError(t, err)

	var tenantIdSelect string
	selectColStr := msg.DBReflect().Descriptor().TenantField().Name
	// Test select from master table
	masterSelectSql := `SELECT %s FROM %s`
	fmtSql := fmt.Sprintf(masterSelectSql, selectColStr, tableName)
	rows, err := pg.DB.Query(ctx, fmtSql)
	require.NoError(t, err)
	defer rows.Close()
	rowCount := 0
	for rows.Next() {
		err = rows.Scan(&tenantIdSelect)
		require.NoError(t, err)
		rowCount += 1
	}
	fmt.Println()
	require.NoError(t, rows.Err())
	require.Equal(t, len(fakeTenantIds), rowCount, "Should have one row per tenant")

	subTables, err := pgdb_v1.ReadPartitionSubTables(ctx, pg.DB, msg.DBReflect().Descriptor())
	require.NoError(t, err)

	// Test select each tenant
	for _, subTable := range subTables {
		selectSql := `SELECT %s FROM %s`
		fmtSql := fmt.Sprintf(selectSql, selectColStr, subTable)
		rows, err := pg.DB.Query(ctx, fmtSql)
		require.NoError(t, err)
		defer rows.Close()
		rowCount := 0
		for rows.Next() {
			err = rows.Scan(&tenantIdSelect)
			require.NoError(t, err)
			rowCount += 1
		}
		require.NoError(t, rows.Err())
		require.Equal(t, 1, rowCount, "Should have one row per tenant table")
	}
}

func TestTenantIterator(ctx context.Context, tenantList []string) pgdb_v1.TenantIteratorFunc {
	index := 0

	return func(ctx context.Context) (string, error) {
		if index >= len(tenantList) {
			return "", nil
		}
		tenantId := tenantList[index]
		index += 1
		return tenantId, nil
	}

}
