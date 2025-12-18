package pgdb

import (
	"io"
	"strconv"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
)

type autovacuumTemplateContext struct {
	VacuumThreshold    *int32
	VacuumScaleFactor  *float32
	AnalyzeThreshold   *int32
	AnalyzeScaleFactor *float32
	VacuumCostDelay    *int32
	VacuumCostLimit    *int32
	FreezeMinAge       *int64
	FreezeMaxAge       *int64
	FreezeTableAge     *int64
	Fillfactor         *int32
	ToastTupleTarget   *int32
	Enabled            *bool
	HasEnabled         bool
}

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
	HasAutovacuum               bool
	Autovacuum                  *autovacuumTemplateContext
}

func (module *Module) renderDescriptor(ctx pgsgo.Context, w io.Writer, in pgs.File, m pgs.Message, ix *importTracker) error {
	ix.PGDBV1 = true
	ix.ProtobufReflect = true // For proto field metadata in Column struct
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

	if !fext.GetNestedOnly() {
		vf, err = getVersioningField(m)
		if err != nil {
			return err
		}
	}

	// Build autovacuum context if configured
	var autovacCtx *autovacuumTemplateContext
	hasAutovacuum := false
	if av := fext.GetAutovacuum(); av != nil {
		autovacCtx = &autovacuumTemplateContext{}
		if av.HasVacuumThreshold() {
			v := av.GetVacuumThreshold()
			autovacCtx.VacuumThreshold = &v
			hasAutovacuum = true
		}
		if av.HasVacuumScaleFactor() {
			v := av.GetVacuumScaleFactor()
			autovacCtx.VacuumScaleFactor = &v
			hasAutovacuum = true
		}
		if av.HasAnalyzeThreshold() {
			v := av.GetAnalyzeThreshold()
			autovacCtx.AnalyzeThreshold = &v
			hasAutovacuum = true
		}
		if av.HasAnalyzeScaleFactor() {
			v := av.GetAnalyzeScaleFactor()
			autovacCtx.AnalyzeScaleFactor = &v
			hasAutovacuum = true
		}
		if av.HasVacuumCostDelay() {
			v := av.GetVacuumCostDelay()
			autovacCtx.VacuumCostDelay = &v
			hasAutovacuum = true
		}
		if av.HasVacuumCostLimit() {
			v := av.GetVacuumCostLimit()
			autovacCtx.VacuumCostLimit = &v
			hasAutovacuum = true
		}
		if av.HasFreezeMinAge() {
			v := av.GetFreezeMinAge()
			autovacCtx.FreezeMinAge = &v
			hasAutovacuum = true
		}
		if av.HasFreezeMaxAge() {
			v := av.GetFreezeMaxAge()
			autovacCtx.FreezeMaxAge = &v
			hasAutovacuum = true
		}
		if av.HasFreezeTableAge() {
			v := av.GetFreezeTableAge()
			autovacCtx.FreezeTableAge = &v
			hasAutovacuum = true
		}
		if av.HasFillfactor() {
			v := av.GetFillfactor()
			autovacCtx.Fillfactor = &v
			hasAutovacuum = true
		}
		if av.HasToastTupleTarget() {
			v := av.GetToastTupleTarget()
			autovacCtx.ToastTupleTarget = &v
			hasAutovacuum = true
		}
		if av.HasEnabled() {
			v := av.GetEnabled()
			autovacCtx.Enabled = &v
			autovacCtx.HasEnabled = true
			hasAutovacuum = true
		}
		if hasAutovacuum {
			ix.ProtobufProto = true // Need proto package for proto.Int32, etc.
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
		IsPartitioned:               fext.GetPartitioned(),
		IsPartitionedByCreatedAt:    fext.GetPartitionedByCreatedAt(),
		PartitionedByKsuidFieldName: fext.GetPartitionedByKsuidFieldName(),
		PartitionDateRange:          "pgdb_v1.MessageOptions_" + fext.GetPartitionedByDateRange().String(),
		HasAutovacuum:               hasAutovacuum,
		Autovacuum:                  autovacCtx,
	}

	return templates["descriptor.tmpl"].Execute(w, c)
}

func getDescriptorType(ctx pgsgo.Context, m pgs.Message) string {
	return "pgdbDescriptor" + ctx.Name(m).String()
}

type nestedFieldContext struct {
	GoName    string
	TypeName  string
	Prefix    string
	FieldNum  int32  // Proto field number
	FieldName string // Proto field name (snake_case)
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
		fieldNum := *f.Field.Descriptor().Number
		rv = append(rv, &nestedFieldContext{
			GoName:    f.GoName,
			Prefix:    strconv.FormatInt(int64(fieldNum), 10) + "$",
			TypeName:  ix.Type(f.Field).String(),
			FieldNum:  fieldNum,
			FieldName: f.Field.Name().LowerSnakeCase().String(),
		})
	}
	return rv
}
