package pgdb

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
)

const (
	pgTypeJSONB = "jsonb"
	pgTypeInt4  = "int4"
	pgTypeBool  = "bool"
)

type fieldContext struct {
	// denotes a realized/virtual field that comes from multiple fields. in this case, F is nil.
	IsVirtual       bool
	ExcludeNested   bool
	GoName          string
	Field           pgs.Field
	Nested          bool
	DB              *pgdb_v1.Column
	DataType        *pgtype.Type
	Convert         FieldConverter
	QueryTypeName   string
	DBFieldNameDeep string
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
		convertDef.IsArray = isArray
		convertDef.PostgresTypeName = "bytea"
		convertDef.TypeConversion = gtBytes
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
			Name:      pgColName,
			Type:      dbTypeRef.Name,
			Nullable:  nullable,
			Default:   defaultValue,
			Collation: ext.GetCollation(),
		}
		rv.DataType = dbTypeRef
	} else {
		rv.Nested = true
	}
	return rv, nil
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

	tenantIdField := &fieldContext{
		ExcludeNested: true,
		IsVirtual:     true,
		DB: &pgdb_v1.Column{
			Name:     "tenant_id",
			Type:     vcDataType.Name,
			Nullable: false,
		},
		GoName:   "TenantId",
		DataType: vcDataType,
		Convert: &tenantIdDataConvert{
			ctx:     ctx,
			Message: m,
			VarName: vn.String(),
		},
		QueryTypeName: ctx.Name(m).String() + "TenantId" + "QueryType",
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
		},
		GoName:   "PKSK",
		DataType: vcDataType,
		Convert: &pkskDataConvert{
			ctx: ctx,
		},
		QueryTypeName: ctx.Name(m).String() + "PKSK" + "QueryType",
	}
	rv = append(rv, pkskField)

	vn = vn.Next()
	pkField := &fieldContext{
		ExcludeNested: true,
		IsVirtual:     true,
		DB: &pgdb_v1.Column{
			Name:     "pk",
			Type:     vcDataType.Name,
			Nullable: false,
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
	}
	rv = append(rv, pkField)

	vn = vn.Next()
	skField := &fieldContext{
		ExcludeNested: true,
		IsVirtual:     true,
		DB: &pgdb_v1.Column{
			Name:     "sk",
			Type:     vcDataType.Name,
			Nullable: false,
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
	}
	rv = append(rv, skField)

	vn = vn.Next()
	pkskv2Field := &fieldContext{
		ExcludeNested: true,
		IsVirtual:     true,
		DB: &pgdb_v1.Column{
			Name:      "pkskv2",
			Type:      vcDataType.Name,
			Nullable:  true,
			Collation: "C",
		},
		GoName:   "PKSKV2",
		DataType: vcDataType,
		Convert: &pkskDataConvert{
			ctx: ctx,
		},
		QueryTypeName: ctx.Name(m).String() + "PKSKV2" + "QueryType",
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
			Name:     "fts_data",
			Type:     "tsvector",
			Nullable: true,
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
			Name:     "pb_data",
			Type:     byteaDataType.Name,
			Nullable: false,
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
