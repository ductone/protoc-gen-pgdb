package xpq

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"strconv"
	"strings"
)

type Array[T sqlable] []T

// jsonb would be the most useful addition.
var SupportedArrayGoTypes = map[string]bool{
	"string":  true,
	"int8":    true,
	"int16":   true,
	"int32":   true,
	"int64":   true,
	"uint8":   true,
	"uint16":  true,
	"uint32":  true,
	"uint64":  true,
	"float32": true,
	"float64": true,
	"bool":    true,
	"[]byte":  true,
}

type sqlable interface {
	string | int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64 | bool | []byte
}

func (a Array[T]) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	sb := &strings.Builder{}
	_, _ = sb.WriteString("{")
	for i, item := range a {
		if i > 0 {
			_, _ = sb.WriteString(", ")
		}
		anyTime := any(item)
		switch v := anyTime.(type) {
		case string:
			stringQuotedBytes(sb, []byte(v))
		case int8:
			_, _ = sb.WriteString(strconv.FormatInt(int64(v), 10))
		case int16:
			_, _ = sb.WriteString(strconv.FormatInt(int64(v), 10))
		case int32:
			_, _ = sb.WriteString(strconv.FormatInt(int64(v), 10))
		case int64:
			_, _ = sb.WriteString(strconv.FormatInt(v, 10))
		case uint8:
			_, _ = sb.WriteString(strconv.FormatUint(uint64(v), 10))
		case uint16:
			_, _ = sb.WriteString(strconv.FormatUint(uint64(v), 10))
		case uint32:
			_, _ = sb.WriteString(strconv.FormatUint(uint64(v), 10))
		case uint64:
			_, _ = sb.WriteString(strconv.FormatUint(v, 10))
		case float32:
			_, _ = sb.WriteString(strconv.FormatFloat(float64(v), 'f', -1, 32))
		case float64:
			_, _ = sb.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
		case bool:
			_, _ = sb.WriteString(strconv.FormatBool(v))
		case []byte:
			_, _ = sb.WriteString(`'\x`)
			_, _ = sb.WriteString(hex.EncodeToString(v))
			_, _ = sb.WriteString(`'`)
		}
	}
	_, _ = sb.WriteString("}")
	return sb.String(), nil
}

func stringQuotedBytes(sb *strings.Builder, v []byte) {
	_, _ = sb.WriteString(`"`)
	for {
		i := bytes.IndexAny(v, `"\`)
		if i < 0 {
			_, _ = sb.Write(v)
			break
		}
		if i > 0 {
			_, _ = sb.Write(v[:i])
		}
		_ = sb.WriteByte('\\')
		_ = sb.WriteByte(v[i])
		v = v[i+1:]
	}
	_, _ = sb.WriteString(`"`)
}
