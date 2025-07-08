package v1

import (
	"context"
	"fmt"
	"net/netip"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"github.com/ductone/protoc-gen-pgdb/internal/pgtest"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"

	"github.com/segmentio/ksuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	llm_v1 "github.com/ductone/protoc-gen-pgdb/example/models/llm/v1"
)

type testTable struct {
	objects []pgdb_v1.DBReflectMessage
}

func TestSchemaFoodPasta(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin; CREATE EXTENSION IF NOT EXISTS vector;")
	require.NoError(t, err)

	testobjects := []testTable{
		{
			objects: []pgdb_v1.DBReflectMessage{
				Pasta_builder{
					TenantId: "t1",
					Id:       "p1",
				}.Build(),
				Pasta_builder{
					TenantId: "t2",
					Id:       "p2",
				}.Build(),
				Pasta_builder{
					TenantId: "t3",
					Id:       "p3",
				}.Build(),
			},
		},
		{
			objects: []pgdb_v1.DBReflectMessage{
				SauceIngredient_builder{
					TenantId:   "t1",
					Id:         "s1",
					SourceAddr: "127.0.0.1",
				}.Build(),
				SauceIngredient_builder{
					TenantId:   "t2",
					Id:         "s2",
					SourceAddr: "1.2.3.4",
				}.Build(),
				SauceIngredient_builder{
					TenantId:   "t3",
					Id:         "s3",
					SourceAddr: "2001:db8:abcd:12::1",
				}.Build(),
			},
		},
		{
			objects: []pgdb_v1.DBReflectMessage{
				PastaIngredient_builder{
					TenantId: "t1",
					Id:       "pi1",
					ModelEmbeddings: []*PastaIngredient_ModelEmbedding{
						PastaIngredient_ModelEmbedding_builder{
							Embedding: []float32{},
							Model:     llm_v1.Model_MODEL_3DIMS,
						}.Build(),
					},
				}.Build(),
				PastaIngredient_builder{
					TenantId: "t2",
					Id:       "pi2",
					ModelEmbeddings: []*PastaIngredient_ModelEmbedding{
						PastaIngredient_ModelEmbedding_builder{
							Embedding: []float32{4.0, 5.0, 6.0},
							Model:     llm_v1.Model_MODEL_3DIMS,
						}.Build(),
					},
				}.Build(),
				PastaIngredient_builder{
					TenantId: "t3",
					Id:       "pi3",
					ModelEmbeddings: []*PastaIngredient_ModelEmbedding{
						PastaIngredient_ModelEmbedding_builder{
							Embedding: []float32{1.0, 2.0, 3.0},
							Model:     llm_v1.Model_MODEL_3DIMS,
						}.Build(),
						PastaIngredient_ModelEmbedding_builder{
							Embedding: []float32{4.0, 5.0, 6.0},
							Model:     llm_v1.Model_MODEL_3DIMS,
						}.Build(),
					},
				}.Build(),
			},
		},
	}

	for _, testobj := range testobjects {
		smsg := testobj.objects[0]
		schema, err := pgdb_v1.CreateSchema(smsg, pgdb_v1.DialectV13)
		require.NoError(t, err)
		for _, line := range schema {
			// fmt.Printf("%s \n", line)
			_, err := pg.DB.Exec(ctx, line)
			require.NoErrorf(t, err, "TestSchemaFoodPasta: failed to execute sql: '\n%s\n'", line)
		}
		ct := schema[0]
		require.Contains(t, ct, "CREATE TABLE")
		require.Equal(t, 3,
			strings.Count(ct, "$pksk"),
			"Create table should contain only one pksk field + index: %s", ct,
		)

		if smsg.DBReflect(pgdb_v1.DialectV13).Descriptor().IsPartitioned() {
			require.Contains(t, ct, "PARTITION BY LIST")
		} else {
			require.NotContains(t, ct, "PARTITION BY LIST")
		}

		hnswIndexCount := 0
		partialIndexCount := 0
		for _, line := range schema {
			if strings.Contains(line, "HNSW") {
				// fmt.Printf("%s \n", line)
				hnswIndexCount += 1
			}
			if strings.Contains(line, "deleted_at IS NULL") {
				partialIndexCount += 1
			}
		}
		// fmt.Printf("hnswIndexCount: %d\n", hnswIndexCount)
		if _, ok := smsg.(*PastaIngredient); ok {
			require.Equal(t, 2, hnswIndexCount, "Should have 2 hnsw indexes") // 2 enums = 2 indexes
			require.Equal(t, 1, partialIndexCount, "Should have 1 partial index")
		} else {
			require.Equal(t, 0, hnswIndexCount, "Should have 0 hnsw indexes")
			require.Equal(t, 0, partialIndexCount, "Should have 0 partial index")
		}

		require.Equal(t, 1,
			strings.Count(ct, "fts_data"),
			"Create table should contain only one fts_data field: %s", ct,
		)

		_, err = pg.DB.Exec(ctx, `ALTER TABLE `+smsg.DBReflect(pgdb_v1.DialectV13).Descriptor().TableName()+` DROP COLUMN "pb$id"`)
		require.NoError(t, err, "TestSchemaFoodPasta: failed to drop col id")

		schema, err = pgdb_v1.Migrations(ctx, pg.DB, smsg, pgdb_v1.DialectV13)
		require.NoError(t, err)
		for _, line := range schema {
			// fmt.Println(line)
			_, err := pg.DB.Exec(ctx, line)
			require.NoErrorf(t, err, "TestSchemaFoodPasta: failed to execute sql: '\n%s\n'", line)
		}
		if len(schema) > 0 {
			ct = schema[0]
			require.Contains(t, ct, "ALTER TABLE")
		}

		for _, line := range schema {
			_, err = pg.DB.Exec(ctx, line)
			require.NoErrorf(t, err, "TestSchemaFoodPasta: failed to repair table: '%s'", line)
		}
		schema, err = pgdb_v1.Migrations(ctx, pg.DB, smsg, pgdb_v1.DialectV13)
		require.NoError(t, err)
		require.Len(t, schema, 0, "Should have no migrations after repair")

		fakeTenantIds := []string{"t1", "t2", "t3"}
		protoTableName := smsg.DBReflect(pgdb_v1.DialectV13).Descriptor().TableName()

		if smsg.DBReflect(pgdb_v1.DialectV13).Descriptor().IsPartitioned() {
			verifyMasterPartition(t, pg, protoTableName, fakeTenantIds)
			// Test sub-tables for partitions
			// Create sub-tables
			testCreatePartitionTables(t, pg, smsg, fakeTenantIds)
			verifySubTables(t, pg, protoTableName, fakeTenantIds)
			// Insert data into master table
			testInsertAndVerify(t, pg, protoTableName, fakeTenantIds, testobj.objects)
		}
	}

}

func testCreatePartitionTables(t *testing.T, pg *pgtest.PG, msg pgdb_v1.DBReflectMessage, fakeTenantIds []string) {
	ctx := context.Background()
	// Create sub-tables
	tenantIter := TenantIteratorTest(ctx, fakeTenantIds)
	// Don't really need tenantId in update func but good for logging purposes.
	err := pgdb_v1.TenantPartitionsUpdate(ctx, pg.DB, msg, pgdb_v1.DialectV13, tenantIter, func(ctx context.Context, schema string, args ...interface{}) error {
		_, err := pg.DB.Exec(ctx, schema, args...)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)
}

func verifyMasterPartition(t *testing.T, pg *pgtest.PG, tableName string, fakeTenantIds []string) {
	ctx := context.Background()
	// fmt.Println(tableName)
	// Verify number of master partition tables
	partTablesQuery := `SELECT count(t.tablename), t.tablename
		FROM pg_tables t
		LEFT JOIN pg_partitioned_table p ON p.partrelid = (SELECT oid FROM pg_class WHERE relname = t.tablename)
		WHERE t.schemaname = 'public' AND p.partrelid IS NOT NULL AND t.tablename = $1
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
	require.Equal(t, tableName, queryTableName, "Should have matching table names")

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
		// fmt.Printf("parent: %s, child: %s\n", parentTable, childTable)
		require.NoError(t, err)
		require.Equal(t, tableName, parentTable, "Parent table name did not match proto")
		selectedSubTableNames = append(selectedSubTableNames, childTable)
		rowCount += 1
	}

	require.NoError(t, rows.Err())
	require.Equal(t, len(fakeTenantIds), rowCount, "Should have one sub-partition table per fake tenant: %v", selectedSubTableNames)
}

func testInsertAndVerify(t *testing.T, pg *pgtest.PG, tableName string, fakeTenantIds []string, objects []pgdb_v1.DBReflectMessage) {
	ctx := context.Background()
	// Insert data into master table
	// Verify data in master table
	// Verify data in sub tables
	msg := objects[0]
	sql, args, err := pgdb_v1.Insert(objects[0], pgdb_v1.DialectV13)
	require.NoError(t, err)
	// fmt.Printf("sql: %s\n\n%v\n", sql, args)
	_, err = pg.DB.Exec(ctx, sql, args...)
	require.NoError(t, err)

	sql, args, err = pgdb_v1.Insert(objects[1], pgdb_v1.DialectV13)
	require.NoError(t, err)
	_, err = pg.DB.Exec(ctx, sql, args...)
	require.NoError(t, err, "Failed to insert object: %v", objects[1])

	sql, args, err = pgdb_v1.Insert(objects[2], pgdb_v1.DialectV13)
	require.NoError(t, err)
	_, err = pg.DB.Exec(ctx, sql, args...)
	require.NoError(t, err)

	var tenantIdSelect string
	selectColStr := msg.DBReflect(pgdb_v1.DialectV13).Descriptor().TenantField().Name
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

	subTables, err := readPartitionSubTables(ctx, pg.DB, msg.DBReflect(pgdb_v1.DialectV13).Descriptor())
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

func TenantIteratorTest(ctx context.Context, tenantList []string) pgdb_v1.TenantIteratorFunc {
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

type sqlScanner interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// this
func readPartitionSubTables(ctx context.Context, db sqlScanner, desc pgdb_v1.Descriptor) ([]string, error) {
	dialect := goqu.Dialect("postgres")

	qb := dialect.From("pg_inherits")
	qb = qb.Select("child.relname").As("child")
	qb = qb.Join(goqu.T("pg_class").As("parent"), goqu.On(goqu.I("pg_inherits.inhparent").Eq(goqu.I("parent.oid"))))
	qb = qb.Join(goqu.T("pg_class").As("child"), goqu.On(goqu.I("pg_inherits.inhrelid").Eq(goqu.I("child.oid"))))
	qb = qb.Join(goqu.T("pg_namespace").As("nmsp_parent"), goqu.On(goqu.I("nmsp_parent.oid").Eq(goqu.I("parent.relnamespace"))))
	qb = qb.Join(goqu.T("pg_namespace").As("nmsp_child"), goqu.On(goqu.I("nmsp_child.oid").Eq(goqu.I("child.relnamespace"))))
	qb = qb.Where(goqu.L("parent.relname = ?", desc.TableName()))
	query, params, err := qb.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tables := make([]string, 0)
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

func fixtureSchemaSauceIngredient(t *testing.T, pg *pgtest.PG) {
	ctx := context.Background()
	_, err := pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	schema, err := pgdb_v1.CreateSchema(&SauceIngredient{}, pgdb_v1.DialectV13)
	require.NoError(t, err)
	for _, line := range schema {
		// fmt.Printf("%s \n", line)
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "TestSchemaSauceIngredientNetworkRange: failed to execute sql: '\n%s\n'", line)
	}

	data := []*SauceIngredient{
		SauceIngredient_builder{
			TenantId:   "t1",
			Id:         "s1",
			SourceAddr: "127.0.0.1",
		}.Build(),
		SauceIngredient_builder{
			TenantId:   "t2",
			Id:         "s2",
			SourceAddr: "1.2.3.4",
		}.Build(),
		SauceIngredient_builder{
			TenantId:   "t3",
			Id:         "s3",
			SourceAddr: "2001:db8:abcd:12::1",
		}.Build(),
	}

	for _, row := range data {
		sql, args, err := pgdb_v1.Insert(row, pgdb_v1.DialectV13)
		require.NoError(t, err)
		_, err = pg.DB.Exec(ctx, sql, args...)
		require.NoError(t, err)
	}
}

func TestSchemaSauceIngredientNetworkRange(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	fixtureSchemaSauceIngredient(t, pg)

	nilSI := (*SauceIngredient)(nil)
	tableName := nilSI.DB().TableName()
	fields := nilSI.DB().Query()

	inRange := fields.SourceAddr().InNetworkPrefix(netip.MustParsePrefix("1.2.3.0/24"))

	qb := goqu.Dialect("postgres")
	query, args, err := qb.Select(
		exp.NewAliasExpression(fields.TenantId().Identifier(), "tenant_id"),
		exp.NewAliasExpression(fields.SourceAddr().Identifier(), "source_addr"),
	).From(tableName).
		Where(inRange).
		ToSQL()
	require.NoError(t, err)
	rows, err := pg.DB.Query(ctx, query, args...)
	require.NoError(t, err)
	defer rows.Close()
	count := 0
	for rows.Next() {
		count++
		values, err := rows.Values()
		require.NoError(t, err)
		require.Equal(t, []any{
			"t2",
			netip.MustParsePrefix("1.2.3.4/32"),
		}, values)
	}
	require.Equal(t, 1, count)

}

func TestSchemaSauceIngredientInBehavoirs(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()
	fixtureSchemaSauceIngredient(t, pg)

	nilSI := (*SauceIngredient)(nil)
	tableName := nilSI.DB().TableName()
	fields := nilSI.DB().Query()

	tc := []struct {
		name  string
		count int
		query exp.Expression
	}{
		{
			name:  "emptyInQuery",
			count: 0,
			query: fields.SourceAddr().In([]string{}),
		},
	}

	qb := goqu.Dialect("postgres")

	for _, tx := range tc {
		t.Run(tx.name, func(t *testing.T) {
			query, args, err := qb.Select(
				exp.NewAliasExpression(fields.TenantId().Identifier(), "tenant_id"),
				exp.NewAliasExpression(fields.SourceAddr().Identifier(), "source_addr"),
			).From(tableName).
				Where(tx.query).
				ToSQL()
			require.NoError(t, err)

			rows, err := pg.DB.Query(ctx, query, args...)
			require.NoError(t, err)
			defer rows.Close()

			count := 0
			for rows.Next() {
				count++
				_, err := rows.Values()
				require.NoError(t, err)
			}
			require.Equal(t, tx.count, count)
		})

	}
}

func TestDatePartitionsUpdate(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	msg := GarlicIngredient_builder{
		TenantId: "t1",
		Id:       "pi1",
	}.Build()

	// Create the table first
	schema, err := pgdb_v1.CreateSchema(msg, pgdb_v1.DialectV13)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "failed to execute sql: '\n%s\n'", line)
	}

	// Set up test dates
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)

	// Create the partitions
	err = pgdb_v1.DatePartitionsUpdate(ctx, pg.DB, msg, pgdb_v1.DialectV13, startDate, endDate, func(ctx context.Context, schema string, args ...interface{}) error {
		_, err := pg.DB.Exec(ctx, schema, args...)
		return err
	})
	require.NoError(t, err)

	// Verify the partitions were created
	subTables, err := readPartitionSubTables(ctx, pg.DB, msg.DBReflect(pgdb_v1.DialectV13).Descriptor())
	require.NoError(t, err)

	// Should have 3 partitions (Jan, Feb, Mar)
	require.Equal(t, 3, len(subTables), "Should have created 3 monthly partitions")

	// Verify partition names follow expected pattern
	expectedNames := []string{
		msg.DBReflect(pgdb_v1.DialectV13).Descriptor().TableName() + "_2024_01",
		msg.DBReflect(pgdb_v1.DialectV13).Descriptor().TableName() + "_2024_02",
		msg.DBReflect(pgdb_v1.DialectV13).Descriptor().TableName() + "_2024_03",
	}

	for i, expected := range expectedNames {
		require.Equal(t, expected, subTables[i], "Partition table name mismatch")
	}

	// Test data insertion into partitions
	testData := []*GarlicIngredient{
		GarlicIngredient_builder{
			TenantId:   "t1",
			Id:         "pi1",
			SourceAddr: "127.0.0.1",
			CreatedAt:  timestamppb.New(time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)),
		}.Build(),
		GarlicIngredient_builder{
			TenantId:   "t2",
			Id:         "pi2",
			SourceAddr: "1.2.3.4",
			CreatedAt:  timestamppb.New(time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)),
		}.Build(),
		GarlicIngredient_builder{
			TenantId:   "t3",
			Id:         "pi3",
			SourceAddr: "2001:db8:abcd:12::1",
			CreatedAt:  timestamppb.New(time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)),
		}.Build(),
	}

	// Insert test data
	for _, data := range testData {
		sql, args, err := pgdb_v1.Insert(data, pgdb_v1.DialectV13)
		require.NoError(t, err)
		_, err = pg.DB.Exec(ctx, sql, args...)
		require.NoError(t, err)
	}

	// Verify data distribution across partitions
	for i, subTable := range subTables {
		var count int
		err := pg.DB.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", subTable)).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count, "Each partition should have exactly one row")

		var createdAt time.Time
		err = pg.DB.QueryRow(ctx, fmt.Sprintf("SELECT pb$created_at FROM %s", subTable)).Scan(&createdAt)
		require.NoError(t, err)
		require.Equal(t, testData[i].GetCreatedAt().AsTime().UTC().Format("2006-01"),
			createdAt.UTC().Format("2006-01"),
			"Data should be in correct monthly partition")
	}
}

func TestEventIDPartitionsUpdate(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	msg := CheeseIngredient_builder{
		TenantId: "t1",
		Id:       "pi1",
	}.Build()

	// Create the table first
	schema, err := pgdb_v1.CreateSchema(msg, pgdb_v1.DialectV13)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "failed to execute sql: '\n%s\n'", line)
	}

	// Set up test dates - KSUIDs will be generated within this range
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)

	// Create the partitions
	err = pgdb_v1.EventIDPartitionsUpdate(ctx, pg.DB, msg, pgdb_v1.DialectV13, startDate, endDate, func(ctx context.Context, schema string, args ...interface{}) error {
		_, err := pg.DB.Exec(ctx, schema)
		return err
	})
	require.NoError(t, err)

	// Verify the partitions were created
	subTables, err := readPartitionSubTables(ctx, pg.DB, msg.DBReflect(pgdb_v1.DialectV13).Descriptor())
	require.NoError(t, err)

	// Should have 3 partitions (Jan, Feb, Mar)
	require.Equal(t, 3, len(subTables), "Should have created 3 monthly partitions")

	// Verify partition names follow expected pattern
	expectedNames := []string{
		msg.DBReflect(pgdb_v1.DialectV13).Descriptor().TableName() + "_2024_01",
		msg.DBReflect(pgdb_v1.DialectV13).Descriptor().TableName() + "_2024_02",
		msg.DBReflect(pgdb_v1.DialectV13).Descriptor().TableName() + "_2024_03",
	}

	for i, expected := range expectedNames {
		require.Equal(t, expected, subTables[i], "Partition table name mismatch")
	}

	// Test data insertion into partitions
	testData := []*CheeseIngredient{
		CheeseIngredient_builder{
			TenantId:   "t1",
			Id:         "pi1",
			EventId:    generateKSUIDForTime(time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)), // Jan 15, 2024
			SourceAddr: "127.0.0.1",
		}.Build(),
		CheeseIngredient_builder{
			TenantId:   "t2",
			Id:         "pi2",
			EventId:    generateKSUIDForTime(time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)), // Feb 15, 2024
			SourceAddr: "1.2.3.4",
		}.Build(),
		CheeseIngredient_builder{
			TenantId:   "t3",
			Id:         "pi3",
			EventId:    generateKSUIDForTime(time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)), // Mar 15, 2024
			SourceAddr: "2001:db8:abcd:12::1",
		}.Build(),
	}

	// Insert test data
	for _, data := range testData {
		sql, args, err := pgdb_v1.Insert(data, pgdb_v1.DialectV13)
		require.NoError(t, err)
		_, err = pg.DB.Exec(ctx, sql, args...)
		require.NoError(t, err)
	}

	// Verify data distribution across partitions
	for i, subTable := range subTables {
		var count int
		err := pg.DB.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", subTable)).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count, "Each partition should have exactly one row")

		var eventId string
		err = pg.DB.QueryRow(ctx, fmt.Sprintf("SELECT pb$event_id FROM %s", subTable)).Scan(&eventId)
		require.NoError(t, err)
		require.Equal(t, testData[i].GetEventId(), eventId, "Data should be in correct monthly partition")
	}
}

// generateKSUIDForTime creates a KSUID string for a given time
func generateKSUIDForTime(t time.Time) string {
	// Create a KSUID with the given timestamp
	id, _ := ksuid.NewRandomWithTime(t)
	return id.String()
}

func TestKSUIDCollation(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	// Create a message with KSUID partitioning
	msg := CheeseIngredient_builder{
		TenantId: "t1",
		Id:       "pi1",
		EventId:  generateKSUIDForTime(time.Now()),
	}.Build()

	// Create the table
	schema, err := pgdb_v1.CreateSchema(msg, pgdb_v1.DialectV13)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "failed to execute sql: '\n%s\n'", line)
	}

	// Verify the column collation
	var collation string
	err = pg.DB.QueryRow(ctx, `
		SELECT a.attcollation::regcollation::text
		FROM pg_attribute a
		JOIN pg_class c ON c.oid = a.attrelid
		WHERE c.relname = $1 AND a.attname = $2
	`, msg.DBReflect(pgdb_v1.DialectV13).Descriptor().TableName(), "pb$event_id").Scan(&collation)
	require.NoError(t, err)
	require.Equal(t, "\"C\"", collation, "KSUID column should use C collation")

	startDate := time.Date(2025, 3, 20, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 4, 21, 0, 0, 0, 0, time.UTC)

	// Test data ordering
	testData := []*CheeseIngredient{
		CheeseIngredient_builder{
			TenantId: "t1",
			Id:       "pi1",
			EventId:  generateKSUIDForTime(startDate.Add(time.Hour)), // Earlier timestamp
		}.Build(),
		CheeseIngredient_builder{
			TenantId: "t1",
			Id:       "pi2",
			EventId:  generateKSUIDForTime(endDate.Add(-1 * time.Hour)), // Later timestamp
		}.Build(),
	}

	// Create the partitions
	err = pgdb_v1.EventIDPartitionsUpdate(ctx, pg.DB, msg, pgdb_v1.DialectV13, startDate, endDate, func(ctx context.Context, schema string, args ...interface{}) error {
		_, err := pg.DB.Exec(ctx, schema)
		return err
	})
	require.NoError(t, err)

	// Insert test data
	for _, data := range testData {
		sql, args, err := pgdb_v1.Insert(data, pgdb_v1.DialectV13)
		require.NoError(t, err)
		_, err = pg.DB.Exec(ctx, sql, args...)
		require.NoError(t, err)
	}

	// Verify ordering
	rows, err := pg.DB.Query(ctx, fmt.Sprintf(`
		SELECT pb$event_id 
		FROM %s 
		ORDER BY pb$event_id
	`, msg.DBReflect(pgdb_v1.DialectV13).Descriptor().TableName()))
	require.NoError(t, err)
	defer rows.Close()

	var eventIds []string
	for rows.Next() {
		var eventId string
		err = rows.Scan(&eventId)
		require.NoError(t, err)
		eventIds = append(eventIds, eventId)
	}

	// Verify the order matches the chronological order of the KSUIDs
	require.Equal(t, testData[0].GetEventId(), eventIds[0], "First KSUID should be earlier timestamp")
	require.Equal(t, testData[1].GetEventId(), eventIds[1], "Second KSUID should be later timestamp")
}

func TestPastaIngredientBitVector(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin; CREATE EXTENSION IF NOT EXISTS vector;")
	require.NoError(t, err)

	// Create test data
	testData := make([]*PastaIngredient, 100)
	for i := 0; i < 100; i++ {
		// Create a bit vector with a pattern that varies by index
		minHash := make([]byte, 512) // 512 bytes -> 4096 bits to match the vector size on the message/model.
		for j := 0; j < len(minHash); j++ {
			// Set bits based on the index to create different patterns
			minHash[j] = byte((i + j) % 256)
		}

		testData[i] = PastaIngredient_builder{
			TenantId:     "t1",
			Id:           fmt.Sprintf("pi%d", i),
			PastaId:      fmt.Sprintf("p%d", i),
			IngredientId: fmt.Sprintf("i%d", i),
			MinHash:      minHash,
		}.Build()
	}

	// Create schema
	schema, err := pgdb_v1.CreateSchema(testData[0], pgdb_v1.DialectV13)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "failed to execute sql: '\n%s\n'", line)
	}

	// Create partitions for tenant_id
	fakeTenantIds := []string{"t1"}
	testCreatePartitionTables(t, pg, testData[0], fakeTenantIds)

	// Insert test data
	for _, data := range testData {
		sql, args, err := pgdb_v1.Insert(data, pgdb_v1.DialectV13)
		require.NoError(t, err)
		_, err = pg.DB.Exec(ctx, sql, args...)
		require.NoError(t, err)
	}

	// Test query using Distance()
	// Use the first item's min_hash as the reference point
	referenceMinHashBytes := testData[0].GetMinHash()

	// Build query to find similar items using raw SQL
	pi := testData[0]
	piFields := pi.DB().Query()

	qb := goqu.Dialect("postgres")
	query, args, err := qb.Select(
		piFields.Id().Identifier(),
		piFields.MinHash().Distance(referenceMinHashBytes),
	).From(pi.DB().TableName()).
		ToSQL()
	require.NoError(t, err)

	// Execute query
	rows, err := pg.DB.Query(ctx, query, args...)
	require.NoError(t, err)
	defer rows.Close()

	// Verify results
	var results []struct {
		ID       string
		Distance int
	}
	for rows.Next() {
		var result struct {
			ID       string
			Distance int
		}
		err := rows.Scan(&result.ID, &result.Distance)
		require.NoError(t, err)
		results = append(results, result)
	}
	require.NoError(t, rows.Err())

	// Verify we got some results
	require.Greater(t, len(results), 0, "Should find at least one similar item")

	// Verify the first result is the reference item (distance = 0)
	require.Equal(t, testData[0].GetId(), results[0].ID, "First result should be the reference item")
	require.Equal(t, 0, results[0].Distance, "First result should be the reference item")
}

func TestPastaIngredientBitVectorRetrieval(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin; CREATE EXTENSION IF NOT EXISTS vector;")
	require.NoError(t, err)

	// Create test data with a specific pattern in the min hash
	minHash := make([]byte, 512)
	for i := 0; i < len(minHash); i++ {
		// Create a pattern where each byte alternates between 0xAA and 0x55
		if i%2 == 0 {
			minHash[i] = 0xAA // 10101010
		} else {
			minHash[i] = 0x55 // 01010101
		}
	}

	// fmt.Println(">>> minHash:")
	// fmt.Println(minHash)

	testData := PastaIngredient_builder{
		TenantId:     "t1",
		Id:           "pi1",
		PastaId:      "p1",
		IngredientId: "i1",
		MinHash:      minHash,
	}.Build()

	// Create schema
	schema, err := pgdb_v1.CreateSchema(testData, pgdb_v1.DialectV13)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "failed to execute sql: '\n%s\n'", line)
	}

	// fmt.Println(">>> schema:")
	// fmt.Println(schema)

	// Create partitions for tenant_id
	fakeTenantIds := []string{"t1"}
	testCreatePartitionTables(t, pg, testData, fakeTenantIds)

	// Insert test data
	sql, args, err := pgdb_v1.Insert(testData, pgdb_v1.DialectV13)
	require.NoError(t, err)
	_, err = pg.DB.Exec(ctx, sql, args...)
	require.NoError(t, err)

	// fmt.Println(">>> testData:")
	// fmt.Println(testData)

	// Query the record
	pi := testData
	piFields := pi.DB().Query()

	qb := goqu.Dialect("postgres")
	query, args, err := qb.Select(
		piFields.Id().Identifier(),
		piFields.MinHash().Identifier(),
	).From(pi.DB().TableName()).
		Where(piFields.Id().Eq(testData.GetId())).
		ToSQL()
	require.NoError(t, err)

	// Execute query
	rows, err := pg.DB.Query(ctx, query, args...)
	require.NoError(t, err)
	defer rows.Close()

	// Verify results
	require.True(t, rows.Next(), "Should find the inserted record")
	var retrievedID string
	var retrievedMinHashBits string
	err = rows.Scan(&retrievedID, &retrievedMinHashBits)
	require.NoError(t, err)
	require.NoError(t, rows.Err())

	// fmt.Println(">>> retrievedID:")
	// fmt.Println(retrievedID)
	// fmt.Println(">>> retrievedMinHashBits:")
	// fmt.Println(retrievedMinHashBits)

	// Verify the ID matches
	require.Equal(t, testData.GetId(), retrievedID, "Retrieved ID should match")
	require.Equal(t, len(minHash)*8, len(retrievedMinHashBits), "Retrieved min hash length should match")

	// Verify the min hash matches
	retrievedMinHash := bitStringToBytes(t, retrievedMinHashBits)
	// fmt.Println(">>> retrievedMinHash:")
	// fmt.Println(retrievedMinHash)

	require.Equal(t, len(minHash), len(retrievedMinHash), "Retrieved min hash length should match")
	for i := 0; i < len(minHash); i++ {
		require.Equal(t, minHash[i], retrievedMinHash[i],
			"Retrieved min hash byte at index %d should match", i)
	}

	// Verify no more rows
	require.False(t, rows.Next(), "Should not find any more records")
}

func bitStringToBytes(t *testing.T, bitString string) []byte {
	bytes := make([]byte, len(bitString)/8)
	for i := 0; i < len(bitString); i += 8 {
		byteValue := bitString[i : i+8]
		intValue, err := strconv.ParseInt(byteValue, 2, 64)
		require.NoError(t, err)
		bytes[i/8] = byte(intValue)
	}
	return bytes
}

func TestKSUIDCollationV17(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()
	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin; CREATE EXTENSION IF NOT EXISTS vector;")
	require.NoError(t, err)

	startDate := time.Date(2025, 3, 20, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 4, 21, 0, 0, 0, 0, time.UTC)

	// Create some messages with KSUID fields that are not used for partitioning
	testData := []*PastaIngredient{
		PastaIngredient_builder{
			TenantId:     "t1",
			Id:           "pi1",
			PastaId:      generateKSUIDForTime(startDate.Add(time.Hour)),
			IngredientId: generateKSUIDForTime(startDate.Add(time.Hour)),
		}.Build(),
		PastaIngredient_builder{
			TenantId:     "t1",
			Id:           "pi2",
			PastaId:      generateKSUIDForTime(endDate.Add(-1 * time.Hour)),
			IngredientId: generateKSUIDForTime(endDate.Add(-1 * time.Hour)),
		}.Build(),
	}

	smsg := testData[0]
	schema, err := pgdb_v1.CreateSchema(smsg, pgdb_v1.DialectV17)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "failed to execute sql: '\n%s\n'", line)
	}

	// Verify the column collations
	var collation1 string
	err = pg.DB.QueryRow(ctx, `
		SELECT a.attcollation::regcollation::text
		FROM pg_attribute a
		JOIN pg_class c ON c.oid = a.attrelid
		WHERE c.relname = $1 AND a.attname = $2
	`, smsg.DBReflect(pgdb_v1.DialectV17).Descriptor().TableName(), "pb$pasta_id").Scan(&collation1)
	require.NoError(t, err)
	require.Equal(t, "\"C\"", collation1, "KSUID column should use C collation")

	var collation2 string
	err = pg.DB.QueryRow(ctx, `
		SELECT a.attcollation::regcollation::text
		FROM pg_attribute a
		JOIN pg_class c ON c.oid = a.attrelid
		WHERE c.relname = $1 AND a.attname = $2
	`, smsg.DBReflect(pgdb_v1.DialectV17).Descriptor().TableName(), "pb$ingredient_id").Scan(&collation2)
	require.NoError(t, err)
	require.Equal(t, "\"C\"", collation2, "KSUID column should use C collation")

	// Insert test data
	for _, msg := range testData {
		sql, args, err := pgdb_v1.Insert(msg, pgdb_v1.DialectV17)
		require.NoError(t, err)
		_, err = pg.DB.Exec(ctx, sql, args...)
		require.NoError(t, err)
	}

	// Verify ordering
	rows1, err := pg.DB.Query(ctx, fmt.Sprintf(`
		SELECT pb$pasta_id 
		FROM %s 
		ORDER BY pb$pasta_id
	`, smsg.DBReflect(pgdb_v1.DialectV17).Descriptor().TableName()))
	require.NoError(t, err)
	defer rows1.Close()

	var pastaIDs []string
	for rows1.Next() {
		var pastaID string
		err = rows1.Scan(&pastaID)
		require.NoError(t, err)
		pastaIDs = append(pastaIDs, pastaID)
	}

	// Verify the order matches the chronological order of the KSUIDs
	require.Equal(t, testData[0].GetPastaId(), pastaIDs[0], "First KSUID should be earlier timestamp")
	require.Equal(t, testData[1].GetPastaId(), pastaIDs[1], "Second KSUID should be later timestamp")

	// Verify ordering
	rows2, err := pg.DB.Query(ctx, fmt.Sprintf(`
		SELECT pb$ingredient_id
		FROM %s 
		ORDER BY pb$ingredient_id
	`, smsg.DBReflect(pgdb_v1.DialectV17).Descriptor().TableName()))
	require.NoError(t, err)
	defer rows2.Close()

	var ingredientIDs []string
	for rows2.Next() {
		var ingredientID string
		err = rows2.Scan(&ingredientID)
		require.NoError(t, err)
		ingredientIDs = append(ingredientIDs, ingredientID)
	}

	// Verify the order matches the chronological order of the KSUIDs
	require.Equal(t, testData[0].GetIngredientId(), ingredientIDs[0], "First KSUID should be earlier timestamp")
	require.Equal(t, testData[1].GetIngredientId(), ingredientIDs[1], "Second KSUID should be later timestamp")
}
