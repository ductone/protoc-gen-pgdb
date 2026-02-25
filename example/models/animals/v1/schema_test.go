package v1

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
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

	schema, err := pgdb_v1.CreateSchema(&Pet{}, pgdb_v1.DialectV13)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "TestSchemaPet: failed to execute sql: '\n%s\n'", line)
	}

	// make sure we should have zero migrations after schema create
	m, err := pgdb_v1.Migrations(ctx, pg.DB, &Pet{}, pgdb_v1.DialectV13)
	require.NoError(t, err)
	require.Len(t, m, 0)

	// drop profile for fun
	// it is both a col, and an index!
	_, err = pg.DB.Exec(ctx, `ALTER TABLE pb_pet_models_animals_v1_8a3723d5 DROP COLUMN "pb$profile"`)
	require.NoError(t, err)
	migrations, err := pgdb_v1.Migrations(ctx, pg.DB, &Pet{}, pgdb_v1.DialectV13)
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

	insertMsg := Pet_builder{
		TenantId:      "t1",
		Id:            "obj2",
		CreatedAt:     timestamppb.Now(),
		UpdatedAt:     timestamppb.Now(),
		DisplayName:   "\u0000Lion zoo:animal\u0000",
		Description:   "the coolest \u0000pet, a lion",
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
	}.Build()

	query, params, err := pgdb_v1.Insert(insertMsg, pgdb_v1.DialectV13)
	require.NoError(t, err)
	_, err = pg.DB.Exec(ctx, query, params...)
	// spew.Dump(query, params)
	// spew.Dump(record)
	require.NoError(t, err, "query failed: %s\n\n%+v\n\n", query, params)

	insertMsg2 := Pet_builder{
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
	}.Build()

	// qb := goqu.Dialect("postgres")
	// countQuery, params, err := qb.Select(goqu.COUNT(goqu.Star()).As("count")).From(insertMsg.DBReflect().Descriptor().TableName()).ToSQL()
	// require.NoError(t, err)

	query, params, err = pgdb_v1.Insert(insertMsg2, pgdb_v1.DialectV13)
	require.NoError(t, err)
	_, err = pg.DB.Exec(ctx, query, params...)
	require.NoError(t, err, "query failed: %s\n\n%+v\n\n", query, params)

	query, params, err = pgdb_v1.Delete(insertMsg2, pgdb_v1.DialectV13)
	require.NoError(t, err)
	res, err := pg.DB.Exec(ctx, query, params...)
	require.NoError(t, err, "query failed: %s\n\n%+v\n\n", query, params)
	require.Equal(t, int64(1), res.RowsAffected())

	var sb strings.Builder
	for i := 0; i < 2100; i++ {
		sb.WriteString("a")
	}

	result := sb.String()

	insertMsg3 := Pet_builder{
		TenantId:      "t1",
		Id:            "obj4",
		CreatedAt:     timestamppb.Now(),
		UpdatedAt:     timestamppb.Now(),
		DisplayName:   result,
		Description:   result,
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
	}.Build()

	query, params, err = pgdb_v1.Insert(insertMsg3, pgdb_v1.DialectV13)
	require.NoError(t, err)
	_, err = pg.DB.Exec(ctx, query, params...)
	require.NoError(t, err, "query failed: %s\n\n%+v\n\n", query, params)

	query, params, err = pgdb_v1.Delete(insertMsg3, pgdb_v1.DialectV13)
	require.NoError(t, err)
	res, err = pg.DB.Exec(ctx, query, params...)
	require.NoError(t, err, "query failed: %s\n\n%+v\n\n", query, params)
	require.Equal(t, int64(1), res.RowsAffected())

}

func TestMigrationIndexMutation(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	// Create initial schema
	schema, err := pgdb_v1.CreateSchema(&Pet{}, pgdb_v1.DialectV13)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "failed to execute sql: '\n%s\n'", line)
	}

	// Verify no migrations needed after initial creation
	m, err := pgdb_v1.Migrations(ctx, pg.DB, &Pet{}, pgdb_v1.DialectV13)
	require.NoError(t, err)
	require.Len(t, m, 0, "should have zero migrations after schema create")

	// Find the profile index (GIN index on tenant_id, profile)
	desc := (&Pet{}).DBReflect(pgdb_v1.DialectV13).Descriptor()
	tableName := desc.TableName()
	var profileIdx *pgdb_v1.Index
	for _, idx := range desc.Indexes() {
		if strings.Contains(idx.Name, "profile") {
			profileIdx = idx
			break
		}
	}
	require.NotNil(t, profileIdx, "should find the profile index")

	// Drop the existing index and recreate with different columns (simulating drift)
	_, err = pg.DB.Exec(ctx, `DROP INDEX IF EXISTS "`+profileIdx.Name+`"`)
	require.NoError(t, err)
	_, err = pg.DB.Exec(ctx, `CREATE INDEX "`+profileIdx.Name+`" ON "`+tableName+`" USING GIN ("pb$tenant_id")`)
	require.NoError(t, err)

	// Migrations should detect the drift and emit DROP + CREATE
	migrations, err := pgdb_v1.Migrations(ctx, pg.DB, &Pet{}, pgdb_v1.DialectV13)
	require.NoError(t, err)
	require.Len(t, migrations, 2, "should have 2 migrations: DROP + CREATE for drifted index")
	require.Contains(t, migrations[0], "DROP INDEX")
	require.Contains(t, migrations[0], profileIdx.Name)
	require.Contains(t, migrations[1], "CREATE INDEX")
	require.Contains(t, migrations[1], profileIdx.Name)
	require.Contains(t, migrations[1], "pb$profile")

	// Execute the migrations
	for _, line := range migrations {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "failed to execute migration: '\n%s\n'", line)
	}

	// Verify no more migrations needed
	m, err = pgdb_v1.Migrations(ctx, pg.DB, &Pet{}, pgdb_v1.DialectV13)
	require.NoError(t, err)
	require.Len(t, m, 0, "should have zero migrations after repair")
}

func TestSchemaBook(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	schema, err := pgdb_v1.CreateSchema(&Book{}, pgdb_v1.DialectV13)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		// fmt.Fprintf(os.Stderr, "---------\n%s\n\n", line)
		require.NoErrorf(t, err, "TestSchemaBook: failed to execute sql: '\n%s\n'", line)
	}

	query, params, err := pgdb_v1.Insert(Book_builder{
		TenantId: "t1",
		Id:       "b1",
		Ebook: EBook_builder{
			Size: 4000,
		}.Build(),
	}.Build(), pgdb_v1.DialectV13)
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

	schema, err := pgdb_v1.CreateSchema(&ScalarValue{}, pgdb_v1.DialectV13)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		// fmt.Fprintf(os.Stderr, "---------\n%s\n\n", line)
		require.NoErrorf(t, err, "TestSchemaScalarValue: failed to execute sql: '\n%s\n'", line)
	}

	query, params, err := pgdb_v1.Insert(&ScalarValue{}, pgdb_v1.DialectV13)
	require.NoError(t, err)

	_, err = pg.DB.Exec(ctx, query, params...)
	require.NoError(t, err)

	query, params, err = pgdb_v1.Insert(ScalarValue_builder{
		TenantId:       "t1",
		Id:             "b1",
		RepeatedString: []string{"hello", "world"},
		RepeatedEnum:   []FurType{FurType_FUR_TYPE_LOTS, FurType_FUR_TYPE_NONE},
		RepeatedFloat:  []float32{1.2, 3.2},
		RepeatedBytes: [][]byte{
			[]byte("hi"),
			[]byte("mars"),
		},
	}.Build(), pgdb_v1.DialectV13)
	require.NoError(t, err)

	_, err = pg.DB.Exec(ctx, query, params...)
	// spew.Dump(query, params)
	// spew.Dump(record)
	// fmt.Fprintf(os.Stderr, "---------\n%s\n\n", query)
	require.NoError(t, err, "invalid sql insert: %w %s", err, query)
}

// TestBookOneofNestedAccessors verifies that nested query builder accessors
// work correctly for oneof fields in the Book message.
func TestBookOneofNestedAccessors(t *testing.T) {
	bookFields := (*Book)(nil).DB().Query()
	qb := goqu.Dialect("postgres")
	tableName := (*Book)(nil).DB().TableName()

	tests := []struct {
		name        string
		expr        exp.Expression
		mustContain string
	}{
		{
			name:        "paper.pages",
			expr:        bookFields.Paper().UnsafePages().Gt(int32(100)),
			mustContain: `"pb$50$pages" > 100`,
		},
		{
			name:        "ebook.size",
			expr:        bookFields.Ebook().UnsafeSize().Lt(int64(1024 * 1024)),
			mustContain: `"pb$51$size" <`,
		},
		{
			name:        "news.id",
			expr:        bookFields.News().UnsafeId().Eq("news123"),
			mustContain: `"pb$52$id" = 'news123'`,
		},
		{
			name:        "news.created_at",
			expr:        bookFields.News().UnsafeCreatedAt().IsNotNull(),
			mustContain: `"pb$52$created_at" IS NOT NULL`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, _, err := qb.Select(goqu.L("1")).From(tableName).Where(tt.expr).ToSQL()
			require.NoError(t, err)
			require.Contains(t, sql, tt.mustContain)
		})
	}
}
