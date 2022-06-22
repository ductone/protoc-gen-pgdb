package pgdb

import pgsgo "github.com/lyft/protoc-gen-star/lang/go"

type ftsDataConvert struct {
	ctx     pgsgo.Context
	VarName string
}

func (pbdc *ftsDataConvert) CodeForValue() (string, error) {
	return templateExecToString("field_pbdata.tmpl", pbdc)
}

func (pbdc *ftsDataConvert) VarForValue() (string, error) {
	return pbdc.VarName, nil
}

// func SearchFieldsToQuery(fields []*fieldContext) (exp.LiteralExpression, error) {
// 	if countSearchFields(fields) == 0 {
// 		return exp.L("NULL"), nil
// 	}
// 	edgeNGramTemplate := "edge_gram_tsvector(?::text, ?)"
// 	vectorTemplate := "setweight(to_tsvector(?, ?::text), ?)"

// 	vectors := []string{}
// 	args := []interface{}{}
// 	for _, field := range fields {
// 		if field.Convert.FullTextType == pgdb_v1.FieldOptions_FULL_TEXT_TYPE_UNSPECIFIED {
// 			continue
// 		}

// 		switch field.Convert.FullTextType {
// 		case pgdb_v1.FieldOptions_FULL_TEXT_TYPE_EDGE_NGRAM:
// 			vectors = append(vectors, edgeNGramTemplate)
// 			args = append(args, field.textContents(), field.weight())
// 			fallthrough
// 		case pgdb_v1.FieldOptions_FULL_TEXT_TYPE_ENGLISH:
// 			vectors = append(vectors, vectorTemplate)
// 			args = append(args, "english", field.textContents(), field.weight())
// 			fallthrough
// 		case pgdb_v1.FieldOptions_FULL_TEXT_TYPE_SIMPLE:
// 			vectors = append(vectors, vectorTemplate)
// 			args = append(args, "simple", field.textContents(), field.weight())
// 		}
// 	}

// 	query := strings.Join(vectors, " || ")
// 	return goqu.L(query, args...), nil
// }
