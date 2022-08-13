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
}

func (module *Module) renderDescriptor(ctx pgsgo.Context, w io.Writer, in pgs.File, m pgs.Message, ix *importTracker) error {
	ix.PGDB_v1 = true
	mt := getDescriptorType(m)
	tableName, err := getTableName(m)
	if err != nil {
		return err
	}

	c := &descriptorTemplateContext{
		Type:         mt,
		ReceiverType: "*" + mt,
		Fields:       module.getMessageFields(ctx, m, ix, "m.self"),
		TableName:    tableName,
	}
	return templates["descriptor.tmpl"].Execute(w, c)
}

func getDescriptorType(m pgs.Message) string {
	return "pgdbDescriptor" + m.Name().UpperCamelCase().String()
}
