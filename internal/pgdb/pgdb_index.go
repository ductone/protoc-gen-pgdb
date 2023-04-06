package pgdb

import (
	"fmt"
	"strings"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type indexContext struct {
	DB            pgdb_v1.Index
	ExcludeNested bool
	SourceFields  []string
	RawColumns    []string
	Fields        []*pgs.Field
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

	usedNames := map[string]struct{}{}
	for _, index := range ext.Indexes {
		if _, ok := usedNames[index.Name]; ok {
			panic(fmt.Errorf("pgdb: getFieldIndexes:index name reused  on Message '%s': %s", m.FullyQualifiedName(), index.Name))
		}
		usedNames[index.Name] = struct{}{}
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

	rv.DB.Method = idx.Method
	// tenantIdField, _ := getTenantIDField(m)

	for _, fieldName := range idx.Columns {
		path := strings.Split(fieldName, "🌮")
		message := m
		resolution := ""
		for i, p := range path {
			// TODO: the path could reference a oneOf which is "virtual"
			f := fieldByName(message, p)
			lastP := i == len(path)-1

			if !lastP {
				resolution += getNestedName(f)
				message = f.Type().Embed()
				continue
			}

			// Access to tenant ids is limited by protoc's limited API (can't do extensions for different pgs.Files),
			//   but we don't materialize them for nested models anyway...
			// if i == 0 && p == tenantIdField {
			// 	// don't snake tenant_id
			// 	resolution += p
			// } else {
			pgColName, err := getColumnName(f)

			if err != nil {
				panic(err)
			}
			resolution += pgColName
			// }

			// module.Logf("old source fields: %s vs %s", ctx.Name(f).String(), resolution)
			rv.SourceFields = append(rv.SourceFields, ctx.Name(f).String())
			rv.DB.Columns = append(rv.DB.Columns, resolution)
			rv.Fields = append(rv.Fields, &f)
		}
	}
	rv.RawColumns = idx.Columns
	return rv
}

func getCommonIndexes(ctx pgsgo.Context, m pgs.Message) ([]*indexContext, error) {
	fext := pgdb_v1.MessageOptions{}
	_, err := m.Extension(pgdb_v1.E_Msg, &fext)
	if err != nil {
		panic(err)
	}

	// nested only currently don't have any of the common fields.
	if fext.NestedOnly {
		return nil, nil
	}

	primaryIndexName, err := getIndexName(m, "pksk")
	if err != nil {
		return nil, err
	}
	var tenantIdField pgs.Field
	if fext.TenantIdField != "" {
		tenantIdField = fieldByName(m, fext.TenantIdField)
	} else {
		tenantIdField = fieldByName(m, "tenant_id")
	}

	primaryIndex := &indexContext{
		ExcludeNested: true,
		DB: pgdb_v1.Index{
			Name:      primaryIndexName,
			IsPrimary: true,
			IsUnique:  true,
			Method:    pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:   []string{"tenant_id", "pksk"},
		},
		SourceFields: []string{"TenantId", "PKSK"},
		Fields:       []*pgs.Field{&tenantIdField, nil},
	}

	// So, we learned early in our deployment that having a second unique index
	// doesn't work well with Upserts in Postgres.
	//
	// So here we "drop" this unique index by name
	pkskIndexNameBroken, err := getIndexName(m, "pksk_split")
	if err != nil {
		return nil, err
	}
	pkskIndexBroken := &indexContext{
		ExcludeNested: true,
		DB: pgdb_v1.Index{
			Name:      pkskIndexNameBroken,
			IsUnique:  false,
			IsDropped: true,
			Method:    pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:   []string{"tenant_id", "pk", "sk"},
		},
		SourceFields: []string{"TenantId", "PK", "SK"},
		Fields:       []*pgs.Field{&tenantIdField, nil, nil},
	}

	pkskIndexName, err := getIndexName(m, "pksk_split2")
	if err != nil {
		return nil, err
	}
	pkskIndex := &indexContext{
		ExcludeNested: true,
		DB: pgdb_v1.Index{
			Name:    pkskIndexName,
			Method:  pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns: []string{"tenant_id", "pk", "sk"},
		},
		SourceFields: []string{"TenantId", "PK", "SK"},
		Fields:       []*pgs.Field{&tenantIdField, nil, nil},
	}

	ftsIndexName, err := getIndexName(m, "fts_data")
	if err != nil {
		return nil, err
	}

	ftsIndex := &indexContext{
		ExcludeNested: true,
		DB: pgdb_v1.Index{
			Name:    ftsIndexName,
			Method:  pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE_GIN,
			Columns: []string{"tenant_id", "fts_data"},
		},
		SourceFields: []string{"FTSData"},
		Fields:       []*pgs.Field{&tenantIdField, nil},
	}

	return []*indexContext{primaryIndex, pkskIndexBroken, pkskIndex, ftsIndex}, nil
}
