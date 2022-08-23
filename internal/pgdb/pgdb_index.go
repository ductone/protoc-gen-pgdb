package pgdb

import (
	"fmt"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type indexContext struct {
	DB pgdb_v1.Index
}

func (module *Module) getMessageIndexes(ctx pgsgo.Context, m pgs.Message, ix *importTracker) []*indexContext {
	ext := pgdb_v1.MessageOptions{}
	_, err := m.Extension(pgdb_v1.E_Msg, &ext)
	if err != nil {
		panic(fmt.Errorf("pgdb: getFieldIndexes: failed to extract Message extension from '%s': %w", m.FullyQualifiedName(), err))
	}

	rv := make([]*indexContext, 0)
	cf, err := getCommonIndexes(ctx, m)
	if err != nil {
		panic(err)
	}
	rv = append(rv, cf...)

	for _, index := range ext.Indexes {
		rv = append(rv, module.extraIndexes(ctx, m, ix, index))
	}
	return rv
}

func (module *Module) extraIndexes(ctx pgsgo.Context, m pgs.Message, ix *importTracker, idx *pgdb_v1.MessageOptions_Index) *indexContext {
	indexName, err := getIndexName(m, idx.Name)
	if err != nil {
		panic(err)
	}
	rv := &indexContext{
		DB: pgdb_v1.Index{
			Name: indexName,
		},
	}
	if idx.Dropped {
		rv.DB.IsDropped = true
		return rv
	}

	rv.DB.Method = idx.Nethod
	for _, fieldName := range idx.Columns {
		f := fieldByName(m, fieldName)
		pgColName, err := getColumnName(f, nil)
		if err != nil {
			panic(err)
		}
		rv.DB.Columns = append(rv.DB.Columns, pgColName)
	}
	return rv
}

func getCommonIndexes(ctx pgsgo.Context, m pgs.Message) ([]*indexContext, error) {
	primaryIndexName, err := getIndexName(m, "pksk")
	if err != nil {
		return nil, err
	}

	primaryIndex := &indexContext{
		DB: pgdb_v1.Index{
			Name:      primaryIndexName,
			IsPrimary: true,
			IsUnique:  true,
			Method:    pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:   []string{"tenant_id", "pk", "sk"},
		},
	}

	ftsIndexName, err := getIndexName(m, "fts_data")
	if err != nil {
		return nil, err
	}
	ftsIndex := &indexContext{
		DB: pgdb_v1.Index{
			Name:    ftsIndexName,
			Method:  pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE_GIN,
			Columns: []string{"tenant_id", "fts_data"},
		},
	}

	return []*indexContext{primaryIndex, ftsIndex}, nil
}