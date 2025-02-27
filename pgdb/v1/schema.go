package v1

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/ductone/protoc-gen-pgdb/internal/slice"
	"github.com/jackc/pgx/v5"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/segmentio/ksuid"
)

func CreateSchema(msg DBReflectMessage) ([]string, error) {
	dbr := msg.DBReflect()
	desc := dbr.Descriptor()

	// Validate only one partitioning strategy is used
	partitioningStrategies := 0
	if desc.IsPartitioned() {
		partitioningStrategies++
	}
	if desc.IsPartitionedByCreatedAt() {
		partitioningStrategies++
	}
	if desc.GetPartitionedByKsuidFieldName() != "" {
		partitioningStrategies++
	}
	if partitioningStrategies > 1 {
		return nil, fmt.Errorf("table %s has multiple partitioning strategies defined", desc.TableName())
	}

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

	switch {
	case desc.IsPartitioned():
		_, _ = buf.WriteString("PARTITION BY LIST(")
		_, _ = buf.WriteString(desc.TenantField().Name)
		_, _ = buf.WriteString(")\n")
	case desc.IsPartitionedByCreatedAt():
		_, _ = buf.WriteString("PARTITION BY RANGE(pb$created_at)\n")
	case desc.GetPartitionedByKsuidFieldName() != "":
		_, _ = buf.WriteString(fmt.Sprintf("PARTITION BY RANGE(pb$%s)\n", desc.GetPartitionedByKsuidFieldName()))
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

func readStats(ctx context.Context, db sqlScanner, desc Descriptor) (map[string]struct{}, error) {
	dialect := goqu.Dialect("postgres")

	/*
		SELECT
		  se.stxname AS statistics_name,
		  n.nspname AS schema_name,
		  c.relname AS table_name
		FROM
		  pg_statistic_ext se
		JOIN
		  pg_class c ON c.oid = se.stxrelid
		JOIN
		  pg_namespace n ON n.oid = c.relnamespace
		WHERE
		  c.relname = 'pb_pasta_ingredient_models_food_v1_0565c036'
		  AND n.nspname = 'public';
	*/
	qb := dialect.From("pg_statistic_ext")
	qb = qb.Select("pg_statistic_ext.stxname")
	qb = qb.Join(goqu.T("pg_class"), goqu.On(goqu.I("pg_class.oid").Eq(goqu.I("pg_statistic_ext.stxrelid"))))
	qb = qb.Join(goqu.T("pg_namespace"), goqu.On(goqu.I("pg_namespace.oid").Eq(goqu.I("pg_class.relnamespace"))))
	qb = qb.Where(goqu.L("pg_class.relname = ?", desc.TableName()))
	qb = qb.Where(goqu.L("pg_namespace.nspname = ?", "public"))
	query, params, err := qb.ToSQL()
	if err != nil {
		return nil, err
	}
	// spew.Dump(query, params)

	rows, err := db.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	// CREATE STATISTICS IF NOT EXISTS "pq_test_stat" ON "pb$tenant_id","pb$lifecycle" FROM "pb_tenant_c1_models_innkeeper_v1_d0c77352"
	statNames := make(map[string]struct{})
	for rows.Next() {
		var stName string
		err = rows.Scan(&stName)
		if err != nil {
			return nil, err
		}
		statNames[stName] = struct{}{}
	}
	// spew.Dump(statNames)
	return statNames, nil
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

	existingStats, err := readStats(ctx, db, desc)
	if err != nil {
		return nil, err
	}

	for _, st := range desc.Statistics() {
		_, exists := existingStats[st.Name]
		query := statistics2sql(desc, st)

		if st.IsDropped {
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

// DatePartitionsUpdate creates partitions for date ranges based on the partitioning scheme.
func DatePartitionsUpdate(ctx context.Context, db sqlScanner, msg DBReflectMessage, startDate, endDate time.Time, updateFunc SchemaUpdateFunc) error {
	desc := msg.DBReflect().Descriptor()
	tableName := desc.TableName()

	// Validate that created_at column exists
	columns, err := readColumns(ctx, db, desc)
	if err != nil {
		return err
	}
	if _, hasCreatedAt := columns["pb$created_at"]; !hasCreatedAt {
		return fmt.Errorf("table %s is configured for date partitioning but missing created_at column", tableName)
	}

	isParentPartition, err := tableIsParentPartition(ctx, db, tableName)
	if err != nil {
		return err
	}

	// The table exists but is not a parent partition
	if !isParentPartition {
		return nil
	}

	// Create partition schema template
	createPartitionSchema := `CREATE TABLE IF NOT EXISTS %s PARTITION OF %s FOR VALUES FROM ($1) TO ($2);`

	// Get partition interval
	interval := desc.GetPartitionDateRange()

	// Iterate through the date range and create partitions
	current := startDate
	for current.Before(endDate) {
		var nextDate time.Time
		var partitionTableName string

		// Calculate the next partition boundary based on interval
		switch interval {
		case MessageOptions_PARTITIONED_BY_DATE_RANGE_DAY:
			nextDate = current.AddDate(0, 0, 1)
			partitionTableName = fmt.Sprintf("%s_%s", tableName, current.Format("2006_01_02"))
		case MessageOptions_PARTITIONED_BY_DATE_RANGE_MONTH:
			nextDate = current.AddDate(0, 1, 0)
			partitionTableName = fmt.Sprintf("%s_%s", tableName, current.Format("2006_01"))
		case MessageOptions_PARTITIONED_BY_DATE_RANGE_YEAR:
			nextDate = current.AddDate(1, 0, 0)
			partitionTableName = fmt.Sprintf("%s_%s", tableName, current.Format("2006"))
		default:
			return fmt.Errorf("unsupported partition interval: %v", interval)
		}

		builtSchema := fmt.Sprintf(createPartitionSchema, partitionTableName, tableName)
		err := updateFunc(ctx, builtSchema, current, nextDate)
		if err != nil {
			return err
		}

		current = nextDate
	}

	return nil
}

// EventIDPartitionsUpdate creates partitions for event ID ranges based on their embedded timestamps.
func EventIDPartitionsUpdate(ctx context.Context, db sqlScanner, msg DBReflectMessage, startDate, endDate time.Time, updateFunc SchemaUpdateFunc) error {
	desc := msg.DBReflect().Descriptor()
	tableName := desc.TableName()

	// Validate that event_id column exists
	columns, err := readColumns(ctx, db, desc)
	if err != nil {
		return err
	}
	if _, hasEventID := columns["pb$event_id"]; !hasEventID {
		return fmt.Errorf("table %s is configured for event_id partitioning but missing event_id column", tableName)
	}

	isParentPartition, err := tableIsParentPartition(ctx, db, tableName)
	if err != nil {
		return err
	}

	// The table exists but is not a parent partition
	if !isParentPartition {
		return nil
	}

	// Create partition schema template
	createPartitionSchema := `CREATE TABLE IF NOT EXISTS %s PARTITION OF %s 
		FOR VALUES FROM (
			'%s'  -- Start KSUID for this time range
		) TO (
			'%s'  -- End KSUID for this time range
		);`

	// Get partition interval
	interval := desc.GetPartitionDateRange()

	// Iterate through the date range and create partitions
	current := startDate
	for current.Before(endDate) {
		var nextDate time.Time
		var partitionTableName string

		// Calculate the next partition boundary based on interval
		switch interval {
		case MessageOptions_PARTITIONED_BY_DATE_RANGE_DAY:
			nextDate = current.AddDate(0, 0, 1)
			partitionTableName = fmt.Sprintf("%s_%s", tableName, current.Format("2006_01_02"))
		case MessageOptions_PARTITIONED_BY_DATE_RANGE_MONTH:
			nextDate = current.AddDate(0, 1, 0)
			partitionTableName = fmt.Sprintf("%s_%s", tableName, current.Format("2006_01"))
		case MessageOptions_PARTITIONED_BY_DATE_RANGE_YEAR:
			nextDate = current.AddDate(1, 0, 0)
			partitionTableName = fmt.Sprintf("%s_%s", tableName, current.Format("2006"))
		default:
			return fmt.Errorf("unsupported partition interval: %v", interval)
		}

		// Generate KSUIDs for the partition boundaries
		minParts := []byte{}
		maxParts := []byte{}
		for i := 0; i < 16; i++ {
			minParts = append(minParts, 0)
			maxParts = append(maxParts, 255)
		}
		startKSUID, err := ksuid.FromParts(current.Add(time.Second), minParts)
		if err != nil {
			return fmt.Errorf("failed to generate start KSUID: %w", err)
		}
		endKSUID, err := ksuid.FromParts(nextDate, maxParts)
		if err != nil {
			return fmt.Errorf("failed to generate end KSUID: %w", err)
		}

		builtSchema := fmt.Sprintf(createPartitionSchema,
			partitionTableName,
			tableName,
			startKSUID.String(),
			endKSUID.String())
		// fmt.Println(builtSchema)

		err = updateFunc(ctx, builtSchema)
		if err != nil {
			return err
		}

		current = nextDate
	}

	return nil
}
