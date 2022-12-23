package v1

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/ductone/protoc-gen-pgdb/internal/pgtest"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	insertMsg := &Pet{
		TenantId:      "t1",
		Id:            "obj2",
		CreatedAt:     timestamppb.Now(),
		UpdatedAt:     timestamppb.Now(),
		DisplayName:   "Lion",
		Description:   "the coolest pet, a lion",
		SystemBuiltin: false,
		Elapsed:       durationpb.New(time.Hour),
		Profile:       &structpb.Struct{},
		Cuteness:      1.0,
		Price:         9000.0,
	}

	dbr := insertMsg.DBReflect()
	tableName := dbr.Descriptor().TableName()

	record, err := dbr.Record()
	require.NoError(t, err)

	qb := goqu.Dialect("postgres")

	q := qb.Insert(tableName).Prepared(true).Rows(
		record,
	)

	if false {
		q = q.OnConflict(goqu.DoUpdate(`"pb$tenant_id", "pb$pk", "pb$sk"`,
			exp.Record{
				"pb$tenant_id": exp.NewIdentifierExpression("", "EXCLUDED", "pb$tenant_id"),
				"pb$pk":        exp.NewIdentifierExpression("", "EXCLUDED", "pb$pk"),
				// "sk":         goqu.L("EXCLUDED.sk"),
				// "gsi1pk":     goqu.L("EXCLUDED.gsi1pk"),
				// "gsi1sk":     goqu.L("EXCLUDED.gsi1sk"),
				// "created_at": goqu.L("EXCLUDED.created_at"),
				// "updated_at": goqu.L("EXCLUDED.updated_at"),
				// "deleted_at": goqu.L("EXCLUDED.deleted_at"),
				// "search":     goqu.L("EXCLUDED.search"),
				// "data":       goqu.L("EXCLUDED.data"),
			},
		).Where(exp.NewIdentifierExpression("", "EXCLUDED", "pb$updated_at").Gte(insertMsg.DB().Columns().UpdatedAt())))
	}

	query, params, err := q.ToSQL()
	require.NoError(t, err)

	_, err = pg.DB.Exec(ctx, query, params...)
	spew.Dump(query, params)
	fmt.Fprintf(os.Stderr, "---------\n%s\n\n", query)
	require.NoError(t, err)
}
