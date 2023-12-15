package v1

import (
	"context"
	"testing"

	"github.com/ductone/protoc-gen-pgdb/internal/pgtest"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	"github.com/stretchr/testify/require"
)

func TestSchemaAttractions(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	schema, err := pgdb_v1.CreateSchema(&Attractions{})
	require.NoError(t, err)
	for _, line := range schema {

		//fmt.Printf("%s \n", line)
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
		"bad resolution for medium medium: %s", numId.column,
	)
}
