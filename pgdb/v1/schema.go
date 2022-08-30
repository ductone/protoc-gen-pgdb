package v1

import (
	"bytes"
	"strings"

	"github.com/davecgh/go-spew/spew"
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
		if field.OverrideExpression != "" {
			buf.WriteString(field.OverrideExpression)
		} else {
			buf.WriteString(field.Type)
			if !field.Nullable {
				buf.WriteString(" NOT NULL")
			}
		}
	}
	buf.WriteString("\n")
	for _, idx := range desc.Indexes() {
		if !idx.IsPrimary {
			continue
		}
		buf.WriteString(",\n")

		buf.WriteString("PRIMARY KEY (")
		buf.WriteString(strings.Join(slice.Convert(idx.Columns, func(in string) string {
			return `"` + in + `"`
		}), ","))
		buf.WriteString(")\n")
	}
	buf.WriteString(")\n")
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
	dbr := msg.DBReflect()
	desc := dbr.Descriptor()
	indexes := desc.Indexes()
	rv := make([]string, 0, len(indexes))
	for _, idx := range indexes {
		spew.Dump(idx)
		buf := &bytes.Buffer{}
		if idx.IsPrimary {
			// we only support doing primary indexes in the create table, and don't support changing them, so bye bye.
			continue
		}
		if idx.IsDropped {
			_, _ = buf.WriteString("DROP INDEX")
			// WARNING: unique indexes cannot be dropped
			// concurrently.  Maybe unsafe?
			if !idx.IsUnique {
				buf.WriteString(" CONCURRENTLY")
			}
			buf.WriteString(" IF EXISTS ")
			pgWriteString(buf, idx.Name)
			rv = append(rv, buf.String())
			continue
		}
		_, _ = buf.WriteString("CREATE")
		if idx.IsUnique {
			_, _ = buf.WriteString(" UNIQUE")
		}
		_, _ = buf.WriteString(" INDEX CONCURRENTLY IF NOT EXISTS\n  ")
		pgWriteString(buf, idx.Name)
		_, _ = buf.WriteString("\nON\n  ")
		pgWriteString(buf, desc.TableName())
		_, _ = buf.WriteString("\nUSING\n  ")
		switch idx.Method {
		case MessageOptions_Index_INDEX_METHOD_BTREE:
			_, _ = buf.WriteString("BTREE")
		case MessageOptions_Index_INDEX_METHOD_GIN:
			_, _ = buf.WriteString("GIN")
		case MessageOptions_Index_INDEX_METHOD_BTREE_GIN:
			// btree gin just means we can index
			// col types in a multi-col index that aren't
			// noramlly supporte dy gin, eg, varchar,
			// but its not actually a new index type!
			_, _ = buf.WriteString("GIN")
		}
		_, _ = buf.WriteString("\n(\n")
		if idx.OverrideExpression != "" {
			_, _ = buf.WriteString(idx.OverrideExpression)
		} else {
			buf.WriteString(strings.Join(slice.Convert(idx.Columns, func(in string) string {
				return `  "` + in + `"`
			}), ", \n"))
		}
		_, _ = buf.WriteString("\n)\n")
		if idx.WherePredicate != "" {
			_, _ = buf.WriteString("WHERE ")
			_, _ = buf.WriteString(idx.WherePredicate)
			_, _ = buf.WriteString("\n")
		}
		rv = append(rv, buf.String())
	}
	return rv, nil
}

func pgWriteString(buf *bytes.Buffer, input string) {
	_, _ = buf.WriteString(`"`)
	// TODO(pquerna): not completely correct escaping
	_, _ = buf.WriteString(input)
	_, _ = buf.WriteString(`"`)
}
