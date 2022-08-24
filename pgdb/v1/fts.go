package v1

import (
	"fmt"
	"strings"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/filters/ascii"
	"github.com/clipperhouse/jargon/filters/stemmer"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/ductone/protoc-gen-pgdb/internal/stackoverflow"
)

type SearchContent struct {
	Type   FieldOptions_FullTextType
	Weight FieldOptions_FullTextWeight
	Value  string
}

// FullTextSearchVectors converts a set of input documents
// into a ::tsvector. Note: this function may generally ignore errors in input text, to be robust to
// untrusted inputs, and will do its "best", for some value of "best"
func FullTextSearchVectors(docs []SearchContent, additionalFilters ...jargon.Filter) exp.Expression {
	edgeGramFilter := edgegramStream(3)
	filters := []jargon.Filter{lowerCaseFilter, ascii.Fold, stackoverflow.Tags}
	filters = append(filters, additionalFilters...)
	rv := make([]string, 0, 8)
	pos := 1
	for _, doc := range docs {
		ts := jargon.TokenizeString(doc.Value).Filter(filters...)
		for ts.Scan() {
			token := ts.Token()
			v := token.String()
			if token.IsPunct() || token.IsSpace() {
				pos += len(v)
				continue
			}

			// always add the full word, without stemming (besides stackoverflow normalization)
			rv = append(rv, pgLexeme(v, pos, doc.Weight))
			switch doc.Type {
			case FieldOptions_FULL_TEXT_TYPE_EXACT:
				// no additional indexing for exact
			case FieldOptions_FULL_TEXT_TYPE_ENGLISH, FieldOptions_FULL_TEXT_TYPE_UNSPECIFIED:
				// for english, we also index "edge-grams" to make partial word matching work better.
				grams, _ := jargon.TokenizeString(v).Filter(edgeGramFilter).Words().ToSlice()
				gramWeight := lowerWeight(doc.Weight)
				for _, gram := range grams {
					rv = append(rv, pgLexeme(gram.String(), pos, gramWeight))
				}

				// we also apply stemming. yay?
				stemmed, _ := jargon.TokenizeString(v).Filter(stemmer.English).Words().ToSlice()
				for _, stemmedWord := range stemmed {
					rv = append(rv, pgLexeme(stemmedWord.String(), pos, gramWeight))
				}
			}
			pos += len(v)
			// wrap around guard for position
			if pos > 2^15 {
				pos = 1
			}
		}
		if err := ts.Err(); err != nil {
			// we eat the error on purpose
			_ = err
		}
	}
	return exp.NewLiteralExpression("?::tsvector", strings.Join(rv, " "))
}

func FullTextSerachQuery(input string, additionalFilters ...jargon.Filter) exp.Expression {
	filters := []jargon.Filter{lowerCaseFilter, ascii.Fold, stackoverflow.Tags}
	filters = append(filters, additionalFilters...)

	terms, _ := jargon.TokenizeString(input).Filter(filters...).String()
	stemmedTerms, _ := jargon.TokenizeString(input).Filter(stemmer.English).String()
	return exp.NewLiteralExpression(
		"(websearch_to_tsquery('simple', ?) || websearch_to_tsquery('simple', ?))",
		terms, stemmedTerms)
}

func pgLexeme(value string, pos int, weight FieldOptions_FullTextWeight) string {
	return fmt.Sprintf("%s:%d%s", value, pos, weightToString(weight))
}

func weightToString(weight FieldOptions_FullTextWeight) string {
	switch weight {
	case FieldOptions_FULL_TEXT_WEIGHT_HIGH:
		return "A"
	case FieldOptions_FULL_TEXT_WEIGHT_MED:
		return "B"
	case FieldOptions_FULL_TEXT_WEIGHT_LOW:
		return "C"
	default:
		return "D"
	}
}

func lowerWeight(weight FieldOptions_FullTextWeight) FieldOptions_FullTextWeight {
	switch weight {
	case FieldOptions_FULL_TEXT_WEIGHT_HIGH:
		return FieldOptions_FULL_TEXT_WEIGHT_MED
	case FieldOptions_FULL_TEXT_WEIGHT_MED:
		return FieldOptions_FULL_TEXT_WEIGHT_LOW
	case FieldOptions_FULL_TEXT_WEIGHT_LOW:
		return FieldOptions_FULL_TEXT_WEIGHT_UNSPECIFIED
	default:
		return FieldOptions_FULL_TEXT_WEIGHT_UNSPECIFIED
	}
}
