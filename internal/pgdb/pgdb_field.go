package pgdb

import (
	"fmt"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	"github.com/jackc/pgtype"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

const (
	pgTypeJSONB = "jsonb"
)

type fieldContext struct {
	// denotes a realized/virtual field that comes from multiple fields. in this case, F is nil.
	IsVirtual bool
	GoName    string
	Field     pgs.Field
	DB        pgdb_v1.Column
	DataType  *pgtype.DataType
	Convert   FiledConverter
}

type FiledConverter interface {
	CodeForValue() (string, error)
	VarForValue() (string, error)
}

func (module *Module) getField(ctx pgsgo.Context, f pgs.Field, vn *varNamer, ix *importTracker, goPrefix string) *fieldContext {
	ext := pgdb_v1.FieldOptions{}
	_, err := f.Extension(pgdb_v1.E_Options, &ext)
	if err != nil {
		panic(fmt.Errorf("pgdb: getField: failed to extract Message extension from '%s': %w", f.FullyQualifiedName(), err))
	}

	isArray := f.Type().ProtoLabel() == pgs.Repeated
	pt := f.Type().ProtoType()

	// TODO(pquerna): nested fields/messages
	pgColName, err := getColumnName(f, nil)
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
	// https://developers.google.com/protocol-buffers/docs/proto3#scalar
	switch pt {
	case pgs.DoubleT:
		// aka float64
		convertDef.PostgresTypeName = "float8"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = gtFloat64
	case pgs.FloatT:
		// aka float32
		convertDef.PostgresTypeName = "float4"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = gtFloat32
	case pgs.Int32T, pgs.SInt32, pgs.SFixed32:
		convertDef.PostgresTypeName = "int4"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = gtInt32
	case pgs.Int64T, pgs.SInt64, pgs.SFixed64:
		convertDef.PostgresTypeName = "int8"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = gtInt64
	case pgs.UInt32T, pgs.Fixed32T:
		convertDef.PostgresTypeName = "int8"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = gtUint32
	case pgs.UInt64T, pgs.Fixed64T:
		// not ideal, but postgres only has signed types for int8.
		convertDef.PostgresTypeName = "numeric"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = gtUint64
	case pgs.BoolT:
		convertDef.PostgresTypeName = "bool"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = gtBool
	case pgs.StringT:
		// TODO(pquerna): annotations for max size
		convertDef.PostgresTypeName = "text"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = gtString
		convertDef.FullTextType = ext.FullTextType
		convertDef.FullTextWeight = ext.FullTextWeight
	case pgs.MessageT:
		// TODO(pquerna): handle nested messages defined locally and in other modules
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
		default:
			switch ext.MessageBehavoir {
			case pgdb_v1.FieldOptions_MESSAGE_BEHAVOIR_OMIT:
				// explict option to just not store this in postgres
				return nil
			case pgdb_v1.FieldOptions_MESSAGE_BEHAVOIR_EXPAND:
				// getMessageFields(ctx)
			case pgdb_v1.FieldOptions_MESSAGE_BEHAVOIR_JSONB:
				convertDef.IsArray = isArray
				convertDef.PostgresTypeName = pgTypeJSONB
				convertDef.TypeConversion = gtPbGenericMsg
			default:
				panic(fmt.Errorf("pgdb: unsupported message field type: %v: %s (of type %s)",
					pt, f.FullyQualifiedName(), f.Descriptor().GetType()))
			}
		}

	case pgs.BytesT:
		// single bytes and repeated bytes we store the same way
		convertDef.PostgresTypeName = "bytea"
		convertDef.TypeConversion = gtBytes
	case pgs.EnumT:
		convertDef.PostgresTypeName = "int4"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = gtEnum
	case pgs.GroupT:
		panic(fmt.Errorf("pgdb: unsupported field type: Group: %s", f.FullyQualifiedName()))
	default:
		panic(fmt.Errorf("pgdb: unsupported field type: %v: %s (of type %s)",
			pt, f.FullyQualifiedName(), f.Descriptor().GetType()))
	}

	dbTypeRef, ok := pgDataTypeForName(convertDef.PostgresTypeName)
	if !ok {
		panic(fmt.Errorf("pgdb: unsupported field type: %v: %s (of type %s): pgDataTypeForName(%s) NOT FOUND",
			pt, f.FullyQualifiedName(), f.Descriptor().GetType(), convertDef.PostgresTypeName))
	}

	rv := &fieldContext{
		IsVirtual: false,
		GoName:    ctx.Name(f).String(),
		Field:     f,
		DB: pgdb_v1.Column{
			Name: pgColName,
			Type: dbTypeRef.Name,
		},
		DataType: dbTypeRef,
		Convert:  convertDef,
	}
	return rv
}

const (
	// fts_data, pb_data, pk, sk.
	lenCommonFields = 4
)

func getCommonFields(ctx pgsgo.Context, m pgs.Message) ([]*fieldContext, error) {
	vn := &varNamer{prefix: "cfv", offset: 0}
	_ = vn
	vcDataType, _ := pgDataTypeForName("varchar")

	byteaDataType, _ := pgDataTypeForName("bytea")
	tenantIdField := &fieldContext{
		IsVirtual: true,
		DB: pgdb_v1.Column{
			Name: "tenant_id",
			Type: vcDataType.Name,
		},
		GoName:   "TenantID",
		DataType: vcDataType,
		Convert: &tenantIdDataConvert{
			ctx:     ctx,
			Message: m,
			VarName: vn.String(),
		},
	}

	vn = vn.Next()
	pkField := &fieldContext{
		IsVirtual: true,
		DB: pgdb_v1.Column{
			Name: "pk",
			Type: vcDataType.Name,
		},
		GoName:   "PK",
		DataType: vcDataType,
		Convert: &dynamoKeyDataConvert{
			ctx:     ctx,
			VarName: vn.String(),
			Message: m,
			KeyType: DynamoKeyTypePartition,
		},
	}

	vn = vn.Next()
	skField := &fieldContext{
		IsVirtual: true,
		DB: pgdb_v1.Column{
			Name: "sk",
			Type: vcDataType.Name,
		},
		GoName:   "SK",
		DataType: vcDataType,
		Convert: &dynamoKeyDataConvert{
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
		IsVirtual: true,
		DB: pgdb_v1.Column{
			Name: "fts_data",
			Type: "tsvector",
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
		IsVirtual: true,
		DB: pgdb_v1.Column{
			Name: "pb_data",
			Type: byteaDataType.Name,
		},
		GoName:   "PBData",
		DataType: byteaDataType,
		Convert: &pbDataConvert{
			VarName: vn.String(),
		},
	}
	return []*fieldContext{tenantIdField, pkField, skField, ftsDataField, pbDataField}, nil
}
