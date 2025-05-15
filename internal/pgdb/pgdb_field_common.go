package pgdb

import (
	"fmt"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
	"google.golang.org/protobuf/types/descriptorpb"
)

type fieldConvert struct {
	// starts as "m.self.", nested messages append
	goPrefix     string
	ctx          pgsgo.Context
	ix           *importTracker
	F            pgs.Field
	varName      string
	NestedPrefix string
	// https://www.postgresql.org/docs/current/datatype-numeric.html
	PostgresTypeName string
	IsArray          bool
	TypeConversion   goTypeConversion
	FullTextType     pgdb_v1.FieldOptions_FullTextType
	FullTextWeight   pgdb_v1.FieldOptions_FullTextWeight
}

type goTypeConversion int64

const (
	//nolint:deadcode,varcheck // i like unsued unspecified
	gtUnspecified      goTypeConversion = 0
	gtFloat32          goTypeConversion = 1
	gtFloat64          goTypeConversion = 2
	gtInt32            goTypeConversion = 3
	gtInt64            goTypeConversion = 4
	gtUint32           goTypeConversion = 5
	gtUint64           goTypeConversion = 6
	gtBool             goTypeConversion = 7
	gtString           goTypeConversion = 8
	gtBytes            goTypeConversion = 9
	gtEnum             goTypeConversion = 10
	gtPbWktAny         goTypeConversion = 11
	gtPbWktTimestamp   goTypeConversion = 12
	gtPbWktDuration    goTypeConversion = 13
	gtPbWktStruct      goTypeConversion = 14
	gtPbGenericMsg     goTypeConversion = 15
	gtPbNestedMsg      goTypeConversion = 16
	gtPbWktBoolValue   goTypeConversion = 17
	gtPbWktStringValue goTypeConversion = 18
	gtInetAddr         goTypeConversion = 19
)

type formatContext struct {
	VarName      string
	InputName    string
	CastType     string
	IsArray      bool
	NestedPrefix string
}

const fieldConvertString = "string"

func (fc *fieldConvert) GoType() (string, error) {
	switch fc.TypeConversion {
	case gtFloat32:
		return "float32", nil
	case gtFloat64:
		return "float64", nil
	case gtInt32:
		return "int32", nil
	case gtInt64:
		return "int64", nil
	case gtUint32:
		return "uint32", nil
	case gtUint64:
		return "uint64", nil
	case gtBool:
		return "bool", nil
	case gtString:
		return fieldConvertString, nil
	case gtBytes:
		return "[]byte", nil
	case gtEnum:
		return "int32", nil
	case gtPbWktAny, gtPbWktStruct, gtPbGenericMsg:
		// objects are stored as JSONB, take input as interface{}, convert to
		return "any", nil
	case gtPbWktTimestamp:
		return "time.Time", nil
	case gtPbWktDuration:
		return "time.Duration", nil
	case gtPbWktBoolValue:
		return "invalid", nil
	case gtPbWktStringValue:
		return "invalid", nil
	case gtPbNestedMsg:
		return string(fc.ctx.Type(fc.F)), nil
	case gtInetAddr:
		return fieldConvertString, nil
	default:
		panic(fmt.Errorf("pgdb: Implement fieldConvert.GoType for %v", fc.TypeConversion))
	}
}

func opaqueFieldGetter(f pgs.Field) string {
	return "Get" + pgsgo.PGGUpperCamelCase(f.Name()).String() + "()"
}

func (fc *fieldConvert) CodeForValue() (string, error) {
	// selfName := fc.goPrefix + ".Get" + fc.ctx.Name(fc.F).String() + "()"
	selfName := fc.goPrefix + "." + opaqueFieldGetter(fc.F)
	switch fc.TypeConversion {
	case gtFloat32:
		return templateExecToString("proto_format_cast.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  "float32",
			IsArray:   fc.IsArray,
		})
	case gtFloat64:
		return templateExecToString("proto_format_cast.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  "float64",
			IsArray:   fc.IsArray,
		})
	case gtInt32:
		return templateExecToString("proto_format_cast.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  "int32",
			IsArray:   fc.IsArray,
		})
	case gtInt64:
		return templateExecToString("proto_format_cast.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  "int64",
			IsArray:   fc.IsArray,
		})
	case gtUint32:
		return templateExecToString("proto_format_cast.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  "uint32",
			IsArray:   fc.IsArray,
		})
	case gtUint64:
		return templateExecToString("proto_format_cast.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  "uint64",
			IsArray:   fc.IsArray,
		})
	case gtBool:
		return templateExecToString("proto_format_cast.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  pgTypeBool,
			IsArray:   fc.IsArray,
		})
	case gtString:
		return templateExecToString("proto_format_cast.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  fieldConvertString,
			IsArray:   fc.IsArray,
		})
	case gtInetAddr:
		fc.ix.NetIP = true
		return templateExecToString("proto_format_inet.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
		})
	case gtBytes:
		return templateExecToString("proto_format_cast.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  "[]byte",
			IsArray:   fc.IsArray,
		})
	case gtEnum:
		return templateExecToString("proto_format_cast.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  "int32",
			IsArray:   fc.IsArray,
		})
	case gtPbWktTimestamp:
		fc.ix.Time = true
		return templateExecToString("proto_format_time.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
		})
	case gtPbWktDuration:
		fc.ix.PgType = true
		return templateExecToString("proto_format_duration.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
		})
	case gtPbWktStringValue:
		return templateExecToString("proto_format_wrapper.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
			CastType:  fieldConvertString,
		})
	case gtPbWktBoolValue:
		return templateExecToString("proto_format_wrapper.tmpl", &formatContext{
			VarName:   fc.varName,
			InputName: selfName,
		})
	case gtPbWktStruct, gtPbWktAny, gtPbGenericMsg:
		useProtoJSON := false

		ft := fc.F.Type()
		pt := ft.ProtoType()

		if pt.Proto() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
			if !ft.IsMap() {
				useProtoJSON = true
			}
		}

		if useProtoJSON {
			if fc.IsArray {
				fc.ix.Bytes = true
			}
			return templateExecToString("proto_format_json_proto.tmpl", &formatContext{
				VarName:   fc.varName,
				InputName: selfName,
				IsArray:   fc.IsArray,
			})
		} else {
			fc.ix.JSON = true
			return templateExecToString("proto_format_json_std.tmpl", &formatContext{
				VarName:   fc.varName,
				InputName: selfName,
				IsArray:   fc.IsArray,
			})
		}
	case gtPbNestedMsg:
		return templateExecToString("proto_format_nested.tmpl", &formatContext{
			VarName:      fc.varName,
			InputName:    selfName,
			IsArray:      fc.IsArray,
			NestedPrefix: fc.NestedPrefix,
		})
	default:
		panic(fmt.Errorf("pgdb: Implement CodeForValue for %v", fc.TypeConversion))
	}
}

func (fc *fieldConvert) VarForValue() (string, error) {
	if fc.TypeConversion == gtPbNestedMsg {
		return "", nil
	}
	return fc.varName, nil
}

func (fc *fieldConvert) VarForAppend() (string, error) {
	if fc.TypeConversion != gtPbNestedMsg {
		return "", nil
	}
	return fc.varName, nil
}

func (fc *fieldConvert) EnumForValue() (string, error) {
	return "", nil
}

type stringFormatContext struct {
	IsFloat  bool
	IsInt    bool
	IsUint   bool
	IsString bool
	VarName  string
}

func typeToString(ix *importTracker, pt pgs.ProtoType, varName string) (string, error) {
	c := stringFormatContext{
		VarName: varName,
	}
	switch pt {
	case pgs.DoubleT, pgs.FloatT:
		ix.Strconv = true
		c.IsFloat = true
	case pgs.Int64T, pgs.SFixed64, pgs.SInt64, pgs.Int32T, pgs.SFixed32, pgs.SInt32, pgs.EnumT:
		ix.Strconv = true
		c.IsInt = true
	case pgs.UInt64T, pgs.Fixed64T, pgs.UInt32T, pgs.Fixed32T:
		ix.Strconv = true
		c.IsUint = true
	case pgs.StringT:
		c.IsString = true
	default:
		panic("typeToString: need to implement for your type")
	}
	return templateExecToString("proto_convert_string.tmpl", c)
}
