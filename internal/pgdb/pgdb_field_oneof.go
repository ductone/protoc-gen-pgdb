package pgdb

import (
	"fmt"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func (module *Module) getOneOf(ctx pgsgo.Context, oneof pgs.OneOf, vn *varNamer, ix *importTracker, goPrefix string) *fieldContext {
	pgColName, err := getColumnOneOfName(oneof)
	if err != nil {
		panic(fmt.Errorf("pgdb: getColumnOneOfName failed for: %s: %w",
			oneof.FullyQualifiedName(), err))
	}

	dbTypeRef := pgDataTypeForName("int4")
	rv := &fieldContext{
		IsVirtual: false,
		GoName:    ctx.Name(oneof).String(),
		DB: &pgdb_v1.Column{
			Name:     pgColName,
			Type:     dbTypeRef.Name,
			Nullable: false,
		},
		DataType: dbTypeRef,
		Convert: &oneofDataConvert{
			ctx:     ctx,
			ix:      ix,
			oneof:   oneof,
			VarName: vn.String(),
		},
	}

	return rv
}

type oneofDataConvert struct {
	ctx     pgsgo.Context
	ix      *importTracker
	VarName string
	oneof   pgs.OneOf
}

type oneofMemberField struct {
	FieldNumber uint32
	GoType      string
	Field       pgs.Field
}

type oneofFieldContext struct {
	VarName string
	GoName  string
	Fields  []*oneofMemberField
}

func (tidc *oneofDataConvert) GoType() (string, error) {
	return pgTypeInt4, nil
}

func (fdc *oneofDataConvert) CodeForValue() (string, error) {
	c := &oneofFieldContext{
		VarName: fdc.VarName,
		GoName:  fdc.ctx.Name(fdc.oneof).String(),
	}
	for _, field := range fdc.oneof.Fields() {
		c.Fields = append(c.Fields, &oneofMemberField{
			FieldNumber: uint32(*field.Descriptor().Number),
			GoType:      fdc.ctx.OneofOption(field).String(),
			Field:       field,
		})
	}

	if len(c.Fields) == 0 {
		return fdc.VarName + ` := uint32(0)`, nil
	}

	return templateExecToString("field_oneof.tmpl", c)
}

func (fdc *oneofDataConvert) VarForValue() (string, error) {
	return fdc.VarName, nil
}

func (fdc *oneofDataConvert) VarForAppend() (string, error) {
	return "", nil
}
