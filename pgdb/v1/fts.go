package v1

import (
	"fmt"
	"strings"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/filters/ascii"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/ductone/protoc-gen-pgdb/internal/stackoverflow"
)

type SearchContent struct {
	Type   FieldOptions_FullTextType
	Weight FieldOptions_FullTextWeight
	Value  string
}

// FullTextSearchVectors converts a set of input documents
// into a ::tsvector.
func FullTextSearchVectors(docs []SearchContent, additionalFilters ...jargon.Filter) (string, error) {
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
			rv = append(rv, pgLexeme(v, pos, doc.Weight))

			grams, err := jargon.TokenizeString(v).Filter(edgeGramFilter).Words().ToSlice()
			if err != nil {
				return "", err
			}
			for _, gram := range grams {
				rv = append(rv, pgLexeme(gram.String(), pos, doc.Weight))
			}
			pos += len(v)
		}
		if err := ts.Err(); err != nil {
			return "", err
		}

	}
	//
	// TODO: add a "transcode" version for FTS data field
	// __transcode_version:1
	//
	// https://github.com/clipperhouse/jargon
	// desc := msg.DBReflect().Descriptor()
	// option 1:
	// do ngrams
	// do split
	// do stemming AND non-stemmed
	// "aaa:3 abb:3 "::tsvector
	// + READ side needs function
	//    websearch_to_tsquery
	//
	// option 2:
	//   ... do ngrams?
	//  use webserach()
	// to_tsvector?
	// (more or less what we do today)
	return strings.Join(rv, " "), nil
}

// expose function to do same stemming
// returns exp including websearch_to_tsquery()

func FullTextSerachQuery(input string, filters ...jargon.Filter) exp.Expression {

	return nil
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
