package xpq

import (
	"database/sql/driver"
	"strconv"
	"strings"
)

type VectorValue []float32

func (v VectorValue) Value() (driver.Value, error) {
	if len(v) == 0 {
		return "[]", nil
	}

	buf := strings.Builder{}
	buf.Grow(len(v) * 16)
	_, _ = buf.WriteString("[")
	needsComma := false
	for _, v := range v {
		if needsComma {
			_, _ = buf.WriteString(",")
		} else {
			needsComma = true
		}
		s32 := strconv.FormatFloat(float64(v), 'f', -1, 32)
		_, _ = buf.WriteString(s32)
	}
	_, _ = buf.WriteString("]")
	return buf.String(), nil
}
