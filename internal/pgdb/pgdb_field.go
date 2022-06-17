package pgdb

import pgs "github.com/lyft/protoc-gen-star"

type fieldContext struct {
	F pgs.Field
}

func getField(f pgs.Field) *fieldContext {
	return &fieldContext{
		F: f,
	}
}

func (fc *fieldContext) ColumnName() string {
	return string(fc.F.Name())
}

func (fc *fieldContext) ColumnValueExp() string {
	return "m.self." + string(fc.F.Name().UpperCamelCase())
}
