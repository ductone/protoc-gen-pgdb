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
	vector := FullTextSearchVectors([]*SearchContent{
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
			// filter all tokens < 3 characters out?
			Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH_LONG,
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

	// an empty doc, containing just spaces and punctuation
	vector := FullTextSearchVectors([]*SearchContent{
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

func TestSearchCamelCase(t *testing.T) {
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	vector := FullTextSearchVectors([]*SearchContent{
		{
			Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
			Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
			Value:  "AWSProdUSEastAuthSNSTestSQS FooBarLOLFoobarSOSTestNOPE CheezeBreeze",
		},
	})

	requireQueryTrue(t, pg, vector, "awsproduseastauthsnstestsqs")
	requireQueryTrue(t, pg, vector, "prod")
	requireQueryTrue(t, pg, vector, "east")
	requireQueryTrue(t, pg, vector, "auth")
	requireQueryTrue(t, pg, vector, "test")
	requireQueryTrue(t, pg, vector, "foo")
	requireQueryTrue(t, pg, vector, "bar")
	requireQueryTrue(t, pg, vector, "foobar")
	requireQueryTrue(t, pg, vector, "test")
	requireQueryTrue(t, pg, vector, "cheeze")
	requireQueryTrue(t, pg, vector, "breeze")
	requireQueryTrue(t, pg, vector, "aws")
	requireQueryTrue(t, pg, vector, "sns")
	requireQueryTrue(t, pg, vector, "sqs")
	requireQueryTrue(t, pg, vector, "lol")
	requireQueryTrue(t, pg, vector, "sos")
	requireQueryTrue(t, pg, vector, "nope")

	requireQueryFalse(t, pg, vector, "github")
	requireQueryFalse(t, pg, vector, "us")
	requireQueryFalse(t, pg, vector, "spro")
	requireQueryFalse(t, pg, vector, "snste")
	requireQueryFalse(t, pg, vector, "zebre")
	requireQueryFalse(t, pg, vector, "lfoo")
	requireQueryFalse(t, pg, vector, "seast")
	requireQueryFalse(t, pg, vector, "easta")
	requireQueryFalse(t, pg, vector, "arsoste")
	requireQueryFalse(t, pg, vector, "zebree")
	requireQueryFalse(t, pg, vector, "testsq")
}

func TestCamelSplitDoc(t *testing.T) {
	testCases := []struct {
		searchContent   *SearchContent
		expectedLexemes []lexeme
	}{
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "AWSProdUSEastAuthSNSTestSQS FooBarLOLFoobarSOSTestNOPE CheezeBreeze",
			},
			expectedLexemes: []lexeme{
				{"prod", 7, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"east", 13, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"auth", 17, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"test", 24, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"foo", 31, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"bar", 34, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"foobar", 43, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"test", 50, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"cheeze", 61, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"breeze", 67, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "AWSProdUSEastAuthSNSTestSQS",
			},
			expectedLexemes: []lexeme{
				{"prod", 7, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"east", 13, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"auth", 17, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"test", 24, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "ProdEastAuthTest",
			},
			expectedLexemes: []lexeme{
				{"prod", 4, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"east", 8, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"auth", 12, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"test", 16, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "PrEastTe",
			},
			expectedLexemes: []lexeme{
				{"east", 6, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "PEastT",
			},
			expectedLexemes: []lexeme{
				{"east", 5, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "PPPPPPPEastTTTTTT",
			},
			expectedLexemes: []lexeme{
				{"east", 11, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "EastTe",
			},
			expectedLexemes: []lexeme{
				{"east", 4, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "PrEast",
			},
			expectedLexemes: []lexeme{
				{"east", 6, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "PrEaTe",
			},
			expectedLexemes: []lexeme{},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "Pr East TeReat",
			},
			expectedLexemes: []lexeme{
				{"east", 7, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"reat", 14, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
	}

	for _, tc := range testCases {
		lexemes := camelSplitDoc(tc.searchContent.Value.(string), tc.searchContent)
		require.Equal(t, tc.expectedLexemes, lexemes)
	}
}

func TestAcronymSplitDoc(t *testing.T) {

	testCases := []struct {
		searchContent   *SearchContent
		expectedLexemes []lexeme
	}{
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "AWSProdUSEastAuthSNSTestSQS FooBarLOLFoobarSOSTestNOPE CheezeBreeze",
			},
			expectedLexemes: []lexeme{
				{"aws", 4, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"sns", 21, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"sqs", 27, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"lol", 38, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"sos", 47, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"nope", 54, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "AWS US SNS SQS",
			},
			expectedLexemes: []lexeme{
				{"aws", 3, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"sns", 10, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"sqs", 14, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "FeAWSFo",
			},
			expectedLexemes: []lexeme{
				{"aws", 6, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "FeAWS ",
			},
			expectedLexemes: []lexeme{
				{"aws", 5, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  " AWSFo",
			},
			expectedLexemes: []lexeme{
				{"aws", 5, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
	}

	for _, tc := range testCases {
		lexemes := acronymSplitDoc(tc.searchContent.Value.(string), tc.searchContent)
		require.Equal(t, tc.expectedLexemes, lexemes)
	}

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
		_ = FullTextSearchVectors([]*SearchContent{
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
