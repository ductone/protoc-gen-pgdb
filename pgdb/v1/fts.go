package v1

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

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

const (
	minWordSize          = 3
	kiloByte             = 1000
	lexemeMaxBytes       = kiloByte * 2
	tsvectorMaxMegabytes = kiloByte * 1000
)

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

// snakeSubTokensSplitDoc tokenizes []string{bar, baz, quo} from foo_bar_baz quix_quo.
func snakeSubTokensSplitDoc(docValue string, wordBuffer bytes.Buffer, doc *SearchContent) []lexeme {
	wordBuffer.Reset()
	rv := make([]lexeme, 0, 8)
	var pos = 1
	var prev rune
	for _, r := range docValue {
		if prev == 0 {
			prev = r
			continue
		}
		if prev == '_' {
			if wordBuffer.Len() == 0 { // no current word
				if unicode.IsPunct(r) || unicode.IsSpace(r) || unicode.IsControl(r) || unicode.IsSymbol(r) {
					prev = r
					continue
				}
				// starting a new word
				if _, e := wordBuffer.WriteRune(r); e != nil {
					wordBuffer.Reset()
					continue
				}
			}
		} else if wordBuffer.Len() > 0 {
			// in a current word, do we append or end?
			switch {
			case r == '_' || unicode.IsPunct(r) || unicode.IsSpace(r) || unicode.IsControl(r) || unicode.IsSymbol(r):
				if utf8.RuneCount(wordBuffer.Bytes()) >= minWordSize {
					// have a word, current is a rune that doesn't continue the word so end current word
					rv = append(rv, lexeme{strings.ToLower(wordBuffer.String()), pos, doc.Weight})
				}
				wordBuffer.Reset()
			default:
				// in word and current rune continues word so continue appending
				if _, e := wordBuffer.WriteRune(r); e != nil {
					wordBuffer.Reset()
					continue
				}
			}
		}
		prev = r
		pos += 1
	}
	if utf8.RuneCount(wordBuffer.Bytes()) >= minWordSize {
		rv = append(rv, lexeme{strings.ToLower(wordBuffer.String()), pos, doc.Weight})
	}
	return rv
}

// snakeFullTokensSplitDoc tokenizes []string{foobarbaz, quixquo} from foo_bar_baz quix_quo.
func snakeFullTokensSplitDoc(docValue string, wordBuffer bytes.Buffer, doc *SearchContent) []lexeme {
	wordBuffer.Reset()
	rv := make([]lexeme, 0, 8)
	pos := 1
	hasUnderscore := false
	for _, r := range docValue {
		if wordBuffer.Len() == 0 { // no current word
			if r == '_' {
				hasUnderscore = true
			}
			if unicode.IsPunct(r) || unicode.IsSpace(r) || unicode.IsControl(r) {
				continue
			}
			// starting a new word
			if _, e := wordBuffer.WriteRune(r); e != nil {
				wordBuffer.Reset()
				continue
			}
		} else if wordBuffer.Len() > 0 {
			// in a current word, do we append or end?
			switch {
			case r != '_' && (unicode.IsPunct(r) || unicode.IsSpace(r) || unicode.IsControl(r)):
				if hasUnderscore && utf8.RuneCount(wordBuffer.Bytes()) >= minWordSize {
					// have a word, current is a rune that doesn't continue the word so end current word
					rv = append(rv, lexeme{strings.ToLower(wordBuffer.String()), pos, doc.Weight})
				}
				hasUnderscore = false
				wordBuffer.Reset()
			default:
				// in word and current rune continues word so continue appending
				if r == '_' {
					hasUnderscore = true
					continue
				}
				if _, e := wordBuffer.WriteRune(r); e != nil {
					wordBuffer.Reset()
					hasUnderscore = false
					continue
				}
			}
		}
		pos += 1
	}

	if hasUnderscore && utf8.RuneCount(wordBuffer.Bytes()) >= minWordSize {
		rv = append(rv, lexeme{strings.ToLower(wordBuffer.String()), pos, doc.Weight})
	}
	return rv
}

func camelSplitDoc(docValue string, wordBuffer bytes.Buffer, doc *SearchContent) []lexeme {
	wordBuffer.Reset()
	rv := make([]lexeme, 0, 8)
	var pos = 1
	var prev rune
	for _, r := range docValue {
		if prev == 0 {
			prev = r
			continue
		}
		if unicode.IsUpper(prev) {
			if wordBuffer.Len() == 0 { // no current word
				if unicode.IsLower(r) {
					// got a upper case in prev and current is lower, starting a new word
					if _, e := wordBuffer.WriteRune(prev); e != nil {
						wordBuffer.Reset()
						continue
					}
					if _, e := wordBuffer.WriteRune(r); e != nil {
						wordBuffer.Reset()
						continue
					}
				}
			}
		} else if wordBuffer.Len() > 0 {
			// in a current word, do we append or end?
			switch {
			case unicode.IsLower(r):
				// in word and lower so continue appending
				if _, e := wordBuffer.WriteRune(r); e != nil {
					wordBuffer.Reset()
					continue
				}
			case utf8.RuneCount(wordBuffer.Bytes()) >= minWordSize:
				// have a word, current is not lower so end current word
				rv = append(rv, lexeme{strings.ToLower(wordBuffer.String()), pos, doc.Weight})
				wordBuffer.Reset()
			default:
				wordBuffer.Reset()
			}
		}
		prev = r
		pos += 1
	}
	if utf8.RuneCount(wordBuffer.Bytes()) >= minWordSize {
		rv = append(rv, lexeme{strings.ToLower(wordBuffer.String()), pos, doc.Weight})
	}
	return rv
}

func acronymSplitDoc(docValue string, wordBuffer bytes.Buffer, doc *SearchContent) []lexeme {
	wordBuffer.Reset()
	rv := make([]lexeme, 0, 8)
	var pos = 1
	var prev rune
	for _, r := range docValue {
		if prev == 0 {
			prev = r
			continue
		}
		if unicode.IsUpper(prev) {
			switch {
			case unicode.IsLower(r):
				// only append previous if it is upper case and and current is not lower case (i.e. don't append T in AWSTest).
				if utf8.RuneCount(wordBuffer.Bytes()) >= minWordSize {
					rv = append(rv, lexeme{strings.ToLower(wordBuffer.String()), pos, doc.Weight})
				}
				wordBuffer.Reset()
			case !unicode.IsUpper(r):
				// finish acronym if there is one of min length if we encounter space.
				if _, e := wordBuffer.WriteRune(prev); e != nil {
					wordBuffer.Reset()
					continue
				}
				if utf8.RuneCount(wordBuffer.Bytes()) >= minWordSize {
					rv = append(rv, lexeme{strings.ToLower(wordBuffer.String()), pos, doc.Weight})
				}
				wordBuffer.Reset()
			default:
				if _, e := wordBuffer.WriteRune(prev); e != nil {
					wordBuffer.Reset()
					continue
				}
			}
		}
		prev = r
		pos += 1
	}
	// finish acronym if there is one of min length.
	if wordBuffer.Len() > 0 {
		if unicode.IsUpper(prev) {
			if _, e := wordBuffer.WriteRune(prev); e != nil {
				return rv
			}
		}
		if utf8.RuneCount(wordBuffer.Bytes()) >= minWordSize {
			rv = append(rv, lexeme{strings.ToLower(wordBuffer.String()), pos, doc.Weight})
		}
	}
	return rv
}

// split foo-bar, foo.bar, and foo/bar into [foo bar]
func punctuationSplitDoc(docValue string, wordBuffer bytes.Buffer, doc *SearchContent) []lexeme {
	wordBuffer.Reset()
	rv := make([]lexeme, 0, 8)
	var pos = 1

	for i, r := range docValue {
		if unicode.IsPunct(r) {
			if utf8.RuneCount(wordBuffer.Bytes()) >= minWordSize {
				rv = append(rv, lexeme{strings.ToLower(wordBuffer.String()), pos, doc.Weight})
				pos = i + 2 // i is zero indexed, and we want to skip this current rune, so add 2 for the true start of the next token
			}
			wordBuffer.Reset()
		} else {
			if unicode.IsLetter(r) || unicode.IsDigit(r) { // || unicode.IsSpace(r) {
				_, e := wordBuffer.WriteRune(r)
				if e != nil {
					wordBuffer.Reset()
					continue
				}
			}
		}
	}

	// leftover since last append
	if utf8.RuneCount(wordBuffer.Bytes()) >= minWordSize {
		rv = append(rv, lexeme{strings.ToLower(wordBuffer.String()), pos, doc.Weight})
	}

	return rv
}

// normalizeVectorDocs - converts a set of input documents into a set of lexemes matching common patterns such as camel case, snake case and accronyms.
func normalizeVectorDocs(docs []*SearchContent) []lexeme {
	rv := make([]lexeme, 0, 8)
	for _, doc := range docs {
		if doc.Type == FieldOptions_FULL_TEXT_TYPE_ENGLISH_LONG {
			continue
		}
		docValue := interfaceToValue(doc.Value)
		var wordBuffer bytes.Buffer
		rv = append(rv, camelSplitDoc(docValue, wordBuffer, doc)...)
		rv = append(rv, snakeSubTokensSplitDoc(docValue, wordBuffer, doc)...)
		rv = append(rv, snakeFullTokensSplitDoc(docValue, wordBuffer, doc)...)
		rv = append(rv, acronymSplitDoc(docValue, wordBuffer, doc)...)
	}
	return rv
}

// FullTextSearchVectors converts a set of input documents
// into a ::tsvector. Note: this function may generally ignore errors in input text, to be robust to
// untrusted inputs, and will do its "best", for some value of "best".
func FullTextSearchVectors(docs []*SearchContent, additionalFilters ...jargon.Filter) exp.Expression {
	rv := make([]lexeme, 0, 8)

	rv = append(rv, lemmatizeDocs(docs, additionalFilters...)...)
	rv = append(rv, normalizeVectorDocs(docs)...)

	if len(rv) == 0 {
		return exp.NewLiteralExpression("''::tsvector")
	}

	sb := strings.Builder{}
	for _, v := range rv {
		// Tsvector must be less than 1 mb
		if sb.Len() > tsvectorMaxMegabytes {
			break
		}

		_, _ = sb.WriteString(pgLexeme(v.value, v.pos, v.weight))
		_, _ = sb.WriteString(" ")
	}

	return exp.NewLiteralExpression("?::tsvector", sb.String())
}

func FullTextSearchQuery(input string, additionalFilters ...jargon.Filter) exp.Expression {
	filters := []jargon.Filter{lowerCaseFilter, ascii.Fold, stackoverflow.Tags}
	filters = append(filters, additionalFilters...)

	terms, _ := jargon.TokenizeString(input).Filter(filters...).String()
	stemmedTerms, _ := jargon.TokenizeString(input).Filter(stemmer.English).String()

	terms = cleanToken(terms)
	stemmedTerms = cleanToken(stemmedTerms)

	// often for simple queries, once fully stemmed, we get the same values!
	if terms == stemmedTerms {
		return exp.NewLiteralExpression(
			"(websearch_to_tsquery('simple', ?))",
			terms)
	}

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
	rv := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) || unicode.IsLetter(r) || unicode.IsSpace(r) {
			return r
		}
		return specialReplaceChar
	}, strings.ReplaceAll(in, "_", ""))
	return rv
}

const specialReplaceChar = '�'

func pgLexeme(value string, pos int, weight FieldOptions_FullTextWeight) string {
	value = cleanToken(value)

	p := strconv.FormatInt(int64(pos), 10)
	w := weightToString(weight)

	// Count the bytes to be added to format the lexeme
	extraBytes := len(p) + len(w) + len("'") + len("'") + len(":")

	// Tsvector must be less than 2kb
	totalLength := len(value) + extraBytes
	if totalLength > lexemeMaxBytes {
		// Truncate the lexeme to fit in 2kb (minus the extra bytes which will be added later)
		value = value[:lexemeMaxBytes-extraBytes]
	}

	sb := strings.Builder{}
	_, _ = sb.WriteString("'")
	_, _ = sb.WriteString(value)
	_, _ = sb.WriteString("'")
	_, _ = sb.WriteString(":")
	_, _ = sb.WriteString(p)
	_, _ = sb.WriteString(w)

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
