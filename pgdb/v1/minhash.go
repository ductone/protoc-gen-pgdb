package v1

import (
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/ductone/protoc-gen-pgdb/pgdb/v1/xpq"
)

// In order to insert a minhash type it needs to be of the form '(1,0,0,0,...)::bit'.
func MinHashToBitVector(in []byte) exp.LiteralExpression {
	return exp.NewLiteralExpression("?::bit", xpq.MinHashValue(in))
}
