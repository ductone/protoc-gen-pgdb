package v1

import (
	"context"
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
)

// CatalogSnapshot holds pre-fetched catalog metadata for all tables,
// eliminating the need for per-table catalog queries during migration planning.
type CatalogSnapshot struct {
	// columns maps table_name -> set of column names
	columns map[string]map[string]struct{}
	// indexes maps table_name -> set of index names
	indexes map[string]map[string]struct{}
	// stats maps table_name -> set of statistic names
	stats map[string]map[string]struct{}
	// storageParams maps table_name -> storage parameter key-value pairs
	storageParams map[string]map[string]string
}

// ReadCatalogSnapshot fetches columns, indexes, statistics, and storage parameters
// for all tables in the public schema in 4 bulk queries, rather than 4 queries per table.
func ReadCatalogSnapshot(ctx context.Context, db sqlScanner) (*CatalogSnapshot, error) {
	snap := &CatalogSnapshot{
		columns:       make(map[string]map[string]struct{}),
		indexes:       make(map[string]map[string]struct{}),
		stats:         make(map[string]map[string]struct{}),
		storageParams: make(map[string]map[string]string),
	}

	if err := snap.readAllColumns(ctx, db); err != nil {
		return nil, fmt.Errorf("catalog snapshot: reading columns: %w", err)
	}
	if err := snap.readAllIndexes(ctx, db); err != nil {
		return nil, fmt.Errorf("catalog snapshot: reading indexes: %w", err)
	}
	if err := snap.readAllStats(ctx, db); err != nil {
		return nil, fmt.Errorf("catalog snapshot: reading stats: %w", err)
	}
	if err := snap.readAllStorageParams(ctx, db); err != nil {
		return nil, fmt.Errorf("catalog snapshot: reading storage params: %w", err)
	}

	return snap, nil
}

func (s *CatalogSnapshot) readAllColumns(ctx context.Context, db sqlScanner) error {
	// Use pg_attribute + pg_class + pg_namespace instead of information_schema.columns
	// for significantly better performance. information_schema.columns is a complex view
	// that joins many catalog tables; querying pg_attribute directly is much faster.
	query := `SELECT c.relname, a.attname
FROM pg_attribute a
JOIN pg_class c ON c.oid = a.attrelid
JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE n.nspname = 'public'
  AND a.attnum > 0
  AND NOT a.attisdropped
  AND c.relkind IN ('r', 'p')
ORDER BY c.relname`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName, columnName string
		if err := rows.Scan(&tableName, &columnName); err != nil {
			return err
		}
		if s.columns[tableName] == nil {
			s.columns[tableName] = make(map[string]struct{})
		}
		s.columns[tableName][columnName] = struct{}{}
	}
	return rows.Err()
}

func (s *CatalogSnapshot) readAllIndexes(ctx context.Context, db sqlScanner) error {
	// Query pg_index + pg_class directly instead of pg_indexes view.
	query := `SELECT ct.relname, ci.relname
FROM pg_index i
JOIN pg_class ct ON ct.oid = i.indrelid
JOIN pg_class ci ON ci.oid = i.indexrelid
JOIN pg_namespace n ON n.oid = ct.relnamespace
WHERE n.nspname = 'public'
ORDER BY ct.relname`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName, indexName string
		if err := rows.Scan(&tableName, &indexName); err != nil {
			return err
		}
		if s.indexes[tableName] == nil {
			s.indexes[tableName] = make(map[string]struct{})
		}
		s.indexes[tableName][indexName] = struct{}{}
	}
	return rows.Err()
}

func (s *CatalogSnapshot) readAllStats(ctx context.Context, db sqlScanner) error {
	dialect := goqu.Dialect("postgres")

	qb := dialect.From("pg_statistic_ext")
	qb = qb.Select("pg_class.relname", "pg_statistic_ext.stxname")
	qb = qb.Join(goqu.T("pg_class"), goqu.On(goqu.I("pg_class.oid").Eq(goqu.I("pg_statistic_ext.stxrelid"))))
	qb = qb.Join(goqu.T("pg_namespace"), goqu.On(goqu.I("pg_namespace.oid").Eq(goqu.I("pg_class.relnamespace"))))
	qb = qb.Where(goqu.L("pg_namespace.nspname = ?", "public"))
	query, params, err := qb.ToSQL()
	if err != nil {
		return err
	}

	rows, err := db.Query(ctx, query, params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName, stName string
		if err := rows.Scan(&tableName, &stName); err != nil {
			return err
		}
		if s.stats[tableName] == nil {
			s.stats[tableName] = make(map[string]struct{})
		}
		s.stats[tableName][stName] = struct{}{}
	}
	return rows.Err()
}

func (s *CatalogSnapshot) readAllStorageParams(ctx context.Context, db sqlScanner) error {
	dialect := goqu.Dialect("postgres")

	qb := dialect.From("pg_class")
	qb = qb.Select("pg_class.relname", "pg_class.reloptions")
	qb = qb.Join(goqu.T("pg_namespace"), goqu.On(goqu.I("pg_namespace.oid").Eq(goqu.I("pg_class.relnamespace"))))
	qb = qb.Where(goqu.L("pg_namespace.nspname = ?", "public"))
	qb = qb.Where(goqu.L("pg_class.reloptions IS NOT NULL"))
	query, params, err := qb.ToSQL()
	if err != nil {
		return err
	}

	rows, err := db.Query(ctx, query, params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		var reloptions []string
		if err := rows.Scan(&tableName, &reloptions); err != nil {
			return err
		}
		if reloptions == nil {
			continue
		}
		params := make(map[string]string)
		for _, opt := range reloptions {
			parts := strings.SplitN(opt, "=", 2)
			if len(parts) == 2 {
				params[parts[0]] = parts[1]
			}
		}
		if len(params) > 0 {
			s.storageParams[tableName] = params
		}
	}
	return rows.Err()
}

// columnsForTable returns the columns for a table, or nil if the table doesn't exist.
func (s *CatalogSnapshot) columnsForTable(tableName string) map[string]struct{} {
	return s.columns[tableName]
}

// indexesForTable returns the indexes for a table, or an empty map if none exist.
func (s *CatalogSnapshot) indexesForTable(tableName string) map[string]struct{} {
	if m := s.indexes[tableName]; m != nil {
		return m
	}
	return make(map[string]struct{})
}

// statsForTable returns the statistics for a table, or an empty map if none exist.
func (s *CatalogSnapshot) statsForTable(tableName string) map[string]struct{} {
	if m := s.stats[tableName]; m != nil {
		return m
	}
	return make(map[string]struct{})
}

// storageParamsForTable returns storage parameters for a table, or an empty map if none exist.
func (s *CatalogSnapshot) storageParamsForTable(tableName string) map[string]string {
	if m := s.storageParams[tableName]; m != nil {
		return m
	}
	return make(map[string]string)
}

// MigrationsWithCatalog computes the DDL migrations needed for a single table using
// a pre-fetched CatalogSnapshot instead of querying the catalog per table.
// This is functionally identical to Migrations but avoids per-table catalog queries.
func MigrationsWithCatalog(snap *CatalogSnapshot, msg DBReflectMessage, dialect Dialect) ([]string, error) {
	rv := make([]string, 0)
	dbr := msg.DBReflect(dialect)
	desc := dbr.Descriptor()

	haveCols := snap.columnsForTable(desc.TableName())

	if len(haveCols) == 0 {
		return CreateSchema(msg, dialect)
	}

	for _, field := range desc.Fields() {
		if _, ok := haveCols[field.Name]; ok {
			continue
		}
		query := col2alter(desc, field)
		rv = append(rv, query)
	}

	indexes := snap.indexesForTable(desc.TableName())

	for _, idx := range desc.Indexes() {
		if idx.IsPrimary {
			continue
		}

		_, exists := indexes[idx.Name]
		query := index2sql(desc, idx)

		if idx.IsDropped {
			if exists {
				rv = append(rv, query)
			}
			continue
		}

		if !exists {
			rv = append(rv, query)
			continue
		}
	}

	existingStats := snap.statsForTable(desc.TableName())

	for _, st := range desc.Statistics() {
		_, exists := existingStats[st.Name]
		query := statistics2sql(desc, st)

		if st.IsDropped {
			if exists {
				rv = append(rv, query)
			}
			continue
		}

		if !exists {
			rv = append(rv, query)
			continue
		}
	}

	existingStorageParams := snap.storageParamsForTable(desc.TableName())
	if storageParamsAlter := storageParams2alter(desc, existingStorageParams); storageParamsAlter != "" {
		rv = append(rv, storageParamsAlter)
	}

	return rv, nil
}
