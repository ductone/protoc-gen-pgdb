package v1

import (
	"bytes"
	"strings"

	"github.com/ductone/protoc-gen-pgdb/internal/slice"
)

func CreateSchema(msg DBReflectMessage) ([]string, error) {
	dbr := msg.DBReflect()
	desc := dbr.Descriptor()
	buf := &bytes.Buffer{}
	_, _ = buf.WriteString("CREATE TABLE ")
	pgWriteString(buf, desc.TableName())
	_, _ = buf.WriteString(" (\n")
	first := true
	for _, field := range desc.Fields() {
		if first {
			first = false
		} else {
			buf.WriteString(",\n")
		}
		pgWriteString(buf, field.Name)
		buf.WriteString(" ")
		buf.WriteString(field.Type)
		if !field.Nullable {
			buf.WriteString(" NOT NULL")
		}
	}
	buf.WriteString("\n")
	for _, idx := range desc.Indexes() {
		if !idx.IsPrimary {
			continue
		}

		buf.WriteString("PRIMARY KEY (")
		buf.WriteString(strings.Join(slice.Convert(idx.Columns, func(in string) string {
			return `"` + in + `"`
		}), ","))
		buf.WriteString(")\n")
	}
	rv := []string{buf.String()}
	buf.Reset()
	more, err := IndexSchema(msg)
	if err != nil {
		return nil, err
	}
	rv = append(rv, more...)
	return rv, nil
}

func IndexSchema(msg DBReflectMessage) ([]string, error) {
	return nil, nil
}

func pgWriteString(buf *bytes.Buffer, input string) {
	_, _ = buf.WriteString(`"`)
	// TODO(pquerna): not completely correct escaping
	_, _ = buf.WriteString(input)
	_, _ = buf.WriteString(`"`)
}
