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

	for _, fieldName := range idx.Columns {
		path := strings.Split(fieldName, ".")
		message := m
		resolution := ""
		for i, p := range path {
			lastP := i == len(path)-1

			if !lastP {
				f := fieldByName(message, p)
				resolution += getNestedName(f)
				message = f.Type().Embed()
				continue
			}

			name := ""
			// could be a real field!
			if f, ok := tryFieldByName(message, p); ok {
				name, err = getColumnName(f)
				if err != nil {
					panic(err)
				}
			} else {
				// look in oneofs!
				for _, oo := range message.RealOneOfs() {
					if oo.Name().String() == p {
						name, err = getColumnOneOfName(oo)
						if err != nil {
							panic(err)
						}
						break
					}
				}
			}
			if name == "" {
				panic(fmt.Errorf("could not find field for index: %s", path))
			}

			resolution += name
			rv.DB.Columns = append(rv.DB.Columns, resolution)
		}
	}

	if idx.PartialDeletedAtIsNull {
		if f, ok := tryFieldByName(m, "deleted_at"); ok {
			name, err := getColumnName(f)
			if err != nil {
				panic(err)
			}
			_ = name
			rv.DB.WherePredicate = fmt.Sprintf(
				`" + io.ColumnName("%s") + " IS NULL`,
				name,
			)
		} else {
			panic(fmt.Sprintf("%s ould not find field for partial index: deleted_at", m.FullyQualifiedName()))
		}
	}
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

	primaryIndex := &indexContext{
		ExcludeNested: true,
		DB: pgdb_v1.Index{
			Name:      primaryIndexName,
			IsPrimary: true,
			IsUnique:  true,
			Method:    pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:   []string{"tenant_id", "pksk"},
		},
	}

	if fext.PartitionedByCreatedAt {
		primaryIndex.DB.Columns = append(primaryIndex.DB.Columns, "created_at")
	}

	if fext.PartitionedByEventId {
		primaryIndex.DB.Columns = append(primaryIndex.DB.Columns, "event_id")
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
	}

	// again loop through and look for vector / hnsw indexes

	rv := []*indexContext{primaryIndex, pkskIndexBroken, pkskIndex, ftsIndex}

	// iterate message for vector behavior options
	for _, field := range m.Fields() {
		ext := pgdb_v1.FieldOptions{}
		_, err := field.Extension(pgdb_v1.E_Options, &ext)
		if err != nil {
			return nil, fmt.Errorf("pgdb: getField: failed to extract Message extension from '%s': %w", field.FullyQualifiedName(), err)
		}
		if ext.MessageBehavior != pgdb_v1.FieldOptions_MESSAGE_BEHAVIOR_VECTOR {
			continue
		}

		enumField, _, err := GetFieldVectorShape(field)
		if err != nil {
			return nil, err
		}

		pgColName, err := getColumnName(field)
		if err != nil {
			panic(fmt.Errorf("pgdb: getColumnName failed for: %v: %s (of type %s)",
				field.Type().ProtoType(), field.FullyQualifiedName(), field.Descriptor().GetType()))
		}

		// enum values
		for _, enumValue := range enumField.Type().Enum().Values() {
			if enumValue.Value() == 0 {
				// skip the zero value
				continue
			}

			vectorIndexName, err := getIndexName(m, fmt.Sprintf("vector_index_%s", enumValue.Name().String()))
			if err != nil {
				return nil, err
			}

			vectorCol := fmt.Sprintf("%s_%d", pgColName, enumValue.Value())

			tempCtx := &indexContext{
				ExcludeNested: true,
				DB: pgdb_v1.Index{
					Name:   vectorIndexName,
					Method: pgdb_v1.MessageOptions_Index_INDEX_METHOD_HNSW_COSINE,
					Columns: []string{
						vectorCol,
					},
					OverrideExpression: fmt.Sprintf("pb$%s vector_cosine_ops", vectorCol),
				},
			}
			rv = append(rv, tempCtx)
		}

		break
	}

	return rv, nil
}
