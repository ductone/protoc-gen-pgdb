package pgdb

import (
	"fmt"

	pgs "github.com/lyft/protoc-gen-star/v2"
)

func fieldByName(msg pgs.Message, name string) pgs.Field {
	f, ok := tryFieldByName(msg, name)
	if !ok {
		panic(fmt.Sprintf("Failed to find field %s on %s", name, msg.FullyQualifiedName()))
	}
	return f
}

func tryFieldByName(msg pgs.Message, name string) (pgs.Field, bool) {
	for _, f := range msg.Fields() {
		if f.Name().LowerSnakeCase().String() == name {
			return f, true
		}
	}
	return nil, false
}

func getVersioningField(msg pgs.Message) (string, error) {
	possibleFields := []string{
		"updated_at",
		"created_at",
	}
	for _, fn := range possibleFields {
		_, ok := tryFieldByName(msg, fn)
		if ok {
			return fn, nil
		}
	}
	return "", fmt.Errorf("pgdb: getVersioningField: must have one of %v from '%s'", possibleFields, msg.FullyQualifiedName())
}

// Checks if the field is a valid vector field shape and if so returns the field data enum, float array, and vector size.
// Otherwise panics.
func GetFieldVectorShape(field pgs.Field) (pgs.Field, pgs.Field, error) {
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
			floatField = subField
		}
	}

	return enumField, floatField, nil
}

func getFieldMinHashShape(field pgs.Field) pgs.Field {
	subMsg := field.Type().Embed()
	if subMsg == nil {
		panic(fmt.Errorf("pgdb: minhash behavior only supported on message fields: %s", field.FullyQualifiedName()))
	}
	allFields := subMsg.Fields()
	if len(allFields) != 1 {
		panic(fmt.Errorf("pgdb: minhash message must only have byte array field: %s", field.FullyQualifiedName()))
	}

	byteArrayField := allFields[0]
	if byteArrayField.Type().ProtoType() != pgs.BytesT {
		panic(fmt.Errorf("pgdb: minhash message must have byte array field: %s", field.FullyQualifiedName()))
	}
	return byteArrayField
}
