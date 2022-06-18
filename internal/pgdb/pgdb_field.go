package pgdb

import (
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9/exp"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	pgs "github.com/lyft/protoc-gen-star"
)

type fieldContext struct {
	F       pgs.Field
	DB      pgdb_v1.Field
	Convert *fieldConvert
}

type fieldConvert struct {
	PostgesName string
	// https://www.postgresql.org/docs/current/datatype-numeric.html
	PostgresTypeName string
	IsArray          bool
	TypeConversion   goTypeConversion
	FullTextType     pgdb_v1.FieldOptions_FullTextType
	FullTextWeight   pgdb_v1.FieldOptions_FullTextWeight
}

type goTypeConversion int64

const (
	GT_UNSPECIFIED      goTypeConversion = 0
	GT_FLOAT32          goTypeConversion = 1
	GT_FLOAT64          goTypeConversion = 2
	GT_INT32            goTypeConversion = 3
	GT_INT64            goTypeConversion = 4
	GT_UINT32           goTypeConversion = 5
	GT_UINT64           goTypeConversion = 6
	GT_BOOL             goTypeConversion = 7
	GT_STRING           goTypeConversion = 8
	GT_BYTES            goTypeConversion = 9
	GT_ENUM             goTypeConversion = 10
	GT_PB_WKT_ANY       goTypeConversion = 11
	GT_PB_WKT_TIMESTAMP goTypeConversion = 12
	GT_PB_WKT_DURATION  goTypeConversion = 13
	GT_PB_WKT_STRUCT    goTypeConversion = 14
)

func getField(f pgs.Field) *fieldContext {
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

	convertDef := &fieldConvert{}
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
		convertDef.PostgresTypeName = "boolean"
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
		case "google.protobuf.Any":
			if isArray {
				panic(fmt.Errorf("pgdb: unsupported field type: %v: %s: repeated Any not supported", pt, f.FullyQualifiedName()))
			}
			convertDef.PostgresTypeName = "jsonb"
			convertDef.TypeConversion = GT_PB_WKT_ANY
		case "google.protobuf.Timestamp":
			convertDef.PostgresTypeName = "timestamptz"
			convertDef.IsArray = isArray
			convertDef.TypeConversion = GT_PB_WKT_TIMESTAMP
		case "google.protobuf.Duration":
			convertDef.PostgresTypeName = "interval"
			convertDef.IsArray = isArray
			convertDef.TypeConversion = GT_PB_WKT_DURATION
		case "google.protobuf.Struct":
			if isArray {
				panic(fmt.Errorf("pgdb: unsupported field type: %v: %s: repeated Struct not supported", pt, f.FullyQualifiedName()))
			}
			convertDef.PostgresTypeName = "jsonb"
			convertDef.TypeConversion = GT_PB_WKT_STRUCT
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
		panic(fmt.Errorf("pgdb: unsupported field type: %v: %s (of type %s)", pt, f.FullyQualifiedName(), f.Descriptor().GetType()))
	}
	rv := &fieldContext{
		F:       f,
		Convert: convertDef,
	}
	return rv
}

func (fc *fieldContext) ColumnName() string {
	return string(fc.F.Name())
}

func (fc *fieldContext) ColumnValueExp() string {
	return "m.self." + string(fc.F.Name().UpperCamelCase())
}

func SearchFieldsToQuery(fields []*fieldContext) (exp.LiteralExpression, error) {
	if countSearchFields(fields) == 0 {
		return exp.L("NULL"), nil
	}
	edgeNGramTemplate := "edge_gram_tsvector(?::text, ?)"
	vectorTemplate := "setweight(to_tsvector(?, ?::text), ?)"

	vectors := []string{}
	args := []interface{}{}
	for _, field := range fields {
		if field.Convert.FullTextType == pgdb_v1.FieldOptions_FULL_TEXT_TYPE_UNSPECIFIED {
			continue
		}

		switch field.Convert.FullTextType {
		case pgdb_v1.FieldOptions_FULL_TEXT_TYPE_EDGE_NGRAM:
			vectors = append(vectors, edgeNGramTemplate)
			args = append(args, field.textContents(), field.weight())
			fallthrough
		case pgdb_v1.FieldOptions_FULL_TEXT_TYPE_ENGLISH:
			vectors = append(vectors, vectorTemplate)
			args = append(args, "english", field.textContents(), field.weight())
			fallthrough
		case pgdb_v1.FieldOptions_FULL_TEXT_TYPE_SIMPLE:
			vectors = append(vectors, vectorTemplate)
			args = append(args, "simple", field.textContents(), field.weight())
		}
	}

	query := strings.Join(vectors, " || ")
	return goqu.L(query, args...), nil
}
