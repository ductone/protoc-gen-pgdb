package v1

import (
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

func Insert(msg DBReflectMessage) (string, []any, error) {
	dbr := msg.DBReflect()
	desc := dbr.Descriptor()
	tableName := desc.TableName()

	record, err := dbr.Record()
	if err != nil {
		return "", nil, err
	}

	qb := goqu.Dialect("postgres")
	q := qb.Insert(tableName).Prepared(true).Rows(
		record,
	)
	conflictRecords := exp.Record{}
	for k, _ := range record {
		switch k {
		case "pb$pksk":
			continue
		case "pb$tenant_id":
			continue
		default:
			conflictRecords[k] = exp.NewIdentifierExpression("", "excluded", k)
		}
	}

	primaryIndex := desc.IndexPrimaryKey()
	if primaryIndex == nil {
		return "", nil, errors.New("pgdb_v1.Insert: malformed message: primary index missing")
	}

	q = q.OnConflict(goqu.DoUpdate(`ON CONSTRAINT "`+primaryIndex.Name+`"`,
		conflictRecords,
	).Where(
		exp.NewIdentifierExpression("", "excluded", "pb$updated_at").Gte(
			exp.NewLiteralExpression("?::timestamptz", record["pb$updated_at"]),
		),
	))

	return q.ToSQL()
}
