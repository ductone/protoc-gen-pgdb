package pgdb

import (
	"io"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type descriptorTemplateContext struct {
	Type             string
	ReceiverType     string
	TableName        string
	Fields           []*fieldContext
	NestedFieldNames []string
	Indexes          []*indexContext
}

func (module *Module) renderDescriptor(ctx pgsgo.Context, w io.Writer, in pgs.File, m pgs.Message, ix *importTracker) error {
	ix.PGDBV1 = true
	mt := getDescriptorType(ctx, m)
	tableName, err := getTableName(m)
	if err != nil {
		return err
	}

	fields := module.getMessageFields(ctx, m, ix, "m.self")
	nestedFieldNames := getNesteFieldNames(fields)

	c := &descriptorTemplateContext{
		Type:             mt,
		ReceiverType:     mt,
		Fields:           fields,
		NestedFieldNames: nestedFieldNames,
		Indexes:          module.getMessageIndexes(ctx, m, ix),
		TableName:        tableName,
	}
	return templates["descriptor.tmpl"].Execute(w, c)
}

func getDescriptorType(ctx pgsgo.Context, m pgs.Message) string {
	return "pgdbDescriptor" + ctx.Name(m).String()
}
