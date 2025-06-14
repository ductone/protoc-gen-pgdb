package v1

import (
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

func Delete(msg DBReflectMessage) (string, []any, error) {
	dbr := msg.DBReflect()
	desc := dbr.Descriptor()
	tableName := desc.TableName()

	record, err := dbr.Record()
	if err != nil {
		return "", nil, err
	}

	versionField := desc.VersioningField()
	if _, ok := record[versionField.Name]; !ok {
		return "", nil, errors.New("pgdb_v1: updated_at missing from message; unable to delete without " + versionField.Name)
	}

	primaryIndex := desc.IndexPrimaryKey()
	if primaryIndex == nil {
		return "", nil, errors.New("pgdb_v1: malformed message: primary index missing")
	}

	qb := goqu.Dialect("postgres")
	q := qb.Delete(tableName).Prepared(true).Where(
		exp.NewIdentifierExpression("", tableName, versionField.Name).Lte(
			exp.NewLiteralExpression("?::timestamptz", record[versionField.Name]),
		),
	)

	for _, colName := range primaryIndex.Columns {
		if colName == "pb$pksk" {
			pksk, err := generatedPKSK(record)
			if err != nil {
				return "", nil, err
			}
			q = q.Where(
				exp.NewIdentifierExpression("", tableName, colName).Eq(
					pksk,
				),
			)
		} else {
			colValue, ok := record[colName]
			if !ok {
				return "", nil, errors.New("pgdb_v1: primary key missing from message; unable to delete without " + colName)
			}

			q = q.Where(
				exp.NewIdentifierExpression("", tableName, colName).Eq(
					colValue,
				),
			)
		}
	}

	return q.ToSQL()
}

func generatedPKSK(record exp.Record) (string, error) {
	pkAny, ok := record["pb$pk"]
	if !ok {
		return "", errors.New("pgdb_v1: pb$pk missing from message")
	}
	skAny, ok := record["pb$sk"]
	if !ok {
		return "", errors.New("pgdb_v1: pb$pk missing from message")
	}
	pk, ok := pkAny.(string)
	if !ok {
		return "", errors.New("pgdb_v1: pb$pk wrong type from message")
	}
	sk, ok := skAny.(string)
	if !ok {
		return "", errors.New("pgdb_v1: pb$pk wrong type from message")
	}
	return pk + "|" + sk, nil
}
