package pgdb

import (
	"fmt"
	"io"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type messageTemplateContext struct {
	ReceiverType            string
	MessageType             string
	DescriptorType          string
	Fields                  []*fieldContext
	SearchFields            []*searchFieldContext
	Indexes                 []*indexContext
	WantRecordStringBuilder bool
}

func (module *Module) renderMessage(ctx pgsgo.Context, w io.Writer, in pgs.File, m pgs.Message, ix *importTracker) error {
	ext := pgdb_v1.MessageOptions{}
	_, err := m.Extension(pgdb_v1.E_Msg, &ext)
	if err != nil {
		panic(fmt.Errorf("pgdb: getFieldIndexes: failed to extract Message extension from '%s': %w", m.FullyQualifiedName(), err))
	}

	ix.PGDBV1 = true
	ix.GoquExp = true
	wantRecordStringBuilder := false
	if !ext.NestedOnly {
		// used by pk/sk builder
		wantRecordStringBuilder = true
		ix.Strings = true
		// used by pb_data
		ix.ProtobufProto = true
	}
	fields := module.getMessageFields(ctx, m, ix, "m.self")

	c := &messageTemplateContext{
		ReceiverType:            ctx.Name(m).String(),
		MessageType:             getMessageType(ctx, m),
		DescriptorType:          getDescriptorType(ctx, m),
		Fields:                  fields,
		SearchFields:            getSearchFields(ctx, m),
		Indexes:                 module.getMessageIndexes(ctx, m, ix),
		WantRecordStringBuilder: wantRecordStringBuilder,
	}
	return templates["message.tmpl"].Execute(w, c)
}

func getMessageType(ctx pgsgo.Context, m pgs.Message) string {
	return "pgdbMessage" + ctx.Name(m).String()
}

type varNamer struct {
	prefix string
	offset int
}

func (fn *varNamer) Next() *varNamer {
	return &varNamer{offset: fn.offset + 1, prefix: fn.prefix}
}

func (fn *varNamer) String() string {
	return fmt.Sprintf("%s%d", fn.prefix, fn.offset)
}

func (module *Module) getMessageFields(ctx pgsgo.Context, m pgs.Message, ix *importTracker, goPrefix string) []*fieldContext {
	fields := m.Fields()
	rv := make([]*fieldContext, 0, len(fields))
	ix.ProtobufEncodingJSON = true
	cf, err := getCommonFields(ctx, m, ix)
	if err != nil {
		panic(err)
	}
	rv = append(rv, cf...)
	vn := &varNamer{prefix: "v", offset: 0}

	tenantIdField := ""
	fext := pgdb_v1.MessageOptions{}
	_, err = m.Extension(pgdb_v1.E_Msg, &fext)
	if err != nil {
		panic(err)
	}

	if !fext.NestedOnly {
		tenantIdField, err = getTenantIDField(m)
		if err != nil {
			panic(err)
		}
	}

	rv = append(rv, module.getMessageFieldsInner(ctx, m, fields, vn, tenantIdField, ix, goPrefix)...)

	vn = &varNamer{prefix: "oneof", offset: 0}
	for _, oneof := range m.RealOneOfs() {
		vn = vn.Next()
		fieldRep := module.getOneOf(ctx, oneof, vn, ix, goPrefix)
		if fieldRep != nil {
			rv = append(rv, fieldRep)
		}
	}

	// for _, field := range rv {
	// 	if field.Field == nil {
	// 		continue
	// 	}
	// 	ix.AddProtoEntity(field.Field)
	// }
	return rv
}

func (module *Module) getMessageFieldsInner(ctx pgsgo.Context, m pgs.Message, fields []pgs.Field, vn *varNamer, tenantIdField string, ix *importTracker, goPrefix string) []*fieldContext {
	rv := make([]*fieldContext, 0, len(fields))
	for _, field := range fields {
		// tenant_id done via common fields
		if tenantIdField == field.Name().LowerSnakeCase().String() {
			continue
		}
		vn = vn.Next()
		fieldRep := module.getField(ctx, field, vn, ix, goPrefix)
		if fieldRep != nil {
			rv = append(rv, fieldRep)
		}
	}
	return rv
}
