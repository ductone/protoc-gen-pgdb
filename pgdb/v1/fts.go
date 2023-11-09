package v1

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/filters/ascii"
	"github.com/clipperhouse/jargon/filters/stemmer"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/ductone/protoc-gen-pgdb/internal/stackoverflow"
)

type SearchContent struct {
	Type   FieldOptions_FullTextType
	Weight FieldOptions_FullTextWeight
	Value  interface{}
}

func interfaceToValue(in interface{}) string {
	if in == nil {
		return ""
	}

	switch v := in.(type) {
	case bool:
		return strconv.FormatBool(v)
	case int32:
		if v == 0 {
			return ""
		}
		return strconv.FormatInt(int64(v), 10)
	case int64:
		if v == 0 {
			return ""
		}
		return strconv.FormatInt(v, 10)
	case uint32:
		if v == 0 {
			return ""
		}
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		if v == 0 {
			return ""
		}
		return strconv.FormatUint(v, 10)
	case string:
		return v
	case []string:
		return strings.Join(v, " ")
	default:
		return ""
	}
}

func lemmatizeDocs(docs []*SearchContent, additionalFilters ...jargon.Filter) []lexeme {
	edgeGramFilter := edgegramStream(3)
	filters := []jargon.Filter{lowerCaseFilter, ascii.Fold, stackoverflow.Tags}
	filters = append(filters, additionalFilters...)
	rv := make([]lexeme, 0, 8)
	pos := 1
	for _, doc := range docs {
		docValue := interfaceToValue(doc.Value)
		ts := jargon.TokenizeString(docValue).Filter(filters...)
		for ts.Scan() {
			// wrap around guard for position
			if pos > 32767 {
				pos = 1
			}

			token := ts.Token()
			v := token.String()
			if token.IsPunct() || token.IsSpace() {
				pos += len(v)
				continue
			}

			// always add the full word, without stemming (besides stackoverflow normalization)
			rv = append(rv, lexeme{v, pos, doc.Weight})
			switch doc.Type {
			case FieldOptions_FULL_TEXT_TYPE_EXACT:
				// no additional indexing for exact
			case FieldOptions_FULL_TEXT_TYPE_ENGLISH, FieldOptions_FULL_TEXT_TYPE_UNSPECIFIED, FieldOptions_FULL_TEXT_TYPE_ENGLISH_LONG:
				// for english, we also index "edge-grams" to make partial word matching work better.
				if doc.Type != FieldOptions_FULL_TEXT_TYPE_ENGLISH_LONG {
					grams, _ := jargon.TokenizeString(v).Filter(edgeGramFilter).Words().ToSlice()
					gramWeight := lowerWeight(doc.Weight)
					for _, gram := range grams {
						rv = append(rv, lexeme{gram.String(), pos, gramWeight})
					}
				}

				// we also apply stemming. yay?
				stemmed, _ := jargon.TokenizeString(v).Filter(stemmer.English).Words().ToSlice()
				for _, stemmedWord := range stemmed {
					rv = append(rv, lexeme{stemmedWord.String(), pos, doc.Weight})
				}
			}
			pos += len(v)
		}
		if err := ts.Err(); err != nil {
			// we eat the error on purpose
			_ = err
		}
	}
	return rv
}

func camelSplitDocs(docs []*SearchContent) []lexeme {
	rv := make([]lexeme, 0, 8)
	for _, doc := range docs {
		if doc.Type == FieldOptions_FULL_TEXT_TYPE_ENGLISH_LONG {
			continue
		}
		docValue := interfaceToValue(doc.Value)
		var word []rune
		for i, r := range docValue {
			if unicode.IsUpper(r) {
				if len(word) > 0 {
					rv = append(rv, lexeme{strings.ToLower(string(word)), i - len(word) + 1, doc.Weight})
				}
				word = []rune{r}
			} else if len(word) > 0 {
				word = append(word, r)
			}
		}
		if len(word) > 0 {
			rv = append(rv, lexeme{strings.ToLower(string(word)), len(docValue) - len(word) + 1, doc.Weight})
			word = nil
		}
		word = nil
		var prev rune
		for i, r := range docValue {
			if prev == 0 {
				prev = r
				continue
			}
			if unicode.IsUpper(prev) {
				if unicode.IsSpace(r) && len(word) > 0 {
					word = append(word, prev)
					rv = append(rv, lexeme{strings.ToLower(string(word)), i - len(word), doc.Weight})
					word = nil
				} else if !unicode.IsUpper(r) && len(word) > 0 {
					rv = append(rv, lexeme{strings.ToLower(string(word)), i - len(word), doc.Weight})
					word = nil
				} else {
					word = append(word, prev)
				}
			}
			prev = r
		}
		if len(word) > 0 {
			if unicode.IsUpper(prev) {
				word = append(word, prev)
			}
			rv = append(rv, lexeme{strings.ToLower(string(word)), len(docValue) - len(word), doc.Weight})
		}
	}
	return rv
}

// FullTextSearchVectors converts a set of input documents
// into a ::tsvector. Note: this function may generally ignore errors in input text, to be robust to
// untrusted inputs, and will do its "best", for some value of "best".
func FullTextSearchVectors(docs []*SearchContent, additionalFilters ...jargon.Filter) exp.Expression {
	rv := make([]lexeme, 0, 8)

	rv = append(rv, lemmatizeDocs(docs, additionalFilters...)...)
	rv = append(rv, camelSplitDocs(docs)...)

	if len(rv) == 0 {
		return exp.NewLiteralExpression("''::tsvector")
	}

	sb := strings.Builder{}
	for _, v := range rv {
		_, _ = sb.WriteString(pgLexeme(v.value, v.pos, v.weight))
		_, _ = sb.WriteString(" ")
	}

	return exp.NewLiteralExpression("?::tsvector", sb.String())
}

func FullTextSearchQuery(input string, additionalFilters ...jargon.Filter) exp.Expression {
	filters := []jargon.Filter{lowerCaseFilter, ascii.Fold, stackoverflow.Tags}
	filters = append(filters, additionalFilters...)

	jargon.TokenizeString(input)
	terms, _ := jargon.TokenizeString(input).Filter(filters...).String()
	stemmedTerms, _ := jargon.TokenizeString(input).Filter(stemmer.English).String()

	terms = cleanToken(terms)
	stemmedTerms = cleanToken(stemmedTerms)

	return exp.NewLiteralExpression(
		"(websearch_to_tsquery('simple', ?) || websearch_to_tsquery('simple', ?))",
		terms, stemmedTerms)
}

type lexeme struct {
	value  string
	pos    int
	weight FieldOptions_FullTextWeight
}

func cleanToken(in string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) || unicode.IsLetter(r) || unicode.IsSpace(r) {
			return r
		}
		return specialReplaceChar
	}, in)
}

const specialReplaceChar = '�'

func pgLexeme(value string, pos int, weight FieldOptions_FullTextWeight) string {
	value = cleanToken(value)
	sb := strings.Builder{}
	_, _ = sb.WriteString("'")
	_, _ = sb.WriteString(value)
	_, _ = sb.WriteString("'")
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(strconv.FormatInt(int64(pos), 10))
	_, _ = sb.WriteString(weightToString(weight))
	return sb.String()
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
