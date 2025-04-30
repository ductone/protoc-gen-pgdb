package v1

import (
	"context"
	"strings"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	v1 "github.com/ductone/protoc-gen-pgdb/example/models/animals/v1"
	"github.com/ductone/protoc-gen-pgdb/internal/pgtest"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSchemaZooShop(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	_, err = pg.DB.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS btree_gin")
	require.NoError(t, err)

	smsg := &Shop{}
	schema, err := pgdb_v1.CreateSchema(smsg)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "TestSchemaZooShop: failed to execute sql: '\n%s\n'", line)
		// os.Stderr.WriteString(line)
		// os.Stderr.WriteString("\n------\n")
	}
	ct := schema[0]
	require.Contains(t, ct, "CREATE TABLE")
	require.Equal(t, 3, // todo: fix this after postgresql migration
		strings.Count(ct, "$pksk"),
		"Create table should contain one pksk and one pkskv2 field + index: %s", ct,
	)
	require.Equal(t, 1,
		strings.Count(ct, "fts_data"),
		"Create table should contain only one fts_data field: %s", ct,
	)
	_, err = pg.DB.Exec(ctx, "DROP TABLE "+smsg.DBReflect().Descriptor().TableName())
	require.NoErrorf(t, err, "TestSchemaZooShop: failed to drop")

	schema, err = pgdb_v1.Migrations(ctx, pg.DB, smsg)
	require.NoError(t, err)
	for _, line := range schema {
		_, err := pg.DB.Exec(ctx, line)
		require.NoErrorf(t, err, "TestSchemaZooShop: failed to execute sql: '\n%s\n'", line)
		// os.Stderr.WriteString(line)
		// os.Stderr.WriteString("\n------\n")
	}
	ct = schema[0]
	require.Contains(t, ct, "CREATE TABLE")

	s := &Shop{
		TenantId:  "t1",
		Id:        "s1",
		CreatedAt: timestamppb.Now(),
		Fur:       v1.FurType_FUR_TYPE_LOTS,
		Medium: &Shop_Anything{
			Anything: &v1.ScalarValue{
				String_:        "unique",
				RepeatedString: []string{"xyz", "zyx"},
			},
		},
	}
	found := false
	searchData := s.DBReflect().SearchData()
	for _, sd := range searchData {
		if sd.Value == "unique" {
			found = true
		}
	}
	require.True(t, found, "expected string in FTS data: %v", searchData)
	vectors := pgdb_v1.FullTextSearchVectors(searchData)
	qb := goqu.Dialect("postgres")
	sql, _, err := qb.Select(exp.NewAliasExpression(vectors, "vectors")).ToSQL()
	require.NoError(t, err)
	require.Contains(t, sql, "''unique'':")
	require.Contains(t, sql, "''xyz'':")
	require.Contains(t, sql, "''zyx'':")
}
