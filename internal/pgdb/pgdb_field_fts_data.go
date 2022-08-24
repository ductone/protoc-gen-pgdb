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
	Field   pgs.Field
	Ext     *pgdb_v1.FieldOptions
}

func (fdc *ftsDataConvert) CodeForValue() (string, error) {
	fdc.SearchFields = []*searchFieldContext{}
	for _, field := range fdc.m.Fields() {
		ext := &pgdb_v1.FieldOptions{}
		_, err := field.Extension(pgdb_v1.E_Options, ext)
		if err != nil {
			panic(fmt.Errorf("pgdb: getField: failed to extract Message extension from '%s': %w", field.FullyQualifiedName(), err))
		}
		if ext.FullTextType != pgdb_v1.FieldOptions_FULL_TEXT_TYPE_UNSPECIFIED {
			fdc.SearchFields = append(fdc.SearchFields, &searchFieldContext{
				Ext:     ext,
				Field:   field,
				VarName: "m.self." + fdc.ctx.Name(field).String(),
			})
		}
	}

	if len(fdc.SearchFields) == 0 {
		return fdc.VarName + ` := exp.NewLiteralExpression("NULL")`, nil
	}

	return templateExecToString("field_fts_data.tmpl", fdc)
}

func (fdc *ftsDataConvert) VarForValue() (string, error) {
	return fdc.VarName, nil
}
