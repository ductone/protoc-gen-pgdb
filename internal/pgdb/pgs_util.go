package pgdb

import (
	"fmt"

	pgs "github.com/lyft/protoc-gen-star"
)

func fieldByName(msg pgs.Message, name string) pgs.Field {
	for _, f := range msg.Fields() {
		if f.Name().LowerSnakeCase().String() == name {
			return f
		}
	}
	panic(fmt.Sprintf("Failed to find field %s on %s", name, msg.FullyQualifiedName()))
}
