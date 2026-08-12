package v1

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ductone/protoc-gen-pgdb/internal/pgtest"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
)

// A NUL byte in a key field must not reach the generated pk/sk values.
// Postgres text can never store NUL, and under pgx's simple protocol an
// interpolated NUL truncates the wire message's query string, so the server
// rejects the whole statement with 08P01 "invalid message format". Plain
// string columns already strip NUL; pk/sk/PKSK must match.
func TestRecord_NulByteInKeyField(t *testing.T) {
	pet := Pet_builder{
		TenantId:    "t1",
		Id:          "pet\x00nul",
		CreatedAt:   timestamppb.New(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		UpdatedAt:   timestamppb.New(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		DisplayName: "Dog",
		Description: "A friendly dog",
		Elapsed:     durationpb.New(time.Hour),
		Profile:     &structpb.Struct{},
	}.Build()

	for _, dialect := range []pgdb_v1.Dialect{pgdb_v1.DialectV13, pgdb_v1.DialectV17} {
		dbr := pet.DBReflect(dialect)
		record, err := dbr.Record()
		require.NoError(t, err)

		for _, col := range []string{"pb$pk", "pb$sk"} {
			v, ok := record[col].(string)
			require.Truef(t, ok, "%s should be a plain string", col)
			require.NotContainsf(t, v, "\x00", "%s must not contain NUL", col)
		}

		pker, ok := dbr.(pgdb_v1.PrimaryKeyer)
		require.True(t, ok)
		require.NotContains(t, pker.PKSK(), "\x00")
	}
}

func TestInsertMany_NulByteInKeyField(t *testing.T) {
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
		require.NoErrorf(t, err, "failed to execute sql: '\n%s\n'", line)
	}

	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	pets := []pgdb_v1.DBReflectMessage{
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet\x00nul",
			CreatedAt:   timestamppb.New(baseTime),
			UpdatedAt:   timestamppb.New(baseTime),
			DisplayName: "Dog",
			Description: "A friendly dog",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
		}.Build(),
		Pet_builder{
			TenantId:    "t1",
			Id:          "pet2",
			CreatedAt:   timestamppb.New(baseTime),
			UpdatedAt:   timestamppb.New(baseTime),
			DisplayName: "Cat",
			Description: "A curious cat",
			Elapsed:     durationpb.New(time.Hour),
			Profile:     &structpb.Struct{},
		}.Build(),
	}

	sql, args, err := pgdb_v1.InsertMany(pgdb_v1.DialectV13, pets...)
	require.NoError(t, err)
	_, err = pg.DB.Exec(ctx, sql, args...)
	require.NoError(t, err)

	var pk string
	err = pg.DB.QueryRow(ctx,
		"SELECT pb$pk FROM "+pets[0].DBReflect(pgdb_v1.DialectV13).Descriptor().TableName()+" WHERE pb$pk LIKE '%petnul%'",
	).Scan(&pk)
	require.NoError(t, err)
	require.False(t, strings.Contains(pk, "\x00"))
}
