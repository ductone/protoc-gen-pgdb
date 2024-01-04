package pgdb

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	"github.com/jackc/pgx/v5"
	pgs "github.com/lyft/protoc-gen-star"
)

func fieldByName(msg pgs.Message, name string) pgs.Field {
	f, ok := tryFieldByName(msg, name)
	if !ok {
		panic(fmt.Sprintf("Failed to find field %s on %s", name, msg.FullyQualifiedName()))
	}
	return f
}

func tryFieldByName(msg pgs.Message, name string) (pgs.Field, bool) {
	for _, f := range msg.Fields() {
		if f.Name().LowerSnakeCase().String() == name {
			return f, true
		}
	}
	return nil, false
}

func getVersioningField(msg pgs.Message) (string, error) {
	possibleFields := []string{
		"updated_at",
		"created_at",
	}
	for _, fn := range possibleFields {
		_, ok := tryFieldByName(msg, fn)
		if ok {
			return fn, nil
		}
	}
	return "", fmt.Errorf("pgdb: getVersioningField: must have one of %v from '%s'", possibleFields, msg.FullyQualifiedName())
}

type sqlScanner interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// Get a list of the provided descriptor's partition sub tables.

func ReadPartitionSubTables(ctx context.Context, db sqlScanner, desc pgdb_v1.Descriptor) ([]string, error) {
	dialect := goqu.Dialect("postgres")

	qb := dialect.From("pg_inherits")
	qb = qb.Select("child.relname").As("child")
	qb = qb.Join(goqu.T("pg_class").As("parent"), goqu.On(goqu.I("pg_inherits.inhparent").Eq(goqu.I("parent.oid"))))
	qb = qb.Join(goqu.T("pg_class").As("child"), goqu.On(goqu.I("pg_inherits.inhrelid").Eq(goqu.I("child.oid"))))
	qb = qb.Join(goqu.T("pg_namespace").As("nmsp_parent"), goqu.On(goqu.I("nmsp_parent.oid").Eq(goqu.I("parent.relnamespace"))))
	qb = qb.Join(goqu.T("pg_namespace").As("nmsp_child"), goqu.On(goqu.I("nmsp_child.oid").Eq(goqu.I("child.relnamespace"))))
	qb = qb.Where(goqu.L("parent.relname = ?", desc.TableName()))
	query, params, err := qb.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tables := make([]string, 0)
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}
