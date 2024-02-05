package pgdb

import (
	"fmt"
	"strings"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type statsContext struct {
	DB            pgdb_v1.Statistic
	ExcludeNested bool
}

func (module *Module) getMessageStatistics(ctx pgsgo.Context, m pgs.Message, ix *importTracker) []*statsContext {
	ext := pgdb_v1.MessageOptions{}
	_, err := m.Extension(pgdb_v1.E_Msg, &ext)
	if err != nil {
		panic(fmt.Errorf("pgdb: getMessageStatistics: failed to extract Message extension from '%s': %w", m.FullyQualifiedName(), err))
	}

	rv := make([]*statsContext, 0)

	usedNames := map[string]struct{}{}
	for _, st := range ext.Stats {
		if _, ok := usedNames[st.Name]; ok {
			panic(fmt.Errorf("pgdb: getMessageStatistics:index name reused  on Message '%s': %s", m.FullyQualifiedName(), st.Name))
		}
		usedNames[st.Name] = struct{}{}
		rv = append(rv, module.renderStats(ctx, m, ix, st))
	}

	return rv
}

func (module *Module) renderStats(ctx pgsgo.Context, m pgs.Message, ix *importTracker, st *pgdb_v1.MessageOptions_Stat) *statsContext {
	statName, err := getIndexName(m, st.Name)
	if err != nil {
		panic(err)
	}
	rv := &statsContext{
		DB: pgdb_v1.Statistic{
			Name:  statName,
			Kinds: st.Kinds,
		},
	}
	if st.Dropped {
		rv.DB.IsDropped = true
		return rv
	}

	for _, fieldName := range st.Columns {
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
				panic(fmt.Errorf("could not find field for stat: %s", path))
			}

			resolution += name
			rv.DB.Columns = append(rv.DB.Columns, resolution)
		}
	}
	return rv
}
