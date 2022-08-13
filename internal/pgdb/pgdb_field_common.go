package pgdb

import (
	"fmt"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type fieldConvert struct {
	// starts as "m.self.", nested messages append
	goPrefix    string
	ctx         pgsgo.Context
	ix          *importTracker
	F           pgs.Field
	varName     string
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
	GT_PB_GENERIC_MSG   goTypeConversion = 15
)

type formatContext struct {
	VarName   string
	InputName string
	CastType  string
	IsArray   bool
}

func (fc *fieldConvert) CodeForValue() (string, error) {
	selfName := fc.goPrefix + "." + fc.ctx.Name(fc.F).String()
	switch fc.TypeConversion {
	case GT_FLOAT32:
		return templateExecToString("proto_format_cast.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  "float32",
			IsArray:   fc.IsArray,
		})
	case GT_FLOAT64:
		return templateExecToString("proto_format_cast.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  "float64",
			IsArray:   fc.IsArray,
		})
	case GT_INT32:
		return templateExecToString("proto_format_cast.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  "int32",
			IsArray:   fc.IsArray,
		})
	case GT_INT64:
		return templateExecToString("proto_format_cast.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  "int64",
			IsArray:   fc.IsArray,
		})
	case GT_UINT32:
		return templateExecToString("proto_format_cast.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  "uint32",
			IsArray:   fc.IsArray,
		})
	case GT_UINT64:
		return templateExecToString("proto_format_cast.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  "uint64",
			IsArray:   fc.IsArray,
		})
	case GT_BOOL:
		return templateExecToString("proto_format_cast.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  "bool",
			IsArray:   fc.IsArray,
		})
	case GT_STRING:
		return templateExecToString("proto_format_cast.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  "string",
			IsArray:   fc.IsArray,
		})
	case GT_BYTES:
		return fc.varName + " := " + selfName, nil
	case GT_ENUM:
		return templateExecToString("proto_format_cast.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  "int32",
		})
	case GT_PB_WKT_TIMESTAMP:
		fc.ix.Time = true
		return templateExecToString("proto_format_time.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
		})
	case GT_PB_WKT_DURATION:
		return templateExecToString("proto_format_duration.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
		})
	case GT_PB_WKT_STRUCT, GT_PB_WKT_ANY, GT_PB_GENERIC_MSG:
		fc.ix.ProtobufEncodingJson = true
		if fc.IsArray {
			fc.ix.Bytes = true
		}

		return templateExecToString("proto_format_jsonb.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			IsArray:   fc.IsArray,
		})
	default:
		panic(fmt.Errorf("pgdb: Implement CodeForValue for %v", fc.TypeConversion))
	}
}

func (fc *fieldConvert) VarForValue() (string, error) {
	return fc.varName, nil
}

type stringFormatContext struct {
	IsFloat  bool
	IsInt    bool
	IsUint   bool
	IsString bool
	VarName  string
}

func typeToString(pt pgs.ProtoType, varName string) (string, error) {
	c := stringFormatContext{
		VarName: varName,
	}
	switch pt {
	case pgs.DoubleT, pgs.FloatT:
		c.IsFloat = true
	case pgs.Int64T, pgs.SFixed64, pgs.SInt64, pgs.Int32T, pgs.SFixed32, pgs.SInt32, pgs.EnumT:
		c.IsInt = true
	case pgs.UInt64T, pgs.Fixed64T, pgs.UInt32T, pgs.Fixed32T:
		c.IsUint = true
	case pgs.StringT:
		c.IsString = true
	default:
		panic("typeToString: need to implement for your type")
	}
	return templateExecToString("proto_convert_string.tmpl", c)
}
