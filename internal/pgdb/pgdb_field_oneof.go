package pgdb

import (
	"fmt"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

func (module *Module) getOneOf(ctx pgsgo.Context, oneof pgs.OneOf, vn *varNamer, ix *importTracker) *fieldContext {
	pgColName, err := getColumnOneOfName(oneof)
	if err != nil {
		panic(fmt.Errorf("pgdb: getColumnOneOfName failed for: %s: %w",
			oneof.FullyQualifiedName(), err))
	}

	dbTypeRef := pgDataTypeForName("int4")

	importPath := ix.importablePackageName(ix.input, oneof).String()
	goType := ctx.Name(oneof.Message()).String() + ctx.Name(oneof).String() + "Type"
	if importPath != "" {
		goType = importPath + "." + goType
	}

	rv := &fieldContext{
		IsVirtual: false,
		GoName:    ctx.Name(oneof).String(),
		DB: &pgdb_v1.Column{
			Name:     pgColName,
			Type:     dbTypeRef.Name,
			Nullable: false,
			Default:  "0",
		},
		DataType: dbTypeRef,
		Convert: &oneofDataConvert{
			ctx:     ctx,
			ix:      ix,
			oneof:   oneof,
			VarName: vn.String(),
			goType:  goType,
		},
		QueryTypeName: ctx.Name(oneof.Message()).String() + ctx.Name(oneof).String() + "QueryType",
	}

	return rv
}

type oneofDataConvert struct {
	ctx     pgsgo.Context
	ix      *importTracker
	VarName string
	oneof   pgs.OneOf
	goType  string
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

func (ofdc *oneofDataConvert) GoType() (string, error) {
	return ofdc.goType, nil
}

func (ofdc *oneofDataConvert) CodeForValue() (string, error) {
	c := &oneofFieldContext{
		VarName: ofdc.VarName,
		GoName:  ofdc.ctx.Name(ofdc.oneof).String(),
	}
	for _, field := range ofdc.oneof.Fields() {
		c.Fields = append(c.Fields, &oneofMemberField{
			//nolint:gosec // overflow not possible
			FieldNumber: uint32(*field.Descriptor().Number),
			GoType:      ofdc.ctx.OneofOption(field).String(),
			Field:       field,
		})
	}

	if len(c.Fields) == 0 {
		return ofdc.VarName + ` := uint32(0)`, nil
	}

	return templateExecToString("field_oneof.tmpl", c)
}

func (ofdc *oneofDataConvert) VarForValue() (string, error) {
	return ofdc.VarName, nil
}

func (ofdc *oneofDataConvert) VarForAppend() (string, error) {
	return "", nil
}

type oneofFieldEnumContext struct {
	StructName string
	// VarName    string
	GoType string
	Fields []*oneofMemberField
}

func (ofdc *oneofDataConvert) EnumForValue() (string, error) {
	oneof := ofdc.oneof
	ctx := ofdc.ctx

	c := &oneofFieldEnumContext{
		StructName: ctx.Name(oneof.Message()).String() + ctx.Name(oneof).String(),
		GoType:     ofdc.goType,
	}
	for _, field := range ofdc.oneof.Fields() {
		c.Fields = append(c.Fields, &oneofMemberField{
			//nolint:gosec // overflow not possible
			FieldNumber: uint32(*field.Descriptor().Number),
			GoType:      field.Name().UpperCamelCase().String(),
			Field:       field,
		})
	}

	return templateExecToString("field_oneof_enum.tmpl", c)
}
