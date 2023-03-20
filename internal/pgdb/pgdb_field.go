package pgdb

import (
	"fmt"
	"strconv"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	"github.com/jackc/pgx/v5/pgtype"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

const (
	pgTypeJSONB = "jsonb"
	pgTypeInt4  = "int4"
	pgTypeBool  = "bool"
)

type fieldContext struct {
	// denotes a realized/virtual field that comes from multiple fields. in this case, F is nil.
	IsVirtual     bool
	ExcludeNested bool
	GoName        string
	Field         pgs.Field
	Nested        bool
	DB            *pgdb_v1.Column
	DataType      *pgtype.Type
	Convert       FiledConverter
}

type FiledConverter interface {
	GoType() (string, error)
	CodeForValue() (string, error)
	VarForValue() (string, error)
	VarForAppend() (string, error)
}

func (module *Module) getField(ctx pgsgo.Context, f pgs.Field, vn *varNamer, ix *importTracker, goPrefix string) *fieldContext {
	ext := pgdb_v1.FieldOptions{}
	_, err := f.Extension(pgdb_v1.E_Options, &ext)
	if err != nil {
		panic(fmt.Errorf("pgdb: getField: failed to extract Message extension from '%s': %w", f.FullyQualifiedName(), err))
	}

	if ext.MessageBehavoir == pgdb_v1.FieldOptions_MESSAGE_BEHAVOIR_OMIT {
		// explict option to just not store this in postgres
		return nil
	}

	isArray := f.Type().ProtoLabel() == pgs.Repeated && !f.Type().IsMap()
	pt := f.Type().ProtoType()

	pgColName, err := getColumnName(f)
	if err != nil {
		panic(fmt.Errorf("pgdb: getColumnName failed for: %v: %s (of type %s)",
			pt, f.FullyQualifiedName(), f.Descriptor().GetType()))
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
		// TODO(pquerna): annotations for max size
		convertDef.PostgresTypeName = "text"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = gtString
		convertDef.FullTextType = ext.FullTextType
		convertDef.FullTextWeight = ext.FullTextWeight
		defaultValue = "''"
		nullable = false
	case pgs.MessageT:
		switch ext.MessageBehavoir {
		case pgdb_v1.FieldOptions_MESSAGE_BEHAVOIR_UNSPECIFIED:
			switch f.Descriptor().GetTypeName() {
			case ".google.protobuf.Any":
				if isArray {
					panic(fmt.Errorf("pgdb: unsupported field type: %v: %s: repeated Any not supported", pt, f.FullyQualifiedName()))
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
					panic(fmt.Errorf("pgdb: unsupported message field type: %v: %s (of type %s): Arrays cannot be nested; consider jsonb",
						pt, f.FullyQualifiedName(), f.Descriptor().GetType()))
				}
				convertDef.TypeConversion = gtPbNestedMsg
				convertDef.NestedPrefix = strconv.FormatInt(int64(*f.Descriptor().Number), 10) + "$"
			}
		case pgdb_v1.FieldOptions_MESSAGE_BEHAVOIR_OMIT:
			// explict option to just not store this in postgres
			return nil
		case pgdb_v1.FieldOptions_MESSAGE_BEHAVOIR_EXPAND:
			if isArray {
				panic(fmt.Errorf("pgdb: unsupported message field type: %v: %s (of type %s): Arrays cannot be nested; consider jsonb",
					pt, f.FullyQualifiedName(), f.Descriptor().GetType()))
			}
			convertDef.TypeConversion = gtPbNestedMsg
			convertDef.NestedPrefix = strconv.FormatInt(int64(*f.Descriptor().Number), 10) + "$"
		case pgdb_v1.FieldOptions_MESSAGE_BEHAVOIR_JSONB:
			convertDef.IsArray = isArray
			convertDef.PostgresTypeName = pgTypeJSONB
			convertDef.TypeConversion = gtPbGenericMsg
		default:
			panic(fmt.Errorf("pgdb: unsupported message field type: %v: %s (of type %s)",
				pt, f.FullyQualifiedName(), f.Descriptor().GetType()))
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
		panic(fmt.Errorf("pgdb: unsupported field type: Group: %s", f.FullyQualifiedName()))
	default:
		panic(fmt.Errorf("pgdb: unsupported field type: %v: %s (of type %s)",
			pt, f.FullyQualifiedName(), f.Descriptor().GetType()))
	}

	if isArray {
		nullable = true
		defaultValue = ""
		ix.XPQ = true
		convertDef.PostgresTypeName = "_" + convertDef.PostgresTypeName
	}

	rv := &fieldContext{
		IsVirtual: false,
		GoName:    ctx.Name(f).String(),
		Field:     f,
		Convert:   convertDef,
	}

	if convertDef.TypeConversion != gtPbNestedMsg {
		dbTypeRef := pgDataTypeForName(convertDef.PostgresTypeName)
		if !nullable && defaultValue == "" {
			panic(fmt.Errorf("pgdb: nullable column with no default: %s (%s)", pgColName, dbTypeRef.Name))
		}
		rv.DB = &pgdb_v1.Column{
			Name:     pgColName,
			Type:     dbTypeRef.Name,
			Nullable: nullable,
			Default:  defaultValue,
		}
		rv.DataType = dbTypeRef
	} else {
		rv.Nested = true
	}
	return rv
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

	// nested only currently don't have any of the common fields.
	if fext.NestedOnly {
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
	}

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
	}

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
	}

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
	}
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
	}

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
	}
	return []*fieldContext{tenantIdField, pkskField, pkField, skField, ftsDataField, pbDataField}, nil
}
