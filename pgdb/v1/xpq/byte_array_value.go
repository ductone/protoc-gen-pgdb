package xpq

import (
	"database/sql/driver"
	"encoding/hex"
	"strings"
)

type ByteArrayValue []byte

func (b ByteArrayValue) Value() (driver.Value, error) {
	if len(b) == 0 {
		return "[]", nil
	}

	buf := strings.Builder{}
	buf.Grow(len(b) * 16)
	_, _ = buf.WriteString("E'\\x")
	_, _ = buf.WriteString(hex.EncodeToString(b))
	_, _ = buf.WriteString("'")
	return buf.String(), nil
}
