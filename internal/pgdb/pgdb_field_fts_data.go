package pgdb

import (
	"fmt"
	"strings"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type ftsDataConvert struct {
	ctx              pgsgo.Context
	VarName          string
	m                pgs.Message
	SearchFields     []*searchFieldContext
	NestedFieldNames []string
}

type searchFieldContext struct {
	VarName string
	Field   pgs.Field
	Ext     *pgdb_v1.FieldOptions
}

func (tidc *ftsDataConvert) GoType() (string, error) {
	return "string", nil
}

func getSearchFields(ctx pgsgo.Context, m pgs.Message) []*searchFieldContext {
	rv := []*searchFieldContext{}
	for _, field := range m.Fields() {
		ext := &pgdb_v1.FieldOptions{}
		_, err := field.Extension(pgdb_v1.E_Options, ext)
		if err != nil {
			panic(fmt.Errorf("pgdb: getField: failed to extract Message extension from '%s': %w", field.FullyQualifiedName(), err))
		}
		if ext.FullTextType != pgdb_v1.FieldOptions_FULL_TEXT_TYPE_UNSPECIFIED {
			rv = append(rv, &searchFieldContext{
				Ext:     ext,
				Field:   field,
				VarName: "m.self.Get" + ctx.Name(field).String() + "()",
			})
		}
	}
	return rv
}

func (fdc *ftsDataConvert) CodeForValue() (string, error) {
	fdc.SearchFields = []*searchFieldContext{}
	fdc.NestedFieldNames = []string{}
	for _, field := range fdc.m.Fields() {
		ext := &pgdb_v1.FieldOptions{}
		_, err := field.Extension(pgdb_v1.E_Options, ext)
		if err != nil {
			panic(fmt.Errorf("pgdb: getField: failed to extract Message extension from '%s': %w", field.FullyQualifiedName(), err))
		}

		// NOTE: this mirrors logic in pgdb_field.go to find nested fields :(
		pt := field.Type().ProtoType()
		if pt == pgs.MessageT {
			if !strings.HasPrefix(field.Descriptor().GetTypeName(), ".google.protobuf") &&
				(ext.MessageBehavoir == pgdb_v1.FieldOptions_MESSAGE_BEHAVIOR_EXPAND ||
					ext.MessageBehavoir == pgdb_v1.FieldOptions_MESSAGE_BEHAVIOR_UNSPECIFIED) {
				fdc.NestedFieldNames = append(fdc.NestedFieldNames, fdc.ctx.Name(field).String())
			}
		}

		if ext.FullTextType != pgdb_v1.FieldOptions_FULL_TEXT_TYPE_UNSPECIFIED {
			fdc.SearchFields = append(fdc.SearchFields, &searchFieldContext{
				Ext:     ext,
				Field:   field,
				VarName: "m.self.Get" + fdc.ctx.Name(field).String() + "()",
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

func (fdc *ftsDataConvert) VarForAppend() (string, error) {
	return "", nil
}

func (tidc *ftsDataConvert) EnumForValue() (string, error) {
	return "", nil
}
