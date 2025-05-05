package pgdb

import (
	"io"
	"strconv"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type descriptorTemplateContext struct {
	Type                        string
	ReceiverType                string
	TableName                   string
	Fields                      []*fieldContext
	NestedFields                []*nestedFieldContext
	Indexes                     []*indexContext
	Statistics                  []*statsContext
	VersioningField             string
	IsPartitioned               bool
	IsPartitionedByCreatedAt    bool
	PartitionedByKsuidFieldName string
	PartitionDateRange          string
	UsePkskv2Column             bool
}

func (module *Module) renderDescriptor(ctx pgsgo.Context, w io.Writer, in pgs.File, m pgs.Message, ix *importTracker) error {
	ix.PGDBV1 = true
	fext := pgdb_v1.MessageOptions{}
	_, err := m.Extension(pgdb_v1.E_Msg, &fext)
	if err != nil {
		panic(err)
	}

	mt := getDescriptorType(ctx, m)
	tableName, err := getTableName(m)
	if err != nil {
		return err
	}

	fields := module.getMessageFields(ctx, m, ix, "m.self")
	mestedFields := getNesteFields(ctx, fields, ix)

	var vf string

	if !fext.NestedOnly {
		vf, err = getVersioningField(m)
		if err != nil {
			return err
		}
	}

	c := &descriptorTemplateContext{
		Type:                        mt,
		ReceiverType:                mt,
		Fields:                      fields,
		NestedFields:                mestedFields,
		Indexes:                     module.getMessageIndexes(ctx, m, ix),
		Statistics:                  module.getMessageStatistics(ctx, m, ix),
		TableName:                   tableName,
		VersioningField:             vf,
		IsPartitioned:               fext.Partitioned,
		IsPartitionedByCreatedAt:    fext.PartitionedByCreatedAt,
		PartitionedByKsuidFieldName: fext.PartitionedByKsuidFieldName,
		PartitionDateRange:          "pgdb_v1.MessageOptions_" + fext.PartitionedByDateRange.String(),
		UsePkskv2Column:             fext.UsePkskv2Column,
	}

	return templates["descriptor.tmpl"].Execute(w, c)
}

func getDescriptorType(ctx pgsgo.Context, m pgs.Message) string {
	return "pgdbDescriptor" + ctx.Name(m).String()
}

type nestedFieldContext struct {
	GoName   string
	TypeName string
	Prefix   string
}

func getNestedFieldNames(fields []*fieldContext) []string {
	rv := make([]string, 0)
	for _, f := range fields {
		if !f.Nested {
			continue
		}
		rv = append(rv, f.GoName)
	}
	return rv
}

func getNesteFields(ctx pgsgo.Context, fields []*fieldContext, ix *importTracker) []*nestedFieldContext {
	rv := make([]*nestedFieldContext, 0)
	for _, f := range fields {
		if !f.Nested {
			continue
		}
		rv = append(rv, &nestedFieldContext{
			GoName:   f.GoName,
			Prefix:   strconv.FormatInt(int64(*f.Field.Descriptor().Number), 10) + "$",
			TypeName: ix.Type(f.Field).String(),
		})
	}
	return rv
}
