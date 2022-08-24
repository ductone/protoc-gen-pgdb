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

func TestSearchBigQueryDoc(t *testing.T) {
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	// example compound document, a "ID", a "display name", and a "description"
	someObjectID := "2DmNjwzqyfzisCFmt0OrPvwJ3gT"
	vector := FullTextSearchVectors([]SearchContent{
		{
			Type:   FieldOptions_FULL_TEXT_TYPE_EXACT,
			Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
			Value:  someObjectID,
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

	// Test exact ID matching
	requireQueryTrue(t, pg, vector, someObjectID)

	// this is a new random ksuid, and... it shouldn't match!
	requireQueryFalse(t, pg, vector, "2DmnMTq8tK41cqz1b1KSnVkDUmr")

	// make sure basic query works!
	requireQueryTrue(t, pg, vector, "biquery")

	// make sure basic query also misses
	requireQueryFalse(t, pg, vector, "github")
}

func TestSearchEmpty(t *testing.T) {
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	// example compound document, a "ID", a "display name", and a "description"
	vector := FullTextSearchVectors([]SearchContent{
		{
			Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
			Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
			Value:  "       ,",
		},
	})

	require.Equal(t, `''::tsvector`, vector.(exp.LiteralExpression).Literal())
	requireQueryFalse(t, pg, vector, "2DmNjwzqyfzisCFmt0OrPvwJ3gT")
	requireQueryFalse(t, pg, vector, " ")
}

func requireQueryIs(t *testing.T, pg *pgtest.PG, vectors exp.Expression, input string, matched bool) {
	qb := goqu.Dialect("postgres")
	ctx := context.Background()
	q := FullTextSearchQuery(input)
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

func FuzzFullTextSearchQuery(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345", "☃️ snowman!"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		_ = FullTextSearchQuery(orig)
	})
}

func FuzzFullTextSearchVectors(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345", "☃️ snowman!"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		_ = FullTextSearchVectors([]SearchContent{
			{
				Type:   FieldOptions_FULL_TEXT_TYPE_EXACT,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  orig,
			},
			{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  orig,
			},
		})
	})
}
