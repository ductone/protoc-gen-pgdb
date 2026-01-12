package pgdb

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
	"google.golang.org/protobuf/reflect/protoreflect"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
)

const (
	pgTypeJSONB = "jsonb"
	pgTypeInt4  = "int4"
	pgTypeBool  = "bool"
	pgTypeBit   = "bit"
)

type fieldContext struct {
	// denotes a realized/virtual field that comes from multiple fields. in this case, F is nil.
	IsVirtual         bool
	ExcludeNested     bool
	GoName            string
	Field             pgs.Field
	Nested            bool
	DB                *pgdb_v1.Column
	DataType          *pgtype.Type
	Convert           FieldConverter
	QueryTypeName     string
	DBFieldNameDeep   string
	V17FieldOverrides *V17FieldOverrides
}

type V17FieldOverrides struct {
	Collation               string
	ClearOverrideExpression bool
	Disabled                bool
}

type FieldConverter interface {
	GoType() (string, error)
	CodeForValue() (string, error)
	VarForValue() (string, error)
	VarForAppend() (string, error)
	EnumForValue() (string, error)
}

func (module *Module) getFieldSafe(ctx pgsgo.Context, f pgs.Field, vn *varNamer, ix *importTracker, goPrefix string) (*fieldContext, error) {
	ext := pgdb_v1.FieldOptions{}
	_, err := f.Extension(pgdb_v1.E_Options, &ext)
	if err != nil {
		return nil, fmt.Errorf("pgdb: getField: failed to extract Message extension from '%s': %w", f.FullyQualifiedName(), err)
	}

	if ext.GetMessageBehavior() == pgdb_v1.FieldOptions_MESSAGE_BEHAVIOR_OMIT {
		// explict option to just not store this in postgres
		return nil, nil
	}

	isArray := f.Type().ProtoLabel() == pgs.Repeated && !f.Type().IsMap()
	pt := f.Type().ProtoType()

	pgColName, err := getColumnName(f)
	if err != nil {
		return nil, fmt.Errorf("pgdb: getColumnName failed for: %v: %s (of type %s)",
			pt, f.FullyQualifiedName(), f.Descriptor().GetType())
	}

	convertDef := &fieldConvert{
		goPrefix: goPrefix,
		ctx:      ctx,
		ix:       ix,
		F:        f,
		varName:  vn.String(),
	}
	nullable := true
	defaultValue := ""
	overrideExpression := ""
	// https://developers.google.com/protocol-buffers/docs/proto3#scalar
	switch pt {
	case pgs.DoubleT:
		// aka float64
		convertDef.PostgresTypeName = "float8"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = gtFloat64
		defaultValue = "0.0"
		nullable = false
	case pgs.FloatT:
		// aka float32
		convertDef.PostgresTypeName = "float4"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = gtFloat32
		defaultValue = "0.0"
		nullable = false
	case pgs.Int32T, pgs.SInt32, pgs.SFixed32:
		convertDef.PostgresTypeName = pgTypeInt4
		convertDef.IsArray = isArray
		convertDef.TypeConversion = gtInt32
		defaultValue = "0"
		nullable = false
	case pgs.Int64T, pgs.SInt64, pgs.SFixed64:
		convertDef.PostgresTypeName = "int8"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = gtInt64
		defaultValue = "0"
		nullable = false
	case pgs.UInt32T, pgs.Fixed32T:
		convertDef.PostgresTypeName = "int8"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = gtUint32
		defaultValue = "0"
		nullable = false
	case pgs.UInt64T, pgs.Fixed64T:
		// not ideal, but postgres only has signed types for int8.
		convertDef.PostgresTypeName = "numeric"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = gtUint64
		defaultValue = "0"
		nullable = false
	case pgs.BoolT:
		convertDef.PostgresTypeName = pgTypeBool
		convertDef.IsArray = isArray
		convertDef.TypeConversion = gtBool
		defaultValue = "false"
		nullable = false
	case pgs.StringT:
		switch ext.GetMessageBehavior() {
		case pgdb_v1.FieldOptions_MESSAGE_BEHAVIOR_UNSPECIFIED:
			// TODO(pquerna): annotations for max size
			convertDef.PostgresTypeName = "text"
			convertDef.IsArray = isArray
			convertDef.TypeConversion = gtString
			convertDef.FullTextType = ext.GetFullTextType()
			convertDef.FullTextWeight = ext.GetFullTextWeight()
			defaultValue = "''"
			nullable = false
		case pgdb_v1.FieldOptions_MESSAGE_BEHAVIOR_INET_ADDR:
			if isArray {
				return nil, fmt.Errorf("pgdb: unsupported field type: %v: %s: arrays of inet addr not supported", pt, f.FullyQualifiedName())
			}

			convertDef.PostgresTypeName = "inet"
			convertDef.TypeConversion = gtInetAddr
			defaultValue = "NULL"
			nullable = true
		default:
			return nil, fmt.Errorf("pgdb: unsupported field type: %v: %s: MessageBehavior not supported on string type", pt, f.FullyQualifiedName())
		}
	case pgs.MessageT:
		switch ext.GetMessageBehavior() {
		case pgdb_v1.FieldOptions_MESSAGE_BEHAVIOR_UNSPECIFIED:
			switch f.Descriptor().GetTypeName() {
			case ".google.protobuf.Any":
				if isArray {
					return nil, fmt.Errorf("pgdb: unsupported field type: %v: %s: repeated Any not supported", pt, f.FullyQualifiedName())
				}
				convertDef.PostgresTypeName = pgTypeJSONB
				convertDef.TypeConversion = gtPbWktAny
			case ".google.protobuf.Timestamp":
				convertDef.PostgresTypeName = "timestamptz"
				convertDef.IsArray = isArray
				convertDef.TypeConversion = gtPbWktTimestamp
			case ".google.protobuf.Duration":
				convertDef.PostgresTypeName = "interval"
				convertDef.IsArray = isArray
				convertDef.TypeConversion = gtPbWktDuration
			case ".google.protobuf.Struct":
				convertDef.IsArray = isArray
				convertDef.PostgresTypeName = pgTypeJSONB
				convertDef.TypeConversion = gtPbWktStruct
			case ".google.protobuf.BoolValue":
				convertDef.IsArray = isArray
				convertDef.PostgresTypeName = pgTypeBool
				convertDef.TypeConversion = gtPbWktBoolValue
			case ".google.protobuf.StringValue":
				convertDef.IsArray = isArray
				convertDef.PostgresTypeName = "text"
				convertDef.TypeConversion = gtPbWktStringValue

			default:
				if isArray {
					return nil, fmt.Errorf("pgdb: unsupported message field type: %v: %s (of type %s): Arrays cannot be nested; consider jsonb",
						pt, f.FullyQualifiedName(), f.Descriptor().GetType())
				}
				convertDef.TypeConversion = gtPbNestedMsg
				convertDef.NestedPrefix = getNestedName(f)
			}
		case pgdb_v1.FieldOptions_MESSAGE_BEHAVIOR_OMIT:
			// explict option to just not store this in postgres
			return nil, nil
		case pgdb_v1.FieldOptions_MESSAGE_BEHAVIOR_EXPAND:
			if isArray {
				return nil, fmt.Errorf("pgdb: unsupported message field type: %v: %s (of type %s): Arrays cannot be nested; consider jsonb",
					pt, f.FullyQualifiedName(), f.Descriptor().GetType())
			}
			convertDef.TypeConversion = gtPbNestedMsg
			convertDef.NestedPrefix = getNestedName(f)
		case pgdb_v1.FieldOptions_MESSAGE_BEHAVIOR_JSONB:
			convertDef.IsArray = isArray
			convertDef.PostgresTypeName = pgTypeJSONB
			convertDef.TypeConversion = gtPbGenericMsg
		case pgdb_v1.FieldOptions_MESSAGE_BEHAVIOR_VECTOR:
			return nil, nil
		default:
			return nil, fmt.Errorf("pgdb: unsupported message field type: %v: %s (of type %s)",
				pt, f.FullyQualifiedName(), f.Descriptor().GetType())
		}
	case pgs.BytesT:
		switch ext.GetMessageBehavior() {
		case pgdb_v1.FieldOptions_MESSAGE_BEHAVIOR_BITS:
			length := ext.GetBitsSize()
			if length == 0 {
				return nil, fmt.Errorf("pgdb: bits size must be greater than 0: %s", f.FullyQualifiedName())
			}
			convertDef.IsArray = isArray
			convertDef.PostgresTypeName = pgTypeBit
			convertDef.TypeConversion = gtBits
			convertDef.ExpectedBytesLen = length / 8
			overrideExpression = fmt.Sprintf("bit(%d)", length)
		default:
			convertDef.IsArray = isArray
			convertDef.PostgresTypeName = "bytea"
			convertDef.TypeConversion = gtBytes
		}
	case pgs.EnumT:
		convertDef.PostgresTypeName = pgTypeInt4
		convertDef.IsArray = isArray
		convertDef.TypeConversion = gtEnum
		defaultValue = "0" // NOTE: assumes a "UNSPECIFIED = 0" in the enum...
		nullable = false

	case pgs.GroupT:
		return nil, fmt.Errorf("pgdb: unsupported field type: Group: %s", f.FullyQualifiedName())
	default:
		return nil, fmt.Errorf("pgdb: unsupported field type: %v: %s (of type %s)",
			pt, f.FullyQualifiedName(), f.Descriptor().GetType())
	}

	if isArray {
		nullable = true
		defaultValue = ""
		ix.XPQ = true
		convertDef.PostgresTypeName = "_" + convertDef.PostgresTypeName
	}

	rv := &fieldContext{
		IsVirtual:     false,
		GoName:        ctx.Name(f).String(),
		Field:         f,
		Convert:       convertDef,
		QueryTypeName: ctx.Name(f.Message()).String() + ctx.Name(f).String() + "QueryType",
	}

	if convertDef.TypeConversion != gtPbNestedMsg {
		dbTypeRef := pgDataTypeForName(convertDef.PostgresTypeName)
		if !nullable && defaultValue == "" {
			return nil, fmt.Errorf("pgdb: nullable column with no default: %s (%s)", pgColName, dbTypeRef.Name)
		}
		rv.DB = &pgdb_v1.Column{
			Name:               pgColName,
			Type:               dbTypeRef.Name,
			Nullable:           nullable,
			Default:            defaultValue,
			Collation:          ext.GetCollation(),
			OverrideExpression: overrideExpression,
			// Proto metadata - will be extended by caller for nested fields
			SourceKind:     pgdb_v1.ColumnSourceKind_PROTO_FIELD,
			ProtoFieldPath: []int32{f.Descriptor().GetNumber()},
			ProtoPath:      f.Name().LowerSnakeCase().String(),
			ProtoKind:      protoTypeToKind(pt),
			ProtoTypeName:  getProtoTypeName(f),
			IsRepeated:     isArray,
		}
		rv.DataType = dbTypeRef

		// Handle sequence option for auto-incrementing columns
		if seqOpts := ext.GetSequence(); seqOpts != nil {
			// Validate that sequence is only used on integer types
			if !isIntegerType(pt) {
				return nil, fmt.Errorf("pgdb: sequence option can only be used on integer fields: %s (type %s)", f.FullyQualifiedName(), pt.String())
			}
			if isArray {
				return nil, fmt.Errorf("pgdb: sequence option cannot be used on repeated fields: %s", f.FullyQualifiedName())
			}
			rv.DB.Sequence = protoSequenceToColumnSequence(seqOpts)
			// Clear default value since identity columns don't use defaults
			rv.DB.Default = ""
			rv.DB.Nullable = false
		}

		if ext.GetKsuid() {
			rv.V17FieldOverrides = &V17FieldOverrides{
				Collation: "C",
			}
		} else if isIDField(ctx, f) {
			rv.V17FieldOverrides = &V17FieldOverrides{
				Collation: "C",
			}
		}
	} else {
		rv.Nested = true
	}
	return rv, nil
}

func isIDField(ctx pgsgo.Context, f pgs.Field) bool {
	isID := strings.HasSuffix(ctx.Name(f).LowerSnakeCase().String(), "_id") || strings.TrimPrefix(ctx.Name(f).LowerSnakeCase().String(), "pb$") == "id"
	return isID && f.Type().ProtoType() == pgs.StringT
}

// Panics on error.
func (module *Module) getField(ctx pgsgo.Context, f pgs.Field, vn *varNamer, ix *importTracker, goPrefix string) *fieldContext {
	fc, err := module.getFieldSafe(ctx, f, vn, ix, goPrefix)
	if err != nil {
		panic(err)
	}
	return fc
}

func getCommonFields(ctx pgsgo.Context, m pgs.Message, ix *importTracker) ([]*fieldContext, error) {
	vn := &varNamer{prefix: "cfv", offset: 0}
	_ = vn
	vcDataType := pgDataTypeForName("varchar")
	fext := pgdb_v1.MessageOptions{}
	_, err := m.Extension(pgdb_v1.E_Msg, &fext)
	if err != nil {
		panic(err)
	}

	rv := []*fieldContext{}
	// nested only currently don't have any of the common fields.
	if fext.GetNestedOnly() {
		return nil, nil
	}

	byteaDataType := pgDataTypeForName("bytea")

	// Get the actual tenant_id field info for proto metadata
	tenantIdFieldName, err := getTenantIDField(m)
	if err != nil {
		return nil, err
	}
	tenantIdProtoField := fieldByName(m, tenantIdFieldName)

	tenantIdField := &fieldContext{
		ExcludeNested: true,
		IsVirtual:     true,
		DB: &pgdb_v1.Column{
			Name:           "tenant_id",
			Type:           vcDataType.Name,
			Nullable:       false,
			SourceKind:     pgdb_v1.ColumnSourceKind_TENANT,
			ProtoFieldPath: []int32{tenantIdProtoField.Descriptor().GetNumber()},
			ProtoPath:      tenantIdFieldName,
			ProtoKind:      protoreflect.StringKind,
		},
		GoName:   "TenantId",
		DataType: vcDataType,
		Convert: &tenantIdDataConvert{
			ctx:     ctx,
			Message: m,
			VarName: vn.String(),
		},
		QueryTypeName: ctx.Name(m).String() + "TenantId" + "QueryType",
		V17FieldOverrides: &V17FieldOverrides{
			Collation: "C",
		},
	}
	rv = append(rv, tenantIdField)

	vn = vn.Next()
	pkskField := &fieldContext{
		ExcludeNested: true,
		IsVirtual:     true,
		DB: &pgdb_v1.Column{
			Name:               "pksk",
			Type:               vcDataType.Name,
			Nullable:           false,
			OverrideExpression: "varchar GENERATED ALWAYS AS (pb$pk || '|' || pb$sk) STORED",
			SourceKind:         pgdb_v1.ColumnSourceKind_PRIMARY_KEY,
			ProtoKind:          protoreflect.StringKind,
		},
		GoName:   "PKSK",
		DataType: vcDataType,
		Convert: &pkskDataConvert{
			ctx: ctx,
		},
		QueryTypeName: ctx.Name(m).String() + "PKSK" + "QueryType",
		V17FieldOverrides: &V17FieldOverrides{
			Collation:               "C",
			ClearOverrideExpression: true,
		},
	}
	rv = append(rv, pkskField)

	vn = vn.Next()
	pkField := &fieldContext{
		ExcludeNested: true,
		IsVirtual:     true,
		DB: &pgdb_v1.Column{
			Name:       "pk",
			Type:       vcDataType.Name,
			Nullable:   false,
			SourceKind: pgdb_v1.ColumnSourceKind_PRIMARY_KEY,
			ProtoKind:  protoreflect.StringKind,
		},
		GoName:   "PK",
		DataType: vcDataType,
		Convert: &dynamoKeyDataConvert{
			ix:      ix,
			ctx:     ctx,
			VarName: vn.String(),
			Message: m,
			KeyType: DynamoKeyTypePartition,
		},
		QueryTypeName: ctx.Name(m).String() + "PK" + "QueryType",
		V17FieldOverrides: &V17FieldOverrides{
			Collation: "C",
		},
	}
	rv = append(rv, pkField)

	vn = vn.Next()
	skField := &fieldContext{
		ExcludeNested: true,
		IsVirtual:     true,
		DB: &pgdb_v1.Column{
			Name:       "sk",
			Type:       vcDataType.Name,
			Nullable:   false,
			SourceKind: pgdb_v1.ColumnSourceKind_PRIMARY_KEY,
			ProtoKind:  protoreflect.StringKind,
		},
		GoName:   "SK",
		DataType: vcDataType,
		Convert: &dynamoKeyDataConvert{
			ix:      ix,
			ctx:     ctx,
			VarName: vn.String(),
			Message: m,
			KeyType: DynamoKeyTypeSort,
		},
		QueryTypeName: ctx.Name(m).String() + "SK" + "QueryType",
		V17FieldOverrides: &V17FieldOverrides{
			Collation: "C",
		},
	}
	rv = append(rv, skField)

	vn = vn.Next()
	pkskv2Field := &fieldContext{
		ExcludeNested: true,
		IsVirtual:     true,
		DB: &pgdb_v1.Column{
			Name:       "pkskv2",
			Type:       vcDataType.Name,
			Nullable:   true,
			Collation:  "C",
			SourceKind: pgdb_v1.ColumnSourceKind_PRIMARY_KEY,
			ProtoKind:  protoreflect.StringKind,
		},
		GoName:   "PKSKV2",
		DataType: vcDataType,
		Convert: &pkskDataConvert{
			ctx: ctx,
		},
		QueryTypeName: ctx.Name(m).String() + "PKSKV2" + "QueryType",
		V17FieldOverrides: &V17FieldOverrides{
			Disabled: true,
		},
	}
	rv = append(rv, pkskv2Field)

	// https://github.com/jackc/pgtype/issues/150
	// tsvector is not in-tree.  but we use to_tsvector() when inserting, so we just need to have the right type name
	// in the Field{} struct.
	// tsDataType, _ := pgDataTypeForName("tsvector")
	vn = vn.Next()
	ftsDataField := &fieldContext{
		ExcludeNested: true,
		IsVirtual:     true,
		DB: &pgdb_v1.Column{
			Name:       "fts_data",
			Type:       "tsvector",
			Nullable:   true,
			SourceKind: pgdb_v1.ColumnSourceKind_SEARCH,
		},
		GoName:   "FTSData",
		DataType: nil,
		Convert: &ftsDataConvert{
			ctx:     ctx,
			m:       m,
			VarName: vn.String(),
		},
		QueryTypeName: ctx.Name(m).String() + "FTSData" + "QueryType",
	}
	rv = append(rv, ftsDataField)

	vn = vn.Next()
	pbDataField := &fieldContext{
		ExcludeNested: true,
		IsVirtual:     true,
		DB: &pgdb_v1.Column{
			Name:       "pb_data",
			Type:       byteaDataType.Name,
			Nullable:   false,
			SourceKind: pgdb_v1.ColumnSourceKind_DATA,
			ProtoKind:  protoreflect.BytesKind,
		},
		GoName:   "PBData",
		DataType: byteaDataType,
		Convert: &pbDataConvert{
			VarName: vn.String(),
		},
		QueryTypeName: ctx.Name(m).String() + "PBData" + "QueryType",
	}
	rv = append(rv, pbDataField)

	// iterate message for vector behavior options
	for _, field := range m.Fields() {
		ext := pgdb_v1.FieldOptions{}
		_, err := field.Extension(pgdb_v1.E_Options, &ext)
		if err != nil {
			return nil, fmt.Errorf("pgdb: getField: failed to extract Message extension from '%s': %w", field.FullyQualifiedName(), err)
		}
		if ext.GetMessageBehavior() != pgdb_v1.FieldOptions_MESSAGE_BEHAVIOR_VECTOR {
			continue
		}

		enumField, floatField, err := GetFieldVectorShape(field)
		if err != nil {
			return nil, err
		}

		pgColName, err := getColumnName(field)
		if err != nil {
			panic(fmt.Errorf("pgdb: getColumnName failed for: %v: %s (of type %s)",
				field.Type().ProtoType(), field.FullyQualifiedName(), field.Descriptor().GetType()))
		}

		var unspecifiedEnum pgs.EnumValue
		// enum values
		for _, enumValue := range enumField.Type().Enum().Values() {
			if enumValue.Value() == 0 {
				// skip the zero value
				unspecifiedEnum = enumValue
				continue
			}

			enumExt := pgdb_v1.EnumValueOptions{}
			_, err := enumValue.Extension(pgdb_v1.E_Enum, &enumExt)
			if err != nil {
				return nil, fmt.Errorf("pgdb: getField: failed to extract enum extension from '%s': %w", enumValue.FullyQualifiedName(), err)
			}

			toTrim := strings.TrimSuffix(ctx.Name(unspecifiedEnum).String(), "_UNSPECIFIED")

			goNameString := ctx.Name(field).String() + strings.TrimPrefix(ctx.Name(enumValue).String(), toTrim)

			tempCtx := &fieldContext{
				ExcludeNested: true,
				IsVirtual:     true,
				DB: &pgdb_v1.Column{
					Name:               fmt.Sprintf("%s_%d", pgColName, enumValue.Value()),
					Type:               "vector",
					Nullable:           true,
					OverrideExpression: fmt.Sprintf("vector(%d)", enumExt.GetVectorSize()),
					SourceKind:         pgdb_v1.ColumnSourceKind_VECTOR,
					ProtoFieldPath:     []int32{field.Descriptor().GetNumber()},
					ProtoPath:          field.Name().LowerSnakeCase().String(),
					ProtoKind:          protoreflect.MessageKind,
				},
				GoName:   goNameString, // Generated go struct name
				DataType: nil,
				// new struct to implement this
				Convert: &pbVectorConvert{
					VarName:        vn.String(),
					EnumName:       ctx.Name(enumField).String(), // Generated enum name
					GoName:         ctx.Name(field).String(),     // Generated go struct name
					EnumModelValue: ix.EnumValue(field, enumValue).String(),
					FloatArrayName: ctx.Name(floatField).String(), // Generated float array name
				},
				QueryTypeName: ctx.Name(m).String() + goNameString + "QueryType",
			}
			rv = append(rv, tempCtx)
		}

		break
	}

	return rv, nil
}

// protoTypeToKind converts pgs.ProtoType to protoreflect.Kind.
func protoTypeToKind(pt pgs.ProtoType) protoreflect.Kind {
	switch pt {
	case pgs.DoubleT:
		return protoreflect.DoubleKind
	case pgs.FloatT:
		return protoreflect.FloatKind
	case pgs.Int32T:
		return protoreflect.Int32Kind
	case pgs.Int64T:
		return protoreflect.Int64Kind
	case pgs.UInt32T:
		return protoreflect.Uint32Kind
	case pgs.UInt64T:
		return protoreflect.Uint64Kind
	case pgs.SInt32:
		return protoreflect.Sint32Kind
	case pgs.SInt64:
		return protoreflect.Sint64Kind
	case pgs.Fixed32T:
		return protoreflect.Fixed32Kind
	case pgs.Fixed64T:
		return protoreflect.Fixed64Kind
	case pgs.SFixed32:
		return protoreflect.Sfixed32Kind
	case pgs.SFixed64:
		return protoreflect.Sfixed64Kind
	case pgs.BoolT:
		return protoreflect.BoolKind
	case pgs.StringT:
		return protoreflect.StringKind
	case pgs.BytesT:
		return protoreflect.BytesKind
	case pgs.EnumT:
		return protoreflect.EnumKind
	case pgs.MessageT:
		return protoreflect.MessageKind
	case pgs.GroupT:
		return protoreflect.GroupKind
	default:
		return 0 // Unknown
	}
}

// getProtoTypeName returns the fully qualified proto type name for enum/message fields.
// Returns empty string for scalar types.
func getProtoTypeName(f pgs.Field) string {
	ft := f.Type()

	// For repeated fields, look at element type
	if ft.IsRepeated() && !ft.IsMap() {
		el := ft.Element()
		if el.IsEnum() {
			return el.Enum().FullyQualifiedName()
		}
		if el.IsEmbed() {
			return el.Embed().FullyQualifiedName()
		}
		return ""
	}

	if ft.IsEnum() {
		return ft.Enum().FullyQualifiedName()
	}
	if ft.IsEmbed() {
		return ft.Embed().FullyQualifiedName()
	}
	return ""
}

// isIntegerType returns true if the proto type is an integer type suitable for sequences.
func isIntegerType(pt pgs.ProtoType) bool {
	switch pt {
	case pgs.Int32T, pgs.Int64T, pgs.UInt32T, pgs.UInt64T,
		pgs.SInt32, pgs.SInt64, pgs.Fixed32T, pgs.Fixed64T,
		pgs.SFixed32, pgs.SFixed64:
		return true
	default:
		return false
	}
}

// protoSequenceToColumnSequence converts proto SequenceOptions to Column SequenceOptions.
func protoSequenceToColumnSequence(seqOpts *pgdb_v1.FieldOptions_SequenceOptions) *pgdb_v1.SequenceOptions {
	if seqOpts == nil {
		return nil
	}
	rv := &pgdb_v1.SequenceOptions{
		Always:    seqOpts.GetAlways(),
		Start:     seqOpts.GetStart(),
		Increment: seqOpts.GetIncrement(),
		Cycle:     seqOpts.GetCycle(),
	}
	if seqOpts.HasMinValue() {
		minVal := seqOpts.GetMinValue()
		rv.MinValue = &minVal
	}
	if seqOpts.HasMaxValue() {
		maxVal := seqOpts.GetMaxValue()
		rv.MaxValue = &maxVal
	}
	if seqOpts.HasCache() {
		cache := seqOpts.GetCache()
		rv.Cache = &cache
	}
	return rv
}
