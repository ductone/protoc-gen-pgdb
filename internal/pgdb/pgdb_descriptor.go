package pgdb

import (
	"io"
	"strconv"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
)

type storageParamsTemplateContext struct {
	AutovacuumVacuumThreshold    *int32
	AutovacuumVacuumScaleFactor  *float32
	AutovacuumAnalyzeThreshold   *int32
	AutovacuumAnalyzeScaleFactor *float32
	AutovacuumVacuumCostDelay    *int32
	AutovacuumVacuumCostLimit    *int32
	AutovacuumFreezeMinAge       *int64
	AutovacuumFreezeMaxAge       *int64
	AutovacuumFreezeTableAge     *int64
	Fillfactor                   *int32
	ToastTupleTarget             *int32
	AutovacuumEnabled            *bool
	HasAutovacuumEnabled         bool
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
	HasStorageParameters        bool
	StorageParameters           *storageParamsTemplateContext
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

	// Build storage parameters context if configured
	var spCtx *storageParamsTemplateContext
	hasStorageParams := false
	if sp := fext.GetStorageParameters(); sp != nil {
		spCtx = &storageParamsTemplateContext{}
		if sp.HasAutovacuumVacuumThreshold() {
			v := sp.GetAutovacuumVacuumThreshold()
			spCtx.AutovacuumVacuumThreshold = &v
			hasStorageParams = true
		}
		if sp.HasAutovacuumVacuumScaleFactor() {
			v := sp.GetAutovacuumVacuumScaleFactor()
			spCtx.AutovacuumVacuumScaleFactor = &v
			hasStorageParams = true
		}
		if sp.HasAutovacuumAnalyzeThreshold() {
			v := sp.GetAutovacuumAnalyzeThreshold()
			spCtx.AutovacuumAnalyzeThreshold = &v
			hasStorageParams = true
		}
		if sp.HasAutovacuumAnalyzeScaleFactor() {
			v := sp.GetAutovacuumAnalyzeScaleFactor()
			spCtx.AutovacuumAnalyzeScaleFactor = &v
			hasStorageParams = true
		}
		if sp.HasAutovacuumVacuumCostDelay() {
			v := sp.GetAutovacuumVacuumCostDelay()
			spCtx.AutovacuumVacuumCostDelay = &v
			hasStorageParams = true
		}
		if sp.HasAutovacuumVacuumCostLimit() {
			v := sp.GetAutovacuumVacuumCostLimit()
			spCtx.AutovacuumVacuumCostLimit = &v
			hasStorageParams = true
		}
		if sp.HasAutovacuumFreezeMinAge() {
			v := sp.GetAutovacuumFreezeMinAge()
			spCtx.AutovacuumFreezeMinAge = &v
			hasStorageParams = true
		}
		if sp.HasAutovacuumFreezeMaxAge() {
			v := sp.GetAutovacuumFreezeMaxAge()
			spCtx.AutovacuumFreezeMaxAge = &v
			hasStorageParams = true
		}
		if sp.HasAutovacuumFreezeTableAge() {
			v := sp.GetAutovacuumFreezeTableAge()
			spCtx.AutovacuumFreezeTableAge = &v
			hasStorageParams = true
		}
		if sp.HasFillfactor() {
			v := sp.GetFillfactor()
			spCtx.Fillfactor = &v
			hasStorageParams = true
		}
		if sp.HasToastTupleTarget() {
			v := sp.GetToastTupleTarget()
			spCtx.ToastTupleTarget = &v
			hasStorageParams = true
		}
		if sp.HasAutovacuumEnabled() {
			v := sp.GetAutovacuumEnabled()
			spCtx.AutovacuumEnabled = &v
			spCtx.HasAutovacuumEnabled = true
			hasStorageParams = true
		}
		if hasStorageParams {
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
		HasStorageParameters:        hasStorageParams,
		StorageParameters:           spCtx,
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
