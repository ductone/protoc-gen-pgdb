package v1

import (
	"fmt"

	"github.com/doug-martin/goqu/v9/exp"
	"github.com/ductone/protoc-gen-pgdb/pgdb/v1/xpq"
	pgs "github.com/lyft/protoc-gen-star"
)

// In order to insert a vector type it needs to be of the form '[1.0,2.0,3.0,...]'.
func FloatToVector(in []float32) exp.LiteralExpression {
	return exp.NewLiteralExpression("?::vector", xpq.VectorValue(in))
}

// Checks if the field is a valid vector field shape and if so returns the field data enum, float array, and vector size.
// Otherwise panics.
func GetFieldVectorShape(field pgs.Field) (pgs.Field, pgs.Field, int32, error) {
	if !field.Type().IsRepeated() {
		panic(fmt.Errorf("pgdb: vector behavior only supported on repeated fields: %s", field.FullyQualifiedName()))
	}
	subMsg := field.Type().Element().Embed()
	if subMsg == nil {
		panic(fmt.Errorf("pgdb: vector behavior only supported on message fields: %s", field.FullyQualifiedName()))
	}
	allFields := subMsg.Fields()
	if len(allFields) != 2 {
		panic(fmt.Errorf("pgdb: vector message must only have model enum and float array: %s", field.FullyQualifiedName()))
	}

	var enumField pgs.Field
	var floatField pgs.Field
	var vectorSize int32

	for _, subField := range allFields {
		switch subField.Descriptor().GetNumber() {
		case 1:
			// enum
			if subField.Type().ProtoType() != pgs.EnumT {
				panic(fmt.Errorf("pgdb: vector message must have model enum as first field: %s", field.FullyQualifiedName()))
			}
			enumField = subField
		case 2:
			// repeated float
			if !subField.Type().IsRepeated() || subField.Type().Element().ProtoType() != pgs.FloatT {
				panic(fmt.Errorf("pgdb: vector message must have repeated float as second field: %s", field.FullyQualifiedName()))
			}
			subExt := FieldOptions{}
			_, err := subField.Extension(E_Options, &subExt)
			if err != nil {
				return nil, nil, 0, fmt.Errorf("pgdb: getField: failed to extract Message extension from '%s': %w", field.FullyQualifiedName(), err)
			}
			if subExt.VectorSize == 0 {
				panic(fmt.Errorf("pgdb: vector message must have vector_size set on repeated float field: %s", field.FullyQualifiedName()))
			}
			vectorSize = subExt.VectorSize
			floatField = subField
		}
	}

	return enumField, floatField, vectorSize, nil
}
