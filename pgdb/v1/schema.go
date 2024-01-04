package v1

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/ductone/protoc-gen-pgdb/internal/slice"
	"github.com/jackc/pgx/v5"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

func CreateSchema(msg DBReflectMessage) ([]string, error) {
	dbr := msg.DBReflect()
	desc := dbr.Descriptor()
	buf := &bytes.Buffer{}
	_, _ = buf.WriteString("CREATE TABLE IF NOT EXISTS\n  ")
	pgWriteString(buf, desc.TableName())
	_, _ = buf.WriteString("\n(\n")

	_, _ = buf.WriteString(
		strings.Join(
			slice.Convert(desc.Fields(), col2spec),
			",\n",
		),
	)

	if idx := desc.IndexPrimaryKey(); idx != nil {
		_, _ = buf.WriteString(",\n  ")
		_, _ = buf.WriteString("CONSTRAINT ")
		_, _ = buf.WriteString(idx.Name)
		_, _ = buf.WriteString(" PRIMARY KEY (")
		_, _ = buf.WriteString(strings.Join(slice.Convert(idx.Columns, func(in string) string {
			return `"` + in + `"`
		}), ","))
		_, _ = buf.WriteString(")\n")
	}

	_, _ = buf.WriteString(")\n")

	if desc.IsPartitioned() {
		_, _ = buf.WriteString("PARTITION BY LIST(")
		_, _ = buf.WriteString(desc.TenantField().Name)
		_, _ = buf.WriteString(")\n")
	}

	rv := []string{buf.String()}

	more, err := IndexSchema(msg)
	if err != nil {
		return nil, err
	}
	rv = append(rv, more...)

	// for _, r := range rv {
	// 	fmt.Fprintf(os.Stderr, " %s\n", r)
	// }
	return rv, nil
}

func IndexSchema(msg DBReflectMessage) ([]string, error) {
	dbr := msg.DBReflect()
	desc := dbr.Descriptor()
	indexes := desc.Indexes()
	rv := make([]string, 0, len(indexes))
	for _, idx := range indexes {
		if idx.IsPrimary {
			// we only support doing primary indexes in the create table, and don't support changing them, so bye bye.
			continue
		}
		if idx.IsDropped {
			// don't add dropped indexes to new tables
			continue
		}
		rv = append(rv, index2sql(desc, idx))
	}
	return rv, nil
}

type sqlScanner interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func readColumns(ctx context.Context, db sqlScanner, desc Descriptor) (map[string]struct{}, error) {
	dialect := goqu.Dialect("postgres")

	qb := dialect.From("information_schema.columns")
	qb = qb.Select("column_name")
	qb = qb.Where(goqu.L("table_name = ?", desc.TableName()))
	query, params, err := qb.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	haveCols := make(map[string]struct{})
	for rows.Next() {
		var columnName string
		err = rows.Scan(&columnName)
		if err != nil {
			return nil, err
		}
		haveCols[columnName] = struct{}{}
	}
	return haveCols, nil
}

func readIndexes(ctx context.Context, db sqlScanner, desc Descriptor) (map[string]struct{}, error) {
	dialect := goqu.Dialect("postgres")

	qb := dialect.From("pg_indexes")
	qb = qb.Select("indexname")
	qb = qb.Where(goqu.L("tablename = ?", desc.TableName()))
	query, params, err := qb.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	indexes := make(map[string]struct{})
	for rows.Next() {
		var indexName string
		err = rows.Scan(&indexName)
		if err != nil {
			return nil, err
		}
		indexes[indexName] = struct{}{}
	}
	return indexes, nil
}

func tableIsParentPartition(ctx context.Context, db sqlScanner, tableName string) (bool, error) {
	dialect := goqu.Dialect("postgres")

	qb := dialect.From("pg_class")
	qb = qb.Select("relname")
	qb = qb.Join(goqu.T("pg_partitioned_table"), goqu.On(goqu.I("pg_partitioned_table.partrelid").Eq(goqu.I("pg_class.oid"))))
	qb = qb.Where(goqu.L("relname = ?", tableName))
	query, params, err := qb.ToSQL()
	if err != nil {
		return false, err
	}

	rows, err := db.Query(ctx, query, params...)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	return false, nil
}

func Migrations(ctx context.Context, db sqlScanner, msg DBReflectMessage) ([]string, error) {
	rv := make([]string, 0)
	dbr := msg.DBReflect()
	desc := dbr.Descriptor()

	haveCols, err := readColumns(ctx, db, desc)
	if err != nil {
		return nil, err
	}

	if len(haveCols) == 0 {
		return CreateSchema(msg)
	}

	for _, field := range desc.Fields() {
		if _, ok := haveCols[field.Name]; ok {
			continue
		}
		query := col2alter(desc, field)
		rv = append(rv, query)
	}

	indexes, err := readIndexes(ctx, db, desc)
	if err != nil {
		return nil, err
	}

	for _, idx := range desc.Indexes() {
		if idx.IsPrimary {
			continue
		}

		_, exists := indexes[idx.Name]
		query := index2sql(desc, idx)

		if idx.IsDropped {
			// if it should be dropped, and its still here, byeeee
			if exists {
				rv = append(rv, query)
			}
			continue
		}

		// doesn't exist, but should, lets go!
		if !exists {
			rv = append(rv, query)
			continue
		}
	}
	return rv, nil
}

func sha256String(input string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

// Child table names
//
//	pbchild_<parentHash>_<tenantId>
//
// If too long then
//
//	pbchild_<parentHash>_<tenantHash>
func createPartitionTableName(tableName string, tenantId string) string {
	newNameFormat := "pbchild_%s_%s"

	const pgMaxTableNameLen = 63
	tenantHash := sha256String(tenantId)[0:8]

	tableSplit := strings.Split(tableName, "_")
	originalHash := tableSplit[len(tableSplit)-1]

	childTableName := fmt.Sprintf(newNameFormat, originalHash, tenantId)
	// These shouldnt ever be over max length because "pgchild_<8chars>_<27chars>" is 43.
	if len(childTableName) > pgMaxTableNameLen {
		childTableName = fmt.Sprintf(newNameFormat, originalHash, tenantHash)
	}
	return childTableName
}

// This will be passed in in C1.
type TenantIteratorFunc func(ctx context.Context) (string, error)
type SchemaUpdateFunc func(ctx context.Context, schema string, args ...interface{}) error

func TenantPartitionsUpdate(ctx context.Context, db sqlScanner, msg DBReflectMessage, iteratorFunc TenantIteratorFunc, updateFunc SchemaUpdateFunc) error {
	tableName := msg.DBReflect().Descriptor().TableName()

	isParentPartition, err := tableIsParentPartition(ctx, db, tableName)
	if err != nil {
		return err
	}

	// The table exists but is not a parent partition.
	if !isParentPartition {
		return nil
	}

	// We'll only need to attach partitions if the table already exists as a regular table.
	// but this shouldn't happen.
	// As for detaching we'll only need to do that if we want to preserve data in a partitioned table.
	createPartitionSchema := `CREATE TABLE IF NOT EXISTS %s PARTITION OF %s FOR VALUES IN ($1);`

	for {
		tenantId, err := iteratorFunc(ctx)
		if err != nil {
			return err
		}
		if tenantId == "" {
			break
		}
		partitionTableName := createPartitionTableName(tableName, tenantId)
		// fmt.Printf("Creating partition table %s for %s\n", partitionTableName, tenantId)
		builtSchema := fmt.Sprintf(createPartitionSchema, partitionTableName, tableName)
		updateErr := updateFunc(ctx, builtSchema, tenantId)
		if updateErr != nil {
			return updateErr
		}
	}

	return nil
}
