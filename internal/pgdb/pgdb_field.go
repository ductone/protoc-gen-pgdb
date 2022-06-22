package pgdb

import (
	"fmt"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	pgs "github.com/lyft/protoc-gen-star"
)

type fieldContext struct {
	// denotes a realized/virtual field that comes from multiple fields. in this case, F is nil.
	IsVirtual bool
	Field     pgs.Field
	DB        pgdb_v1.Field
	Convert   FiledConverter
}

type FiledConverter interface {
	CodeForValue() (string, error)
	VarForValue() (string, error)
}

func getField(f pgs.Field, vn *varNamer) *fieldContext {
	ext := pgdb_v1.FieldOptions{}
	_, err := f.Extension(pgdb_v1.E_Options, &ext)
	if err != nil {
		panic(fmt.Errorf("pgdb: getField: failed to extract Message extension from '%s': %w", f.FullyQualifiedName(), err))
	}
	if ext.FullTextWeight == pgdb_v1.FieldOptions_FULL_TEXT_WEIGHT_UNSPECIFIED {
		ext.FullTextWeight = pgdb_v1.FieldOptions_FULL_TEXT_WEIGHT_MED
	}

	isArray := f.Type().ProtoLabel() == pgs.Repeated
	pt := f.Type().ProtoType()

	pgColName, err := getColumnName(f)
	if err != nil {
		panic(fmt.Errorf("pgdb: getColumnName failed for: %v: %s (of type %s)",
			pt, f.FullyQualifiedName(), f.Descriptor().GetType()))
	}

	convertDef := &fieldConvert{
		F:       f,
		varName: vn.String(),
	}
	// https://developers.google.com/protocol-buffers/docs/proto3#scalar
	switch pt {
	case pgs.DoubleT:
		// aka float64
		convertDef.PostgresTypeName = "float8"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = GT_FLOAT64
	case pgs.FloatT:
		// aka float32
		convertDef.PostgresTypeName = "float4"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = GT_FLOAT32
	case pgs.Int32T, pgs.SInt32, pgs.SFixed32:
		convertDef.PostgresTypeName = "int4"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = GT_INT32
	case pgs.Int64T, pgs.SInt64, pgs.SFixed64:
		convertDef.PostgresTypeName = "int8"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = GT_INT64
	case pgs.UInt32T, pgs.Fixed32T:
		convertDef.PostgresTypeName = "int8"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = GT_UINT32
	case pgs.UInt64T, pgs.Fixed64T:
		// not ideal, but postgres only has signed types for int8.
		convertDef.PostgresTypeName = "numeric"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = GT_UINT64
	case pgs.BoolT:
		convertDef.PostgresTypeName = "bool"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = GT_BOOL
	case pgs.StringT:
		// TODO(pquerna): annotations for max size
		convertDef.PostgresTypeName = "text"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = GT_STRING
		convertDef.FullTextType = ext.FullTextType
		convertDef.FullTextWeight = ext.FullTextWeight
	case pgs.MessageT:
		// TODO(pquerna): handle nested messages defined locally and in other modules
		switch f.Descriptor().GetTypeName() {
		case ".google.protobuf.Any":
			if isArray {
				panic(fmt.Errorf("pgdb: unsupported field type: %v: %s: repeated Any not supported", pt, f.FullyQualifiedName()))
			}
			convertDef.PostgresTypeName = "jsonb"
			convertDef.TypeConversion = GT_PB_WKT_ANY
		case ".google.protobuf.Timestamp":
			convertDef.PostgresTypeName = "timestamptz"
			convertDef.IsArray = isArray
			convertDef.TypeConversion = GT_PB_WKT_TIMESTAMP
		case ".google.protobuf.Duration":
			convertDef.PostgresTypeName = "interval"
			convertDef.IsArray = isArray
			convertDef.TypeConversion = GT_PB_WKT_DURATION
		case ".google.protobuf.Struct":
			if isArray {
				panic(fmt.Errorf("pgdb: unsupported field type: %v: %s: repeated Struct not supported", pt, f.FullyQualifiedName()))
			}
			convertDef.PostgresTypeName = "jsonb"
			convertDef.TypeConversion = GT_PB_WKT_STRUCT
		default:
			panic(fmt.Errorf("pgdb: unsupported message field type: %v: %s (of type %s)",
				pt, f.FullyQualifiedName(), f.Descriptor().GetType()))
		}

	case pgs.BytesT:
		// single bytes and repeated bytes we store the same way
		convertDef.PostgresTypeName = "bytea"
		convertDef.TypeConversion = GT_BYTES
	case pgs.EnumT:
		convertDef.PostgresTypeName = "int4"
		convertDef.IsArray = isArray
		convertDef.TypeConversion = GT_ENUM
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
		Field:     f,
		DB: pgdb_v1.Field{
			Name:     pgColName,
			DataType: dbTypeRef,
		},
		Convert: convertDef,
	}
	return rv
}

const (
	// fts_data, pb_data, pk, sk
	lenCommonFields = 4
)

func getCommonFields(m pgs.Message) ([]*fieldContext, error) {
	vn := &varNamer{prefix: "cfv", offset: 0}
	_ = vn
	vcDataType, _ := pgDataTypeForName("varchar")
	tsDataType, _ := pgDataTypeForName("tsvector")
	byteaDataType, _ := pgDataTypeForName("bytea")
	pkField := &fieldContext{
		IsVirtual: true,
		DB: pgdb_v1.Field{
			Name:     "pk",
			DataType: vcDataType,
		},
		Convert: &dynamoKeyDataConvert{
			VarName: vn.String(),
			Message: m,
			KeyType: DKT_PK,
		},
	}
	vn = vn.Next()
	skField := &fieldContext{
		IsVirtual: true,
		DB: pgdb_v1.Field{
			Name:     "sk",
			DataType: vcDataType,
		},
		Convert: &dynamoKeyDataConvert{
			VarName: vn.String(),
			Message: m,
			KeyType: DKT_SK,
		},
	}
	vn = vn.Next()
	ftsDataField := &fieldContext{
		IsVirtual: true,
		DB: pgdb_v1.Field{
			Name:     "fts_data",
			DataType: tsDataType,
		},
		Convert: &ftsDataConvert{
			VarName: vn.String(),
		},
	}
	vn = vn.Next()
	pbDataField := &fieldContext{
		IsVirtual: true,
		DB: pgdb_v1.Field{
			Name:     "pb_data",
			DataType: byteaDataType,
		},
		Convert: &pbDataConvert{
			VarName: vn.String(),
		},
	}
	return []*fieldContext{pkField, skField, ftsDataField, pbDataField}, nil
}
