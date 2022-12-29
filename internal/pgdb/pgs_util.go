package pgdb

import (
	"fmt"

	pgs "github.com/lyft/protoc-gen-star"
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
