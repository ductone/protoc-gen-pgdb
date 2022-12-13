package pgdb

import (
	"io"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type descriptorTemplateContext struct {
	Type         string
	ReceiverType string
	TableName    string
	Fields       []*fieldContext
	NestedFields []*nestedFieldContext
	Indexes      []*indexContext
}

func (module *Module) renderDescriptor(ctx pgsgo.Context, w io.Writer, in pgs.File, m pgs.Message, ix *importTracker) error {
	ix.PGDBV1 = true
	mt := getDescriptorType(ctx, m)
	tableName, err := getTableName(m)
	if err != nil {
		return err
	}

	fields := module.getMessageFields(ctx, m, ix, "m.self")
	mestedFields := getNesteFields(ctx, fields)

	c := &descriptorTemplateContext{
		Type:         mt,
		ReceiverType: mt,
		Fields:       fields,
		NestedFields: mestedFields,
		Indexes:      module.getMessageIndexes(ctx, m, ix),
		TableName:    tableName,
	}
	return templates["descriptor.tmpl"].Execute(w, c)
}

func getDescriptorType(ctx pgsgo.Context, m pgs.Message) string {
	return "pgdbDescriptor" + ctx.Name(m).String()
}

type nestedFieldContext struct {
	GoName   string
	TypeName string
}

func getNesteFieldNames(fields []*fieldContext) []string {
	rv := make([]string, 0)
	for _, f := range fields {
		if !f.Nested {
			continue
		}
		rv = append(rv, f.GoName)
	}
	return rv
}

func getNesteFields(ctx pgsgo.Context, fields []*fieldContext) []*nestedFieldContext {
	rv := make([]*nestedFieldContext, 0)
	for _, f := range fields {
		if !f.Nested {
			continue
		}

		rv = append(rv, &nestedFieldContext{
			GoName:   f.GoName,
			TypeName: ctx.Type(f.Field).String(),
		})
		if false == false {
			continue
		}
	}
	return rv
}
