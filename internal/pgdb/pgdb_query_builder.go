package pgdb

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ductone/protoc-gen-pgdb/internal/slice"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type qbContext struct {
	ReceiverType string
	DbType       string
	QueryType    string
	QueryFields  []*safeFieldContext
	UnsafeType   string
	UnsafeFields []*fieldContext
	NestedFields []string
	ColumnType   string
	ColumnFields []*fieldContext
}

// a set of SQL Operations that we believe are safe, depending upon the Indexs on a col.
//
// upstream docs:
// https://www.postgresql.org/docs/14/indexes-types.html
// btree: <   <=   =   >=   >
//
//	BETWEEN, IN, IS NULL, IS NOT NULL
//
// btree_gin: basically gin?  we use btree_gin where the first col is tenant_id
//
// gin: <@   @>   =   &&
//
//		https://www.postgresql.org/docs/14/gin-builtin-opclasses.html#GIN-BUILTIN-OPCLASSES-TABLE
//
//	   Postgres Operators           	Function Name in Go
//		   @> (jsonb,jsonb)            		ObjectContains(interface{}) // marshaled to JSON
//		   @? (jsonb,jsonpath)   			ObjectPathExists(string)
//		   @@ (jsonb,jsonpath)				ObjectPath(string)
//		   ? (jsonb,text)					ObjectKeyExists(string)
//		   ?| (jsonb,text[])				ObjectAnyKeyExists([]string)
//		   ?& (jsonb,text[])				ObjectAllKeyExists([]string)
type safeOps struct {
	// exp.Comparable
	Eq  bool
	Neq bool
	Gt  bool
	Gte bool
	Lt  bool
	Lte bool

	// exp.Inable
	In    bool
	NotIn bool

	// exp.Isable -- we only export a subset to support IS NULL / NOT NULL, use Eq for equality
	IsNull    bool
	IsNotNull bool

	// exp.Rangeable
	Between    bool
	NotBetween bool

	// Postgres JSONB Operators
	ObjectContains     bool
	ObjectPathExists   bool
	ObjectPath         bool
	ObjectKeyExists    bool
	ObjectAnyKeyExists bool
	ObjectAllKeyExists bool

	// never safe, or at least we can't understand it yet:
	// exp.Likeable
	// exp.Distinctable
	// exp.Bitwiseable
}

type safeFieldContext struct {
	Field       *fieldContext
	OpsTypeName string
	InputType   string
	ColName     string
	Op          *safeOps
	Prefix      string
}

func (module *Module) renderQueryBuilder(ctx pgsgo.Context, w io.Writer, in pgs.File, m pgs.Message, ix *importTracker) error {
	c := module.getQueryBuilder(ctx, m, ix)
	return templates["query_builder.tmpl"].Execute(w, c)
}

func (module *Module) getQueryBuilder(ctx pgsgo.Context, m pgs.Message, ix *importTracker) *qbContext {
	nsgFields := module.getMessageFields(ctx, m, ix, "m.self")
	// fullFields := module.getMessageFieldsFull(ctx, m, ix, "m.self", "")
	// for _, f := range fullFields {
	// 	//
	// 	if f.DB != nil && f.Field != nil {
	// 		name, _ := getColumnName(f.Field)
	// 		if f.GoName == "Medium" {
	// 			fmt.Fprintf(os.Stderr, "🌊🌊 %s %v\n", *m.Descriptor().Name, f)
	// 		}
	// 		fmt.Fprintf(os.Stderr, "🌊 %s %s%s\n", *m.Descriptor().Name, f.Prefix_, name)
	// 	} else {
	// 		if f.GoName == "Medium" {
	// 			fmt.Fprintf(os.Stderr, "🌊🌊 %s %v\n", *m.Descriptor().Name, f)
	// 		}
	// 		fmt.Fprintf(os.Stderr, "🌊 %s %s%s\n", *m.Descriptor().Name, f.Prefix_, f.GoName)
	// 	}

	// }

	return &qbContext{
		ReceiverType: ctx.Name(m).String(),
		DbType:       getDbType(ctx, m),
		QueryType:    getQueryType(ctx, m),
		QueryFields:  module.getSafeFields(ctx, m, nsgFields, ix),
		NestedFields: getNesteFieldNames(nsgFields),
		UnsafeType:   getUnsafeType(ctx, m),
		UnsafeFields: nsgFields,
		ColumnType:   getColumnType(ctx, m),
		ColumnFields: nsgFields,
	}
}

func getDbType(ctx pgsgo.Context, m pgs.Message) string {
	return ctx.Name(m).String() + "DB"
}

func getQueryType(ctx pgsgo.Context, m pgs.Message) string {
	return ctx.Name(m).String() + "DBQueryBuilder"
}

func getUnsafeType(ctx pgsgo.Context, m pgs.Message) string {
	return ctx.Name(m).String() + "DBQueryUnsafe"
}

func getColumnType(ctx pgsgo.Context, m pgs.Message) string {
	return ctx.Name(m).String() + "DBColumns"
}

func safeOpCheck(indexMethods map[pgdb_v1.MessageOptions_Index_IndexMethod]bool, methods ...pgdb_v1.MessageOptions_Index_IndexMethod) bool {
	for _, m := range methods {
		if indexMethods[m] {
			return true
		}
	}
	return false
}

func (module *Module) getSafeFields(ctx pgsgo.Context, m pgs.Message, fields []*fieldContext, ix *importTracker) []*safeFieldContext {
	rv := make([]*safeFieldContext, 0, len(fields))
	// todo(pquerna): not ideal, little weird way to do this.
	allIndexes := module.getMessageIndexes(ctx, m, &importTracker{})
	indexByFullName := make(map[string][]*indexContext)
	for _, idx := range allIndexes {
		for _, f := range idx.SourceFields {
			indexByFullName[f] = append(indexByFullName[f], idx)
		}
	}
	missingIndices := map[string]bool{}
	for key := range indexByFullName {
		missingIndices[key] = true
	}

	// getField
	for _, f := range fields {
		inputType, err := f.Convert.GoType()
		if err != nil {
			panic(err)
		}
		// asdf, _ := getColumnName(f.Field)
		// fmt.Fprintf(os.Stderr, "🌮 %s: (safe field) %v %v %v\n", *m.Descriptor().Name, f.GoName, f.DB, f.DataType)
		indexTypes := slice.Convert(indexByFullName[f.GoName], func(ic *indexContext) pgdb_v1.MessageOptions_Index_IndexMethod {
			return ic.DB.Method
		})

		if len(indexTypes) == 0 {
			continue
		}
		delete(missingIndices, f.GoName)
		fieldName := f.GoName
		if !f.IsVirtual {
			fieldName = ctx.Name(f.Field).String()
		}

		ops := safeOpsForIndexTypes(indexTypes)
		if ops.ObjectContains {
			ix.JSON = true
		}
		if ops.ObjectAllKeyExists || ops.ObjectAnyKeyExists {
			ix.XPQ = true
		}
		rv = append(rv, &safeFieldContext{
			InputType:   inputType,
			OpsTypeName: ctx.Name(m).String() + fieldName + "SafeOperators",
			Field:       f,
			ColName:     f.DB.Name,
			Op:          ops,
		})
	}

	if len(missingIndices) == 0 {
		return rv
	}

	allFields := module.getMessageFieldsFull(ctx, m, ix, "m.self", "")
	for key := range missingIndices {
		if key == "tenant_id" {
			continue
		}
		found := false
		for _, f := range allFields {
			if f.Field == nil {
				continue
			}
			name, err := getColumnName(f.Field)
			if err != nil {
				panic(err)
			}
			if key != f.Prefix_+name {
				continue
			}
			inputType, err := f.Convert.GoType()
			if err != nil {
				panic(err)
			}

			for _, ic := range indexByFullName[key] {
				indexTypes := []pgdb_v1.MessageOptions_Index_IndexMethod{ic.DB.Method}
				ops := safeOpsForIndexTypes(indexTypes)
				if ops.ObjectContains {
					ix.JSON = true
				}
				if ops.ObjectAllKeyExists || ops.ObjectAnyKeyExists {
					ix.XPQ = true
				}
				var sb strings.Builder

				_, _ = sb.WriteString(ctx.Name(m).String())
				for i, sf := range ic.RawColumns {
					for _, s := range strings.Split(sf, "🌮") {
						_, _ = sb.WriteString(cases.Title(language.AmericanEnglish).String(s))
					}
					if i != len(ic.RawColumns)-1 {
						_, _ = sb.WriteString("_And_")
					}
				}
				_, _ = sb.WriteString("SafeOperators")
				fieldName := sb.String()

				prefix := ""
				for _, sf := range ic.SourceFields {
					if sf == "tenant_id" {
						continue
					}
					prefix = sf
				}
				fmt.Fprintf(os.Stderr, "🌮 found it!! %s:%s -> (%s.%s)\n", ctx.Name(m).String(), key, ctx.Name(m).String(), fieldName)

				rv = append(rv, &safeFieldContext{
					InputType:   inputType,
					OpsTypeName: ctx.Name(m).String() + fieldName + "SafeOperators",
					Field:       f,
					ColName:     prefix,
					Op:          ops,
					Prefix:      prefix,
				})
			}

			found = true
			break
		}
		if !found {
			panic(fmt.Errorf("cant resolve index: %s", key))
		}
	}

	return rv
}

func safeOpsForIndexTypes(input []pgdb_v1.MessageOptions_Index_IndexMethod) *safeOps {
	indexMethods := make(map[pgdb_v1.MessageOptions_Index_IndexMethod]bool)
	for _, m := range input {
		indexMethods[m] = true
	}
	btree := pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE
	btreeGin := pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE_GIN
	gin := pgdb_v1.MessageOptions_Index_INDEX_METHOD_GIN

	rv := &safeOps{
		Eq:  safeOpCheck(indexMethods, btree, btreeGin, gin),
		Neq: safeOpCheck(indexMethods, btree),
		Gt:  safeOpCheck(indexMethods, btree),
		Gte: safeOpCheck(indexMethods, btree),
		Lt:  safeOpCheck(indexMethods, btree),
		Lte: safeOpCheck(indexMethods, btree),

		In:    safeOpCheck(indexMethods, btree),
		NotIn: safeOpCheck(indexMethods, btree),

		IsNull:    safeOpCheck(indexMethods, btree),
		IsNotNull: safeOpCheck(indexMethods, btree),

		Between:    safeOpCheck(indexMethods, btree),
		NotBetween: safeOpCheck(indexMethods, btree),

		ObjectContains:     safeOpCheck(indexMethods, btreeGin, gin),
		ObjectPathExists:   safeOpCheck(indexMethods, btreeGin, gin),
		ObjectPath:         safeOpCheck(indexMethods, btreeGin, gin),
		ObjectKeyExists:    safeOpCheck(indexMethods, btreeGin, gin),
		ObjectAnyKeyExists: safeOpCheck(indexMethods, btreeGin, gin),
		ObjectAllKeyExists: safeOpCheck(indexMethods, btreeGin, gin),
	}
	return rv
}
