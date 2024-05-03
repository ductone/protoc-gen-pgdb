package v1

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ductone/protoc-gen-pgdb/internal/pgtest"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
)

func TestSchemaPet(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	schema, err := pgdb_v1.CreateSchema(&Pet{})
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "TestSchemaPet: failed to execute sql: '\n%s\n'", line)
	}

	// make sure we should have zero migrations after schema create
	m, err := pgdb_v1.Migrations(ctx, pg.DB, &Pet{})
	require.NoError(t, err)
	require.Len(t, m, 0)

	// drop profile for fun
	// it is both a col, and an index!
	_, err = pg.DB.Exec(ctx, `ALTER TABLE pb_pet_models_animals_v1_8a3723d5 DROP COLUMN "pb$profile"`)
	require.NoError(t, err)
	migrations, err := pgdb_v1.Migrations(ctx, pg.DB, &Pet{})
	require.NoError(t, err)

	require.Len(t, migrations, 2)
	// fmt.Printf("-----\n%s\n", m[0])
	require.Contains(t, migrations[0], "ALTER TABLE")
	require.Contains(t, migrations[0], "pb$profile")
	// fmt.Printf("-----\n%s\n", m[1])
	require.Contains(t, migrations[1], "CREATE INDEX CONCURRENTLY IF NOT EXISTS")
	require.Contains(t, migrations[1], "pb$profile")

	for _, line := range migrations {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "TestSchemaPet: failed to execute migration sql: '\n%s\n'", line)
	}

	insertMsg := &Pet{
		TenantId:      "t1",
		Id:            "obj2",
		CreatedAt:     timestamppb.Now(),
		UpdatedAt:     timestamppb.Now(),
		DisplayName:   "Lion zoo:animal",
		Description:   "the coolest pet, a lion",
		SystemBuiltin: false,
		Elapsed:       durationpb.New(time.Hour),
		Profile:       &structpb.Struct{},
		Cuteness:      1.0,
		Price:         9000.0,
		ExtraProfiles: []*structpb.Struct{
			{
				Fields: map[string]*structpb.Value{
					"foo": {
						Kind: &structpb.Value_BoolValue{
							BoolValue: true,
						},
					},
				},
			},
		},
	}

	query, params, err := pgdb_v1.Insert(insertMsg)
	require.NoError(t, err)
	_, err = pg.DB.Exec(ctx, query, params...)
	// spew.Dump(query, params)
	// spew.Dump(record)
	require.NoError(t, err, "query failed: %s\n\n%+v\n\n", query, params)

	insertMsg2 := &Pet{
		TenantId:      "t1",
		Id:            "obj3",
		CreatedAt:     timestamppb.Now(),
		UpdatedAt:     timestamppb.Now(),
		DisplayName:   "Tiger",
		Description:   "the coolest pet, a Tiger",
		SystemBuiltin: false,
		Elapsed:       durationpb.New(time.Hour),
		Profile:       &structpb.Struct{},
		Cuteness:      1.0,
		Price:         9000.0,
		ExtraProfiles: []*structpb.Struct{
			{
				Fields: map[string]*structpb.Value{
					"foo": {
						Kind: &structpb.Value_BoolValue{
							BoolValue: true,
						},
					},
				},
			},
		},
	}

	// qb := goqu.Dialect("postgres")
	// countQuery, params, err := qb.Select(goqu.COUNT(goqu.Star()).As("count")).From(insertMsg.DBReflect().Descriptor().TableName()).ToSQL()
	// require.NoError(t, err)

	query, params, err = pgdb_v1.Insert(insertMsg2)
	require.NoError(t, err)
	_, err = pg.DB.Exec(ctx, query, params...)
	require.NoError(t, err, "query failed: %s\n\n%+v\n\n", query, params)

	query, params, err = pgdb_v1.Delete(insertMsg2)
	require.NoError(t, err)
	res, err := pg.DB.Exec(ctx, query, params...)
	require.NoError(t, err, "query failed: %s\n\n%+v\n\n", query, params)
	require.Equal(t, int64(1), res.RowsAffected())
}

func TestSchemaBook(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	schema, err := pgdb_v1.CreateSchema(&Book{})
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		// fmt.Fprintf(os.Stderr, "---------\n%s\n\n", line)
		require.NoErrorf(t, err, "TestSchemaBook: failed to execute sql: '\n%s\n'", line)
	}

	query, params, err := pgdb_v1.Insert(&Book{
		TenantId: "t1",
		Id:       "b1",
		Medium: &Book_Ebook{
			Ebook: &EBook{Size: 4000},
		},
	})
	require.NoError(t, err)

	_, err = pg.DB.Exec(ctx, query, params...)
	// spew.Dump(query, params)
	// spew.Dump(record)
	// fmt.Fprintf(os.Stderr, "---------\n%s\n\n", query)
	require.NoError(t, err)
}

func TestSchemaScalarValue(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	schema, err := pgdb_v1.CreateSchema(&ScalarValue{})
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		// fmt.Fprintf(os.Stderr, "---------\n%s\n\n", line)
		require.NoErrorf(t, err, "TestSchemaScalarValue: failed to execute sql: '\n%s\n'", line)
	}

	query, params, err := pgdb_v1.Insert(&ScalarValue{})
	require.NoError(t, err)

	_, err = pg.DB.Exec(ctx, query, params...)
	require.NoError(t, err)

	query, params, err = pgdb_v1.Insert(&ScalarValue{
		TenantId:       "t1",
		Id:             "b1",
		RepeatedString: []string{"hello", "world"},
		RepeatedEnum:   []FurType{FurType_FUR_TYPE_LOTS, FurType_FUR_TYPE_NONE},
		RepeatedFloat:  []float32{1.2, 3.2},
		RepeatedBytes: [][]byte{
			[]byte("hi"),
			[]byte("mars"),
		},
	})
	require.NoError(t, err)

	_, err = pg.DB.Exec(ctx, query, params...)
	// spew.Dump(query, params)
	// spew.Dump(record)
	// fmt.Fprintf(os.Stderr, "---------\n%s\n\n", query)
	require.NoError(t, err, "invalid sql insert: %w %s", err, query)
}
