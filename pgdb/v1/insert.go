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

	versionField := desc.VersioningField()
	if _, ok := record[versionField.Name]; !ok {
		return "", nil, errors.New("pgdb_v1: updated_at missing from message; unable to upsert without " + versionField.Name)
	}

	// postponing backfill until post migration
	/* pkskv2Field := desc.PKSKV2Field()
	if _, ok := record[pkskv2Field.Name]; !ok {
		return "", nil, errors.New("pgdb_v1: pkskv2 missing from message; unable to upsert without " + pkskv2Field.Name)
	}

	pk := record["pb$pk"]
	sk := record["pb$sk"]
	record[pkskv2Field.Name] = fmt.Sprintf("%s|%s", pk.(string), sk.(string))
	*/

	qb := goqu.Dialect("postgres")
	q := qb.Insert(tableName).Prepared(true).Rows(
		record,
	)
	conflictRecords := exp.Record{}
	for k := range record {
		switch k {
		case "pb$pksk":
			continue
		case "pb$tenant_id":
			continue
		case "pb$pkskv2":
			continue
		default:
			conflictRecords[k] = exp.NewIdentifierExpression("", "excluded", k)
		}
	}

	primaryIndex := desc.IndexPrimaryKey()
	if primaryIndex == nil {
		return "", nil, errors.New("pgdb_v1.Insert: malformed message: primary index missing")
	}

	q = q.OnConflict(
		exp.NewDoUpdateConflictExpression(`ON CONSTRAINT "`+primaryIndex.Name+`"`, conflictRecords).Where(
			exp.NewIdentifierExpression("", tableName, versionField.Name).Lte(
				exp.NewLiteralExpression("?::timestamptz", record[versionField.Name]),
			),
		),
	)

	return q.ToSQL()
}
