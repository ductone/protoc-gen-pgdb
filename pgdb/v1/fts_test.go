package v1

import (
	"context"
	"testing"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/ductone/protoc-gen-pgdb/internal/pgtest"
	"github.com/stretchr/testify/require"
)

func TestSearchGithub(t *testing.T) {
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	// example compound document, a "ID", a "display name", and a "description"
	vector, err := FullTextSearchVectors([]SearchContent{
		{
			Type:   FieldOptions_FULL_TEXT_TYPE_EXACT,
			Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
			Value:  "2DmNjwzqyfzisCFmt0OrPvwJ3gT",
		},
		{
			Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
			Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
			Value:  "role/biquery.user",
		},
		{
			Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
			Weight: FieldOptions_FULL_TEXT_WEIGHT_MED,
			Value:  "When applied to a dataset, this role provides the ability to read the dataset's metadata and list tables in the dataset.",
		},
	})
	require.NoError(t, err)

	requireQueryTrue(t, pg, vector, "2DmNjwzqyfzisCFmt0OrPvwJ3gT")
	// this is a new random ksuid, and... it shouldn't match!
	requireQueryFalse(t, pg, vector, "2DmnMTq8tK41cqz1b1KSnVkDUmr")
}

func requireQueryIs(t *testing.T, pg *pgtest.PG, vectors exp.Expression, input string, matched bool) {
	qb := goqu.Dialect("postgres")
	ctx := context.Background()
	q := FullTextSerachQuery(input)
	query, args, err := qb.Select(
		goqu.L("? @@ ?", vectors, q),
	).ToSQL()
	require.NoError(t, err)

	var result bool
	err = pg.DB.QueryRow(ctx, query, args...).Scan(&result)
	require.NoError(t, err)
	require.Equal(t, matched, result, "Expected query matching failed: '%s'", input)
}

func requireQueryTrue(t *testing.T, pg *pgtest.PG, vectors exp.Expression, query string) {
	requireQueryIs(t, pg, vectors, query, true)
}

func requireQueryFalse(t *testing.T, pg *pgtest.PG, vectors exp.Expression, query string) {
	requireQueryIs(t, pg, vectors, query, false)
}

// https://www.code-intelligence.com/blog/fuzzing-golang-1.19
// spew.Dump(query)
// spew.Dump(FullTextSearchVectors([]SearchContent{
// 	{
// 		Type:   FieldOptions_FULL_TEXT_TYPE_EXACT,
// 		Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
// 		Value:  "github.com/pquerna/reponame",
// 	},
// }))

// 	Cols("first_name", "last_name").
// 	Vals(
// 		goqu.Vals{"Greg", "Farley"},
// 		goqu.Vals{"Jimmy", "Stewart"},
// 		goqu.Vals{"Jeff", "Jeffers"},
// 	)
// insertSQL, args, _ := ds.ToSQL()
// fmt.Println(insertSQL, args)
//

// goqu.Insert() exp.Record{
// 	"example": vector,
// }
