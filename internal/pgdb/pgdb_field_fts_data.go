package pgdb

import (
	"fmt"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type ftsDataConvert struct {
	ctx          pgsgo.Context
	VarName      string
	m            pgs.Message
	SearchFields []*searchFieldContext
}

type searchFieldContext struct {
	VarName string
	Weight  string
	Field   pgs.Field
	Ext     *pgdb_v1.FieldOptions
}

func weightToString(weight pgdb_v1.FieldOptions_FullTextWeight) string {
	switch weight {
	case pgdb_v1.FieldOptions_FULL_TEXT_WEIGHT_HIGH:
		return "A"
	case pgdb_v1.FieldOptions_FULL_TEXT_WEIGHT_MED:
		return "B"
	case pgdb_v1.FieldOptions_FULL_TEXT_WEIGHT_LOW:
		return "C"
	default:
		return "D"
	}
}

func (fdc *ftsDataConvert) CodeForValue() (string, error) {
	fdc.SearchFields = []*searchFieldContext{}
	for _, field := range fdc.m.Fields() {
		ext := &pgdb_v1.FieldOptions{}
		_, err := field.Extension(pgdb_v1.E_Options, ext)
		if err != nil {
			panic(fmt.Errorf("pgdb: getField: failed to extract Message extension from '%s': %w", field.FullyQualifiedName(), err))
		}
		if ext.FullTextWeight == pgdb_v1.FieldOptions_FULL_TEXT_WEIGHT_UNSPECIFIED {
			ext.FullTextWeight = pgdb_v1.FieldOptions_FULL_TEXT_WEIGHT_MED
		}
		if ext.FullTextType != pgdb_v1.FieldOptions_FULL_TEXT_TYPE_UNSPECIFIED {
			fdc.SearchFields = append(fdc.SearchFields, &searchFieldContext{
				Ext:     ext,
				Field:   field,
				VarName: "m.self." + fdc.ctx.Name(field).String(),
				Weight:  weightToString(ext.FullTextWeight),
			})
		}
	}

	if len(fdc.SearchFields) <= 0 {
		return fdc.VarName + ` := exp.NewLiteralExpression("NULL")`, nil
	}

	return templateExecToString("field_fts_data.tmpl", fdc)
}

func (fdc *ftsDataConvert) VarForValue() (string, error) {
	return fdc.VarName, nil
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
