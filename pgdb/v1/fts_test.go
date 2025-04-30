package v1

import (
	"bytes"
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/ductone/protoc-gen-pgdb/internal/pgtest"
	"github.com/stretchr/testify/assert"
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

func TestSearchSymbols(t *testing.T) {
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	dotVector := FullTextSearchVectors([]*SearchContent{
		{
			Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
			Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
			Value:  "Example.THING-PROD.Group29 aws.push_group",
		},
	})

	pathVector := FullTextSearchVectors([]*SearchContent{
		{
			Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
			Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
			Value:  "/some/file/path/foo33",
		},
	})

	underscoreVector := FullTextSearchVectors([]*SearchContent{
		{
			Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
			Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
			Value:  "foo_bar_baz_21_quux another_one_is_here_to",
		},
	})

	dashVector := FullTextSearchVectors([]*SearchContent{
		{
			Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
			Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
			Value:  "long-string-separated-by-dashes",
		},
	})

	comboVector := FullTextSearchVectors([]*SearchContent{
		{
			Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
			Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
			Value:  "abcdefg.hijklmnop qrstuv/wxyz superfoo$superbar@combo",
		},
	})

	queries := []struct {
		input   string
		matched bool
		vectors exp.Expression
	}{
		{
			"THING", true, dotVector,
		},
		{
			"THING-EXPECTEDTOFAIL", false, dotVector,
		},
		{
			"THING-PROD", true, dotVector,
		},
		{
			"Group29", true, dotVector,
		},
		{
			"Group30", false, dotVector,
		},
		{
			"Group", true, dotVector,
		},
		{
			"example.thing", true, dotVector,
		},
		{
			"example.thing-prod", true, dotVector,
		},
		{
			"example-thing", true, dotVector,
		},
		{
			"example_thing", true, dotVector,
		},
		{
			"example/thing", true, dotVector,
		},
		{
			"aws", true, dotVector,
		},
		{
			"aws.push", true, dotVector,
		},
		{
			"aws.push_group", true, dotVector,
		},
		{
			"Example", false, pathVector,
		},
		{
			"foo33", true, pathVector,
		},
		{
			"file/path", true, pathVector,
		},
		{
			"file", true, pathVector,
		},
		{
			"path/foo", true, pathVector,
		},
		{
			"path/foo33", true, pathVector,
		},
		{
			"foo", true, underscoreVector,
		},
		{
			"foo_bar", true, underscoreVector,
		},
		{
			"foo_baz", false, underscoreVector,
		},
		{
			"baz", true, underscoreVector,
		},
		{
			"quux", true, underscoreVector,
		},
		{
			"another", true, underscoreVector,
		},
		{
			"long", true, dashVector,
		},
		{
			"string", true, dashVector,
		},
		{
			"by", true, dashVector,
		},
		{
			"qrst", true, comboVector,
		},
		{
			"qrstuv/wxyz", true, comboVector,
		},
		{
			"superfoo", true, comboVector,
		},
		{
			"superbar", true, comboVector,
		},
		{
			"superfoo$superbar", true, comboVector,
		},
		{
			"combo", true, comboVector,
		},
	}

	for _, q := range queries {
		// run each in test suite
		t.Run(q.input, func(t *testing.T) {
			requireQueryIs(t, pg, q.vectors, q.input, q.matched)
		})
	}
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

func TestStemming(t *testing.T) {
	testVectors := map[string][]string{
		"carri": []string{
			"carry",
			"carrying",
			"carried",
			"carries",
		},
		"danc": []string{
			"dance",
			"dancing",
			"danced",
			"dances",
		},
		"confus": []string{
			"confuse",
			"confused",
			"confusing",
			"confusingly",
			"confuser",
			"confuses",
		},
	}

	// Test stemming for tsvector
	for expectedStem, inputs := range testVectors {
		for _, term := range inputs {
			t.Run(fmt.Sprintf("tsvector %s", term), func(t *testing.T) {
				vector := FullTextSearchVectors([]*SearchContent{
					{
						Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
						Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
						Value:  term,
					},
				}).(exp.LiteralExpression)

				for _, term := range vector.Args() {
					assert.Contains(t, term, fmt.Sprintf("'%s':1", expectedStem))
				}
			})
		}
	}

	// Test stemming for tsquery
	for expectedStem, inputs := range testVectors {
		for _, term := range inputs {
			t.Run(fmt.Sprintf("tsquery %s", term), func(t *testing.T) {
				exp := FullTextSearchQuery(term).(exp.LiteralExpression)
				var args []string
				for _, arg := range exp.Args() {
					args = append(args, arg.(string))
				}

				assert.True(t, slices.Contains(args, expectedStem))
			})
		}
	}
}

func TestFullSentence(t *testing.T) {
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	vector := FullTextSearchVectors([]*SearchContent{
		{
			Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
			Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
			Value:  "This is a full sentence with fewer than 100 billion words, I shouldn't have /any/ _problems_ matching any part of it, right?! Yeah!",
		},
	})

	requireQueryTrue(t, pg, vector, "sentence")
	requireQueryTrue(t, pg, vector, "sentences")
	requireQueryTrue(t, pg, vector, "full sentence")
	requireQueryTrue(t, pg, vector, "few")
	requireQueryTrue(t, pg, vector, "fewer")
	requireQueryTrue(t, pg, vector, "100")
	requireQueryTrue(t, pg, vector, "should")
	// requireQueryTrue(t, pg, vector, "shouldn't") // it would be nice if this matched
	requireQueryTrue(t, pg, vector, "any")
	requireQueryTrue(t, pg, vector, "problem")
	requireQueryTrue(t, pg, vector, "match")
	requireQueryTrue(t, pg, vector, "right")
	requireQueryTrue(t, pg, vector, "YEAH")

	requireQueryFalse(t, pg, vector, "false")
	requireQueryFalse(t, pg, vector, "problemz")
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

func TestSearchSnakeCase(t *testing.T) {
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	vector := FullTextSearchVectors([]*SearchContent{
		{
			Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
			Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
			Value:  "AWS_Prod_US_east_auth_sns_test_sqs foo_bar_LOL_foobar_SOS_Test_NOPE_ _Cheeze_breeze_",
		},
	})

	requireQueryTrue(t, pg, vector, "aws_prod_us_east_auth_sns_test_sqs")
	requireQueryTrue(t, pg, vector, "aws_prod_us_east")
	requireQueryTrue(t, pg, vector, "foo_bar_lol_foobar_sos_test_nope_")
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

func TestSearchPathsFull(t *testing.T) {
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	vector := FullTextSearchVectors([]*SearchContent{
		{
			Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
			Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
			Value:  "AWS/Prod/US/east/auth/sns/test/sqs foo/bar/LOL/foobar/SOS/Test/NOPE/ /Cheeze/breeze/",
		},
	})

	requireQueryTrue(t, pg, vector, "aws/prod/us/east/auth/sns/test/sqs")
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

func TestSubSnakeCaseDoc(t *testing.T) {
	testCases := []struct {
		searchContent   *SearchContent
		expectedLexemes []lexeme
	}{
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "AWS_Prod_US_east_auth_sns_test_sqs foo_bar_LOL_foobar_SOS_Test_NOPE_ _Cheeze_breeze_",
			},
			expectedLexemes: []lexeme{
				{"prod", 8, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"east", 16, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"auth", 21, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"sns", 25, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"test", 30, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"sqs", 34, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"bar", 42, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"lol", 46, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"foobar", 53, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"sos", 57, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"test", 62, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"nope", 67, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"cheeze", 75, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"breeze", 82, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "AWS_Prod_US_east_auth_sns_test_sqs-foo_bar_LOL_foobar_SOS_Test_NOPE_:_Cheeze_breeze_",
			},
			expectedLexemes: []lexeme{
				{"prod", 8, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"east", 16, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"auth", 21, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"sns", 25, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"test", 30, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"sqs", 34, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"bar", 42, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"lol", 46, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"foobar", 53, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"sos", 57, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"test", 62, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"nope", 67, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"cheeze", 75, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"breeze", 82, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "AWS_Prod_US_east_auth_sns_test_sqs",
			},
			expectedLexemes: []lexeme{
				{"prod", 8, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"east", 16, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"auth", 21, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"sns", 25, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"test", 30, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"sqs", 34, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "prod_east_auth_test",
			},
			expectedLexemes: []lexeme{
				{"east", 9, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"auth", 14, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"test", 19, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "pr_east_te",
			},
			expectedLexemes: []lexeme{
				{"east", 7, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "p_east_t",
			},
			expectedLexemes: []lexeme{
				{"east", 6, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "_____east____",
			},
			expectedLexemes: []lexeme{
				{"east", 5, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "east_te",
			},
			expectedLexemes: []lexeme{},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "pr_east",
			},
			expectedLexemes: []lexeme{
				{"east", 7, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "pr_et_te",
			},
			expectedLexemes: []lexeme{},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "pr _east_ te_reat",
			},
			expectedLexemes: []lexeme{
				{"east", 8, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"reat", 16, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "Pr$_east_+te_reat",
			},
			expectedLexemes: []lexeme{
				{"east", 8, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"reat", 16, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "_hello, 世𝒜界🌮🌍B🌎_foo_🌏",
			},
			expectedLexemes: []lexeme{
				{"hello", 6, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"foo", 19, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
	}

	symbols := map[rune]struct{}{
		'_': {},
	}

	for _, tc := range testCases {
		var wordBuffer bytes.Buffer
		lexemes := symbolsSubTokensSplitDoc(symbols, tc.searchContent.Value.(string), wordBuffer, tc.searchContent)
		require.Equal(t, tc.expectedLexemes, lexemes)
	}
}

func TestFullSnakeCaseDoc(t *testing.T) {
	testCases := []struct {
		searchContent   *SearchContent
		expectedLexemes []lexeme
	}{
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "AWS_Prod_US_east_auth_sns_test_sqs foo_bar_LOL_foobar_SOS_Test_NOPE_ _Cheeze_breeze_",
			},
			expectedLexemes: []lexeme{
				{"awsproduseastauthsnstestsqs", 28, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"foobarlolfoobarsostestnope", 55, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"cheezebreeze", 68, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "AWS_Prod_US_east_auth_sns_test_sqs-foo_bar_LOL_foobar_SOS_Test_NOPE_:_Cheeze_breeze_",
			},
			expectedLexemes: []lexeme{
				{"awsproduseastauthsnstestsqs", 28, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"foobarlolfoobarsostestnope", 55, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"cheezebreeze", 68, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "AWS_Prod_US_east_auth_sns_test_sqs",
			},
			expectedLexemes: []lexeme{
				{"awsproduseastauthsnstestsqs", 28, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "prod_east_auth_test",
			},
			expectedLexemes: []lexeme{
				{"prodeastauthtest", 17, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "pr_east_te",
			},
			expectedLexemes: []lexeme{
				{"preastte", 9, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "p_east_t",
			},
			expectedLexemes: []lexeme{
				{"peastt", 7, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "_____east____",
			},
			expectedLexemes: []lexeme{
				{"east", 5, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "east_te",
			},
			expectedLexemes: []lexeme{
				{"eastte", 7, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "pr_east",
			},
			expectedLexemes: []lexeme{
				{"preast", 7, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "pr_et_te",
			},
			expectedLexemes: []lexeme{
				{"prette", 7, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "pr _east_ te_reat",
			},
			expectedLexemes: []lexeme{
				{"east", 8, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"tereat", 15, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "Pr$_east_+te_reat",
			},
			expectedLexemes: []lexeme{
				{"pr$east+tereat", 15, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "_hello, 世𝒜界🌮🌍B🌎_foo_🌏",
			},
			expectedLexemes: []lexeme{
				{"hello", 6, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"世𝒜界🌮🌍b🌎foo🌏", 18, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
	}

	symbols := map[rune]struct{}{
		'_': {},
	}

	for _, tc := range testCases {
		var wordBuffer bytes.Buffer
		lexemes := symbolsFullTokensSplitDoc(symbols, tc.searchContent.Value.(string), wordBuffer, tc.searchContent)
		require.Equal(t, tc.expectedLexemes, lexemes)
	}
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
				Value:  "AWSProdUSEastAuthSNSTestSQS-FooBarLOLFoobarSOSTestNOPE:CheezeBreeze",
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
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "Pr$East+TeReat",
			},
			expectedLexemes: []lexeme{
				{"east", 7, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"reat", 14, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "Hello, 世𝒜界🌮🌍B🌎Foo🌏",
			},
			expectedLexemes: []lexeme{
				{"hello", 5, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"foo", 17, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
	}

	for _, tc := range testCases {
		var wordBuffer bytes.Buffer
		lexemes := camelSplitDoc(tc.searchContent.Value.(string), wordBuffer, tc.searchContent)
		require.Equal(t, tc.expectedLexemes, lexemes)
	}
}

func TestDotsSplitDoc(t *testing.T) {
	testCases := []struct {
		searchContent   *SearchContent
		expectedLexemes []lexeme
	}{
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "foo.bar.baz.quux.eek35!",
			},
			expectedLexemes: []lexeme{
				{"bar", 7, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"baz", 11, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"quux", 16, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
				{"eek35", 22, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
	}

	symbols := map[rune]struct{}{
		'.': {},
	}

	for _, tc := range testCases {
		var wordBuffer bytes.Buffer
		lexemes := symbolsSubTokensSplitDoc(symbols, tc.searchContent.Value.(string), wordBuffer, tc.searchContent)
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
				Value:  "AWS:US/SNS|SQS",
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
				Value:  "AWS-US_SNS/SQS",
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
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "FoAW",
			},
			expectedLexemes: []lexeme{},
		},
		{
			searchContent: &SearchContent{
				Type:   FieldOptions_FULL_TEXT_TYPE_ENGLISH,
				Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
				Value:  "Hello, 世𝒜界🌮🌍B🌎FOO🌏",
			},
			expectedLexemes: []lexeme{
				{"foo", 17, FieldOptions_FULL_TEXT_WEIGHT_HIGH},
			},
		},
	}

	for _, tc := range testCases {
		var wordBuffer bytes.Buffer
		lexemes := acronymSplitDoc(tc.searchContent.Value.(string), wordBuffer, tc.searchContent)
		require.Equal(t, tc.expectedLexemes, lexemes)
	}
}

func requireQueryIs(t *testing.T, pg *pgtest.PG, vectors exp.Expression, input string, matched bool) {
	qb := goqu.Dialect("postgres")
	ctx := context.Background()
	q := FullTextSearchQuery(input)
	_ = q
	query, args, err := qb.Select(
		goqu.L("? @@ ?", vectors, q),
	).ToSQL()
	require.NoError(t, err)

	var result bool

	err = pg.DB.QueryRow(ctx, query, args...).Scan(&result)
	require.NoError(t, err, "Failed to execute query: %s\nQUERY: %s\n", err, query)
	assert.Equal(t, matched, result, "Expected query matching failed: '%s'\nQUERY: %s\n\n", input, query)
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
