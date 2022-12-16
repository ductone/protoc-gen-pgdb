package v1

import (
	"context"
	"testing"

	"github.com/ductone/protoc-gen-pgdb/internal/pgtest"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	"github.com/stretchr/testify/require"
)

func TestSchema(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	msg := &Pet{}
	schema, err := pgdb_v1.CreateSchema(msg)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		//		fmt.Printf("%s\n", line)
		require.NoErrorf(t, err, "TestCreateSchema: failed to execute sql: '\n%s\n'", line)
	}

	// make sure we should have zero migrations after schema create
	m, err := pgdb_v1.Migrations(ctx, pg.DB, msg)
	require.NoError(t, err)
	require.Len(t, m, 0)

	// drop profile for fun
	// it is both a col, and an index!
	_, err = pg.DB.Exec(ctx, `ALTER TABLE pb_pet_models_animals_v1_8a3723d5 DROP COLUMN "pb$profile"`)
	require.NoError(t, err)
	m, err = pgdb_v1.Migrations(ctx, pg.DB, msg)
	require.NoError(t, err)

	require.Len(t, m, 2)
	// fmt.Printf("-----\n%s\n", m[0])
	require.Contains(t, m[0], "ALTER TABLE")
	require.Contains(t, m[0], "pb$profile")
	// fmt.Printf("-----\n%s\n", m[1])
	require.Contains(t, m[1], "CREATE INDEX CONCURRENTLY IF NOT EXISTS")
	require.Contains(t, m[1], "pb$profile")
}
