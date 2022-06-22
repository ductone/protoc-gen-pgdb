package pgdb

import (
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	pgs "github.com/lyft/protoc-gen-star"
)

type fieldConvert struct {
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
)

func (fc *fieldConvert) CodeForValue() (string, error) {
	// TODO(pquerna): template this
	return fc.varName + " := m.self." + string(fc.F.Name().UpperCamelCase()), nil
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
	return templateExecToString("proto_format_string.tmpl", c)
}
