package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

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
