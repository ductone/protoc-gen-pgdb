package v1

import (
	"context"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/stretchr/testify/require"

	animals_v1 "github.com/ductone/protoc-gen-pgdb/example/models/animals/v1"
	zoo_v1 "github.com/ductone/protoc-gen-pgdb/example/models/zoo/v1"
	"github.com/ductone/protoc-gen-pgdb/internal/pgtest"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
)

func TestSchemaAttractions(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	schema, err := pgdb_v1.CreateSchema(&Attractions{}, pgdb_v1.DialectV13)
	require.NoError(t, err)
	for _, line := range schema {

		// fmt.Printf("%s \n", line)
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "TestSchemaPet: failed to execute sql: '\n%s\n'", line)
	}
}

func TestNestedIndexes(t *testing.T) {
	fields := (*Attractions)(nil).DB().Query()

	pk := fields.PK()
	require.Equal(t, "pb$pk",
		pk.column,
		"bad field resolution for pk: %s", pk.column,
	)
	pksk := fields.PKSK()
	require.Equal(t, "pb$pksk",
		pksk.column,
		"bad field resolution for pksk: %s", pksk.column,
	)

	ftsData := fields.FTSData()
	require.Equal(t, "pb$fts_data",
		ftsData.column,
		"bad field resolution for ftsData: %s", ftsData.column,
	)

	mediumMedium := fields.MediumMedium()
	require.Equal(t, "pb$12$medium_oneof",
		mediumMedium.column,
		"bad resolution for medium medium: %s", mediumMedium.column,
	)

	zooshopsfixed := fields.ZooShopAnythingSfixed64()
	require.Equal(t, "pb$11$52$sfixed_64",
		zooshopsfixed.column,
		"bad resolution for zooshopsfixed: %s", zooshopsfixed.column,
	)

	zooShopMedium := fields.ZooShopMedium()
	require.Equal(t, "pb$11$medium_oneof",
		zooShopMedium.column,
		"bad resolution for zooShopMedium: %s", zooShopMedium.column,
	)

	numId := fields.Unsafe().Numid()
	require.Equal(t, "pb$numid",
		numId.column,
		"bad resolution for numid: %s", numId.column,
	)

	// Test accessing non-indexed nested fields through the Unsafe accessor pattern
	// zoo_shop.anything.id is NOT indexed - use UnsafeAnything() to access it
	// (only zoo_shop.anything.sfixed_64 is indexed in city.proto)
	zooShopAnythingId := fields.ZooShop().UnsafeAnything().Id()
	require.Equal(t, "pb$11$52$id",
		zooShopAnythingId.column,
		"bad resolution for zoo_shop.anything.id via UnsafeAnything(): %s", zooShopAnythingId.column,
	)

	// zoo_shop.mgr.id is a deeply nested non-indexed field - use UnsafeMgr() to access
	zooShopMgrId := fields.ZooShop().UnsafeMgr().Id()
	require.Equal(t, "pb$11$5$id",
		zooShopMgrId.column,
		"bad resolution for zoo_shop.mgr.id via UnsafeMgr(): %s", zooShopMgrId.column,
	)

	// Verify the indexed path still works (for comparison)
	// zoo_shop.anything.sfixed_64 IS indexed, so it's accessible via the regular Anything() accessor
	zooShopAnythingSfixed := fields.ZooShop().Anything().Sfixed64()
	require.Equal(t, "pb$11$52$sfixed_64",
		zooShopAnythingSfixed.column,
		"bad resolution for indexed zoo_shop.anything.sfixed_64: %s", zooShopAnythingSfixed.column,
	)
}

// TestNestedColumnMetadata verifies that nested columns have correct metadata,
// including the parent's table name and full proto paths.
func TestNestedColumnMetadata(t *testing.T) {
	desc := (*Attractions)(nil).DBReflect(pgdb_v1.DialectV13).Descriptor()
	expectedTableName := desc.TableName()
	columns := desc.Fields()

	// Find columns by their ProtoPath to verify correct metadata
	var tenantCol, zooShopFurCol, zooShopAnythingSfixed64Col *pgdb_v1.Column
	for _, col := range columns {
		switch col.ProtoPath {
		case "tenant_id":
			tenantCol = col
		case "zoo_shop.fur":
			zooShopFurCol = col
		case "zoo_shop.anything.sfixed_64":
			zooShopAnythingSfixed64Col = col
		}
	}

	// Test direct field
	require.NotNil(t, tenantCol, "should find tenant_id column")
	require.Equal(t, expectedTableName, tenantCol.Table, "tenant_id should have Attractions table name")
	require.Equal(t, "pb$tenant_id", tenantCol.Name, "tenant_id column name")
	require.Equal(t, []int32{1}, tenantCol.ProtoFieldPath, "tenant_id proto field path")

	// Test nested field (Shop.fur) - should have Attractions table, not Shop table
	require.NotNil(t, zooShopFurCol, "should find zoo_shop.fur column")
	require.Equal(t, expectedTableName, zooShopFurCol.Table,
		"nested fur column should have Attractions table name, not Shop table name")
	require.Equal(t, "pb$11$fur", zooShopFurCol.Name, "zoo_shop.fur column name with nested prefix")
	require.Equal(t, []int32{11, 4}, zooShopFurCol.ProtoFieldPath, "zoo_shop.fur proto field path")
	require.Equal(t, ".models.animals.v1.FurType", zooShopFurCol.ProtoTypeName,
		"zoo_shop.fur should have enum type name")

	// Test deeply nested field (Shop.anything.sfixed_64)
	require.NotNil(t, zooShopAnythingSfixed64Col, "should find zoo_shop.anything.sfixed_64 column")
	require.Equal(t, expectedTableName, zooShopAnythingSfixed64Col.Table,
		"deeply nested column should have Attractions table name")
	require.Equal(t, "pb$11$52$sfixed_64", zooShopAnythingSfixed64Col.Name,
		"deeply nested column name with nested prefixes")
	require.Equal(t, []int32{11, 52, 14}, zooShopAnythingSfixed64Col.ProtoFieldPath,
		"zoo_shop.anything.sfixed_64 proto field path (field 52=anything, 14=sfixed_64)")
}

// TestColumnIdentifier verifies that Column.Identifier() returns correct goqu expression.
func TestColumnIdentifier(t *testing.T) {
	desc := (*Attractions)(nil).DBReflect(pgdb_v1.DialectV13).Descriptor()
	expectedTableName := desc.TableName()
	columns := desc.Fields()

	// Find a nested column
	var zooShopFurCol *pgdb_v1.Column
	for _, col := range columns {
		if col.ProtoPath == "zoo_shop.fur" {
			zooShopFurCol = col
			break
		}
	}

	require.NotNil(t, zooShopFurCol, "should find zoo_shop.fur column")

	// Verify Identifier() returns correct expression
	ident := zooShopFurCol.Identifier()
	require.Equal(t, expectedTableName, ident.GetTable(),
		"Identifier() should return parent table name for nested columns")
	require.Equal(t, "pb$11$fur", ident.GetCol().(string),
		"Identifier() should return correct column name")
}

// TestNestedAccessorSQLGeneration verifies that nested query builder accessors
// generate correct SQL when used with goqu.
func TestNestedAccessorSQLGeneration(t *testing.T) {
	fields := (*Attractions)(nil).DB().Query()
	tableName := (*Attractions)(nil).DB().TableName()
	qb := goqu.Dialect("postgres")

	tests := []struct {
		name        string
		expr        exp.Expression
		mustContain string
	}{
		{
			name:        "indexed nested field Eq",
			expr:        fields.ZooShop().Anything().Sfixed64().Eq(int64(42)),
			mustContain: `"pb$11$52$sfixed_64" = 42`,
		},
		{
			name:        "unsafe nested field In",
			expr:        fields.ZooShop().UnsafeAnything().Id().In([]string{"a", "b"}),
			mustContain: `"pb$11$52$id" IN ('a', 'b')`,
		},
		{
			name:        "deeply nested manager field",
			expr:        fields.ZooShop().UnsafeMgr().Id().Eq(int32(123)),
			mustContain: `"pb$11$5$id" = 123`,
		},
		{
			name:        "nested time.Time comparison",
			expr:        fields.ZooShop().UnsafeCreatedAt().IsNotNull(),
			mustContain: `"pb$11$created_at" IS NOT NULL`,
		},
		{
			name:        "nested Between range",
			expr:        fields.ZooShop().UnsafeAnything().Double().Between(1.0, 100.0),
			mustContain: `"pb$11$52$double" BETWEEN`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, _, err := qb.Select(goqu.L("1")).From(tableName).Where(tt.expr).ToSQL()
			require.NoError(t, err, "SQL generation should succeed")
			require.Contains(t, sql, tt.mustContain, "SQL should contain expected clause")
		})
	}

}

// TestNestedAccessorQueryExecution is an integration test that verifies
// queries using nested accessors execute correctly against PostgreSQL.
func TestNestedAccessorQueryExecution(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	// Create schema
	schema, err := pgdb_v1.CreateSchema(&Attractions{}, pgdb_v1.DialectV13)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoError(t, err)
	}

	// Insert test data with nested message
	testAttraction := Attractions_builder{
		TenantId: "tenant1",
		Id:       "attr1",
		Numid:    42,
		ZooShop: zoo_v1.Shop_builder{
			TenantId: "tenant1",
			Id:       "shop1",
			Anything: animals_v1.ScalarValue_builder{
				TenantId: "tenant1",
				Id:       "scalar1",
				Sfixed64: 999,
				Double:   3.14,
			}.Build(),
			Mgr: zoo_v1.Shop_Manager_builder{Id: 123}.Build(),
		}.Build(),
	}.Build()

	insertSQL, insertArgs, err := pgdb_v1.Insert(testAttraction, pgdb_v1.DialectV13)
	require.NoError(t, err)
	_, err = pg.DB.Exec(ctx, insertSQL, insertArgs...)
	require.NoError(t, err)

	// Query using nested accessors
	fields := (*Attractions)(nil).DB().Query()
	tableName := (*Attractions)(nil).DB().TableName()
	qb := goqu.Dialect("postgres")

	// Test 1: Query indexed nested field
	query, args, err := qb.Select(
		fields.TenantId().Identifier(),
		fields.ZooShop().Anything().Sfixed64().Identifier(),
	).From(tableName).
		Where(fields.ZooShop().Anything().Sfixed64().Eq(int64(999))).
		ToSQL()
	require.NoError(t, err)

	rows, err := pg.DB.Query(ctx, query, args...)
	require.NoError(t, err)
	defer rows.Close()

	require.True(t, rows.Next(), "should find one row")
	var tenantId string
	var sfixed64 int64
	err = rows.Scan(&tenantId, &sfixed64)
	require.NoError(t, err)
	require.Equal(t, "tenant1", tenantId)
	require.Equal(t, int64(999), sfixed64)
	rows.Close()

	// Test 2: Query unsafe nested field (manager id)
	query2, args2, err := qb.Select(
		fields.ZooShop().UnsafeMgr().Id().Identifier(),
	).From(tableName).
		Where(fields.ZooShop().UnsafeMgr().Id().Eq(int32(123))).
		ToSQL()
	require.NoError(t, err)

	rows2, err := pg.DB.Query(ctx, query2, args2...)
	require.NoError(t, err)
	defer rows2.Close()

	require.True(t, rows2.Next(), "should find manager by id")
	var mgrId int32
	err = rows2.Scan(&mgrId)
	require.NoError(t, err)
	require.Equal(t, int32(123), mgrId)
}

// TestNestedQueryBuilderTableNamePropagation verifies that table names propagate
// correctly through chained nested query builder accessors.
func TestNestedQueryBuilderTableNamePropagation(t *testing.T) {
	// Use WithTable to set a custom table name
	customTable := "my_custom_attractions_table"
	fields := (*Attractions)(nil).DB().Query().WithTable(customTable)

	// Verify table name propagates through chained accessors
	sfixedOps := fields.ZooShop().Anything().Sfixed64()
	require.Equal(t, customTable, sfixedOps.tableName, "table name should propagate to indexed nested field")

	unsafeIdOps := fields.ZooShop().UnsafeAnything().Id()
	require.Equal(t, customTable, unsafeIdOps.tableName, "table name should propagate to unsafe nested field")

	mgrIdOps := fields.ZooShop().UnsafeMgr().Id()
	require.Equal(t, customTable, mgrIdOps.tableName, "table name should propagate to deeply nested field")
}

// TestNestedAccessorEmptySliceHandling verifies that empty slices passed to
// In() and NotIn() operators generate correct SQL (FALSE and TRUE respectively).
func TestNestedAccessorEmptySliceHandling(t *testing.T) {
	fields := (*Attractions)(nil).DB().Query()
	qb := goqu.Dialect("postgres")
	tableName := (*Attractions)(nil).DB().TableName()

	// Empty In() should return FALSE (no matches possible)
	emptyInExpr := fields.ZooShop().UnsafeAnything().Id().In([]string{})
	sql, _, err := qb.Select(goqu.L("1")).From(tableName).Where(emptyInExpr).ToSQL()
	require.NoError(t, err)
	require.Contains(t, sql, "FALSE", "empty In() should generate FALSE clause")

	// Empty NotIn() should return TRUE (all rows match)
	emptyNotInExpr := fields.ZooShop().UnsafeAnything().Id().NotIn([]string{})
	sql2, _, err := qb.Select(goqu.L("1")).From(tableName).Where(emptyNotInExpr).ToSQL()
	require.NoError(t, err)
	require.Contains(t, sql2, "TRUE", "empty NotIn() should generate TRUE clause")
}

// TestNestedOnlyNoStandaloneQueryBuilder verifies that nested_only messages
// do not generate standalone query builders, preventing type name collisions.
//
// Without the fix, the combination of:
//   - AttractionsConfig (nested_only message)
//   - Attractions.config field of type AttractionsConfig
//
// Would cause duplicate type definitions because:
//   - "Attractions" + "Config" = "AttractionsConfig" (parent + field name)
//   - "AttractionsConfig" (the message type name)
//
// Both would generate "AttractionsConfigDetailSafeOperators" and similar types.
// This test passes if the code compiles - duplicate types would cause compilation failure.
func TestNestedOnlyNoStandaloneQueryBuilder(t *testing.T) {
	// Access nested config through Attractions - this should work
	fields := (*Attractions)(nil).DB().Query()

	// Verify we can access nested fields through the parent
	// These use UnsafeDetail() and UnsafeName() because AttractionsConfig fields are not indexed
	configDetail := fields.Config().UnsafeDetail()
	require.NotNil(t, configDetail, "should be able to access nested_only message fields")

	configName := fields.Config().UnsafeName()
	require.NotNil(t, configName, "should be able to access nested_only message fields")

	// Test AttractionsV2 can also access the same nested_only message type
	// This verifies the same nested_only type can be embedded in multiple parents
	fieldsV2 := (*AttractionsV2)(nil).DB().Query()

	configDetailV2 := fieldsV2.Config().UnsafeDetail()
	require.NotNil(t, configDetailV2, "should be able to access nested_only message from second parent")

	configNameV2 := fieldsV2.Config().UnsafeName()
	require.NotNil(t, configNameV2, "should be able to access nested_only message from second parent")

	// Verify we can access the actual Info field (which is the leaf field in Detail message)
	// This demonstrates the full path access works without type collisions
	infoOps := fields.Config().UnsafeDetail().Info()
	require.NotNil(t, infoOps, "should be able to access leaf field through nested_only message")

	infoOpsV2 := fieldsV2.Config().UnsafeDetail().Info()
	require.NotNil(t, infoOpsV2, "should be able to access leaf field through nested_only message from V2")

	// The test passes if this compiles - duplicate types would cause compilation failure
}

// TestNestedOnlyWithOneofThroughMultipleLayers verifies that nested_only messages
// with oneofs can be accessed through multiple layers of non-nested_only messages.
// This tests the fix for incorrect ParentTypeName in unsafe query builder children.
func TestNestedOnlyWithOneofThroughMultipleLayers(t *testing.T) {
	// Access from outer wrapper through middle to the oneof
	wrapperFields := (*NestedOnlyWrapper)(nil).DB().Query()

	// Navigate: Wrapper -> Middle -> OneofField -> oneof cases
	// Middle() returns NestedOnlyWrapperMiddleQueryBuilder
	middleQB := wrapperFields.Middle()
	require.NotNil(t, middleQB, "should access middle query builder")

	// UnsafeOneofField() returns NestedOnlyWrapperMiddleOneofFieldUnsafeQueryBuilder
	oneofQB := middleQB.UnsafeOneofField()
	require.NotNil(t, oneofQB, "should access oneof field via unsafe accessor")

	// UnsafeChoiceA() should be on the Unsafe type, not on a non-existent safe type
	choiceAQB := oneofQB.UnsafeChoiceA()
	require.NotNil(t, choiceAQB, "should access oneof case A")

	// Access leaf field
	valueAOps := choiceAQB.ValueA()
	require.NotNil(t, valueAOps, "should access leaf field in oneof case")

	// Similarly for choice B
	choiceBQB := oneofQB.UnsafeChoiceB()
	require.NotNil(t, choiceBQB, "should access oneof case B")
	valueBOps := choiceBQB.ValueB()
	require.NotNil(t, valueBOps, "should access leaf field in oneof case B")

	// This test passes if it compiles - the bug causes undefined type errors
}

// TestDuplicateTypeBug verifies that nested query builders don't generate
// duplicate type definitions when accessible from both safe and unsafe paths.
// Bug scenario: when a nested message has both indexed fields (creating safe builder)
// and non-indexed content (creating unsafe builder), child types should only be defined once.
func TestDuplicateTypeBug(t *testing.T) {
	fields := (*DuplicateTypeBugOuter)(nil).DB().Query()

	// Access via safe path (indexed field exists)
	nestedQB := fields.Nested()
	require.NotNil(t, nestedQB, "should have safe nested query builder")

	// Access indexed field directly
	indexedOps := nestedQB.IndexedField()
	require.NotNil(t, indexedOps, "should access indexed field from safe builder")

	// Access unindexed field via Unsafe accessor
	unindexedOps := nestedQB.UnsafeUnindexedField()
	require.NotNil(t, unindexedOps, "should access unindexed field via unsafe accessor")

	// Access deeper nested via Unsafe (from safe builder)
	innerFromSafe := nestedQB.UnsafeInner()
	require.NotNil(t, innerFromSafe, "should access inner from safe builder")
	innerValueFromSafe := innerFromSafe.InnerValue()
	require.NotNil(t, innerValueFromSafe, "should access inner value from safe path")

	// Note: When there's an indexed field, the Nested() method returns a safe builder
	// and UnsafeNested() is not generated. The child types (like InnerUnsafeQueryBuilder)
	// are accessible from both the safe builder's UnsafeChildren path and would be
	// duplicated if we also generated them for the unsafe path.
	// This test verifies that the code compiles - duplicate types would cause failure.
}

// TestEmbeddedWithOwnDBNoCollision verifies that embedding a message with its own DB
// does NOT cause duplicate type declarations. When message B (with its own DB) is
// embedded in message A, A should NOT generate types for B's children because B
// will generate those types itself.
func TestEmbeddedWithOwnDBNoCollision(t *testing.T) {
	// This test passes if it compiles - duplicate types would cause compilation failure.

	// Access from parent - should have accessor to embedded but NOT its children
	parentFields := (*ParentWithEmbeddedDB)(nil).DB().Query()
	require.NotNil(t, parentFields, "parent should have query builder")

	// Access from embedded's own query builder
	embeddedFields := (*EmbeddedWithOwnDB)(nil).DB().Query()
	require.NotNil(t, embeddedFields, "embedded should have its own query builder")

	// The inner nested_only field should be accessible from EmbeddedWithOwnDB
	innerQB := embeddedFields.Inner()
	require.NotNil(t, innerQB, "embedded should have accessor to its nested_only child")

	// Access the value field from the inner nested message (non-indexed, so Unsafe prefix)
	valueOps := innerQB.UnsafeValue()
	require.NotNil(t, valueOps, "should access value field from inner nested message")
}
