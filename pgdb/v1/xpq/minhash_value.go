package xpq

import (
	"database/sql/driver"
	"strings"
)

type MinHashValue []byte

func (b MinHashValue) Value() (driver.Value, error) {
	if len(b) == 0 {
		return "()", nil
	}

	buf := strings.Builder{}
	buf.Grow(len(b) * 16)
	_, _ = buf.WriteString("(")
	needsComma := false
	for _, v := range b {
		if needsComma {
			_, _ = buf.WriteString(",")
		} else {
			needsComma = true
		}
		_, _ = buf.WriteString(ByteToBitsString(v))
	}
	_, _ = buf.WriteString(")")
	return buf.String(), nil
}

func ByteToBitsString(b byte) string {
	var bits [8]string
	for i := 0; i < 8; i++ {
		if (b & (1 << i)) != 0 {
			bits[7-i] = "1"
		} else {
			bits[7-i] = "0"
		}
	}
	return strings.Join(bits[:], ",")
}
