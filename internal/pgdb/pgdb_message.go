package pgdb

import (
	"fmt"
	"io"
	"strings"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
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
	if !ext.GetNestedOnly() {
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

func (module *Module) getMessageFieldsDeep(ctx pgsgo.Context, m pgs.Message, ix *importTracker, goPrefix string, dbPrefix string, humanPrefix string) []*fieldContext {
	// Call the internal implementation with empty proto path (for top-level) and no oneof
	return module.getMessageFieldsDeepInternal(ctx, m, ix, goPrefix, dbPrefix, humanPrefix, nil, nil, "")
}

func (module *Module) getMessageFieldsDeepInternal(ctx pgsgo.Context, m pgs.Message, ix *importTracker, goPrefix string, dbPrefix string, humanPrefix string, protoPath []int32, protoNamePath []string, oneofName string) []*fieldContext {
	fields := m.Fields()
	rv := make([]*fieldContext, 0, len(fields))
	tenantIdField := "tenant_id"

	// only top level embed gets a common field and has a full protoc API
	if dbPrefix == "" {
		ix.ProtobufEncodingJSON = true
		cfs, err := getCommonFields(ctx, m, ix)
		if err != nil {
			panic(err)
		}
		for _, cf := range cfs {
			cf.DBFieldNameDeep = cf.DB.Name
			rv = append(rv, cf)
		}

		fext := pgdb_v1.MessageOptions{}
		_, err = m.Extension(pgdb_v1.E_Msg, &fext)
		if err != nil {
			panic(err)
		}

		if !fext.GetNestedOnly() {
			tenantIdField, err = getTenantIDField(m)
			if err != nil {
				panic(err)
			}
		}
	}

	vn := &varNamer{prefix: "oneof", offset: 0}
	for _, oneof := range m.RealOneOfs() {
		vn = vn.Next()
		fc := module.getOneOf(ctx, oneof, vn, ix)
		if fc == nil {
			continue
		}
		name, err := getColumnOneOfName(oneof)
		if err != nil {
			panic(err)
		}

		fc.DBFieldNameDeep = dbPrefix + name

		// For nested oneof discriminators, update the proto path prefix
		if len(protoPath) > 0 && fc.DB != nil {
			fc.DB.ProtoPath = buildProtoPath(protoNamePath, fc.DB.OneofName)
		}

		if humanPrefix != "" {
			fc.GoName = humanPrefix + fc.GoName
		}
		rv = append(rv, fc)
	}

	vn = &varNamer{prefix: "v", offset: 0}
	for _, field := range fields {
		// tenant_id done via common fields and shouldn't be nested
		if tenantIdField == field.Name().LowerSnakeCase().String() {
			continue
		}
		vn = vn.Next()
		fc, err := module.getFieldSafe(ctx, field, vn, ix, goPrefix)
		if fc == nil || err != nil {
			continue
		}
		name, err := getColumnName(fc.Field)
		if err != nil {
			panic(err)
		}

		fc.DBFieldNameDeep = dbPrefix + name

		// Build the full proto path for this field
		fieldNum := int32(field.Descriptor().GetNumber())
		fieldName := field.Name().LowerSnakeCase().String()
		newProtoPath := appendInt32Slice(protoPath, fieldNum)
		newProtoNamePath := appendStringSlice(protoNamePath, fieldName)

		// Determine oneofName for this field - either inherited from parent or from this field's oneof
		fieldOneofName := oneofName
		if oo := field.OneOf(); oo != nil {
			fieldOneofName = oo.Name().LowerSnakeCase().String()
		}

		// Update the Column's proto paths
		if fc.DB != nil {
			fc.DB.ProtoFieldPath = newProtoPath
			fc.DB.ProtoPath = strings.Join(newProtoNamePath, ".")
			fc.DB.OneofName = fieldOneofName
		}

		rv = append(rv, fc)
		var embededMessage pgs.Message
		if fc.Field != nil {
			embededMessage = fc.Field.Type().Embed()
		}

		if embededMessage == nil {
			if humanPrefix != "" {
				fc.GoName = humanPrefix + fc.GoName
			}
			continue
		}
		// NOTE: humanPrefixes need to avoid exponential growth of prefixes (for two deep or lower).
		nextHumanPrefix := humanPrefix + fc.GoName
		if humanPrefix != "" {
			fc.GoName = humanPrefix
		}

		pre := getNestedName(fc.Field)

		// Recurse with extended proto paths, passing the oneof name for nested fields
		rv = append(rv, module.getMessageFieldsDeepInternal(ctx, embededMessage, ix, goPrefix, dbPrefix+pre, nextHumanPrefix, newProtoPath, newProtoNamePath, fieldOneofName)...)
	}

	return rv
}

// appendInt32Slice creates a new slice with the value appended (doesn't modify original).
func appendInt32Slice(slice []int32, val int32) []int32 {
	result := make([]int32, len(slice)+1)
	copy(result, slice)
	result[len(slice)] = val
	return result
}

// appendStringSlice creates a new slice with the value appended (doesn't modify original).
func appendStringSlice(slice []string, val string) []string {
	result := make([]string, len(slice)+1)
	copy(result, slice)
	result[len(slice)] = val
	return result
}

// buildProtoPath constructs a dot-delimited path from parts.
func buildProtoPath(parts []string, suffix string) string {
	if len(parts) == 0 {
		return suffix
	}
	return strings.Join(parts, ".") + "." + suffix
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

	if !fext.GetNestedOnly() {
		tenantIdField, err = getTenantIDField(m)
		if err != nil {
			panic(err)
		}
	}

	rv = append(rv, module.getMessageFieldsInner(ctx, fields, vn, tenantIdField, ix, goPrefix)...)

	vn = &varNamer{prefix: "oneof", offset: 0}
	for _, oneof := range m.RealOneOfs() {
		vn = vn.Next()
		fieldRep := module.getOneOf(ctx, oneof, vn, ix)
		if fieldRep != nil {
			rv = append(rv, fieldRep)
		}
	}

	return rv
}

func (module *Module) getMessageFieldsInner(ctx pgsgo.Context, fields []pgs.Field, vn *varNamer, tenantIdField string, ix *importTracker, goPrefix string) []*fieldContext {
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
