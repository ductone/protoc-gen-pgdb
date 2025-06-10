package v1

import (
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/ductone/protoc-gen-pgdb/pgdb/v1/xpq"
)

// In order to insert a byte array type it needs to be of the form 'E\'\\x010203...\'::bytea'.
func ByteArrayToHex(in []byte) exp.LiteralExpression {
	return exp.NewLiteralExpression("?::bytea", xpq.ByteArrayValue(in))
}
