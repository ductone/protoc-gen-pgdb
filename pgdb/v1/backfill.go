package v1

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

func BackfillPKSKV2(msg DBReflectMessage) (string, []any, error) {
	dbr := msg.DBReflect()
	desc := dbr.Descriptor()
	tableName := desc.TableName()
	record, err := dbr.Record()
	if err != nil {
		return "", nil, err
	}

	pkskField := desc.PKSKField()
	if _, ok := record[pkskField.Name]; !ok {
		return "", nil, fmt.Errorf("pgdb_v1: cannot backfill pb$pkskv2 field: %s not found in record", pkskField.Name)
	}
	pksk := record[pkskField.Name].(string)

	pkskv2Field := desc.PKSKV2Field()

	qb := goqu.Dialect("postgres")
	q := qb.Update(tableName).Prepared(true).Set(exp.Record{
		pkskv2Field.Name: goqu.C(pkskField.Name),
	}).Where(
		goqu.C(pkskField.Name).Eq(
			pksk,
		),
		goqu.C(pkskv2Field.Name).IsNull(),
	)

	return q.ToSQL()
}
