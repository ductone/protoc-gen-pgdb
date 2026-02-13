package v1

import (
	"github.com/doug-martin/goqu/v9/exp"

	"github.com/ductone/protoc-gen-pgdb/pgdb/v1/xpq"
)

// FloatToVector to insert a vector type it needs to be of the form '[1.0,2.0,3.0,...]'.
func FloatToVector(in []float32) exp.LiteralExpression {
	return exp.NewLiteralExpression("?::vector", xpq.VectorValue(in))
}

// FloatToHalfVector inserts a halfvec (float16) type using the same '[1.0,2.0,3.0,...]' format.
func FloatToHalfVector(in []float32) exp.LiteralExpression {
	return exp.NewLiteralExpression("?::halfvec", xpq.VectorValue(in))
}
