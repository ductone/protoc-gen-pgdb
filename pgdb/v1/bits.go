package v1

import (
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/ductone/protoc-gen-pgdb/pgdb/v1/xpq"
)

// In order to insert a bits type it needs to be of the form B'10101010::bit'.
func BytesToBitVector(in []byte) exp.LiteralExpression {
	return exp.NewLiteralExpression("?", xpq.BitsValue(in))
}
