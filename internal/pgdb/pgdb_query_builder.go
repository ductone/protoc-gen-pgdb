package pgdb

import (
	"fmt"
	"io"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	"github.com/ductone/protoc-gen-pgdb/pgdb/v1/xpq"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
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
//
//		   && (anyarray,anyarray)			Overlap
//		   @> (anyarray,anyarray)			ArrayContains
//		   <@ (anyarray,anyarray)			ArrayIsContained
//		   = (anyarray,anyarray)			ArrayEqual (not aliased, use equal)
//
//		   tsvector_ops	@@ (tsvector,tsquery)
//		   @@@ (tsvector,tsquery)
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

	// Array Ops
	ArrayOverlap     bool
	ArrayContains    bool
	ArrayIsContained bool
	// ArrayEqual       bool  -- covered by equal

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
}

func (module *Module) renderQueryBuilder(ctx pgsgo.Context, w io.Writer, in pgs.File, m pgs.Message, ix *importTracker) error {
	c := module.getQueryBuilder(ctx, m, ix)
	return templates["query_builder.tmpl"].Execute(w, c)
}

func (module *Module) getQueryBuilder(ctx pgsgo.Context, m pgs.Message, ix *importTracker) *qbContext {
	msgFields := module.getMessageFields(ctx, m, ix, "m.self")
	return &qbContext{
		ReceiverType: ctx.Name(m).String(),
		DbType:       getDbType(ctx, m),
		QueryType:    getQueryType(ctx, m),
		QueryFields:  module.getSafeFields(ctx, m, ix),
		NestedFields: getNestedFieldNames(msgFields),
		UnsafeType:   getUnsafeType(ctx, m),
		UnsafeFields: msgFields,
		ColumnType:   getColumnType(ctx, m),
		ColumnFields: msgFields,
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

func (module *Module) getSafeFields(ctx pgsgo.Context, m pgs.Message, ix *importTracker) []*safeFieldContext {
	fields := module.getMessageFieldsDeep(ctx, m, ix, "m.self", "", "")
	rv := make([]*safeFieldContext, 0, len(fields))

	// uniquify index by path or else we'll double generate methods
	allIndexes := module.getMessageIndexes(ctx, m, &importTracker{})
	indexByFullName := make(map[string][]pgdb_v1.MessageOptions_Index_IndexMethod)
	missingIndices := map[string]bool{}
	for _, idx := range allIndexes {
		for _, f := range idx.DB.Columns {
			if _, ok := indexByFullName[f]; !ok {
				indexByFullName[f] = []pgdb_v1.MessageOptions_Index_IndexMethod{}
			}
			indexByFullName[f] = append(indexByFullName[f], idx.DB.Method)
			missingIndices[f] = true
		}
	}

	for _, f := range fields {
		methods := indexByFullName[f.DBFieldNameDeep]
		if len(methods) == 0 {
			continue
		}
		if f.GoName == "" {
			panic(fmt.Errorf("missing goName for field context: %s:\n%v", m.Name(), f))
		}

		inputType, err := f.Convert.GoType()
		if err != nil {
			panic(err)
		}

		delete(missingIndices, f.DBFieldNameDeep)

		isArray := false
		isJSONB := false

		if f.DB != nil {
			isArray = f.DB.Type[0] == '_'
			isJSONB = f.DB.Type == "jsonb"
		}
		_, isSupportedArrayType := xpq.SupportedArrayGoTypes[inputType]
		ops := safeOpsForIndexTypes(methods, isArray && isSupportedArrayType, isJSONB)

		ix.JSON = ix.JSON || isJSONB
		ix.XPQ = ix.XPQ || ops.ObjectAllKeyExists || ops.ObjectAnyKeyExists || (isArray && isSupportedArrayType)

		rv = append(rv, &safeFieldContext{
			InputType:   inputType,
			OpsTypeName: ctx.Name(m).String() + f.GoName + "SafeOperators",
			Field:       f,
			ColName:     f.DBFieldNameDeep,
			Op:          ops,
		})
	}
	if len(missingIndices) != 0 {
		panic(fmt.Errorf("did not find some indexes: %v", missingIndices))
	}
	return rv
}

func safeOpsForIndexTypes(input []pgdb_v1.MessageOptions_Index_IndexMethod, isSuportedArrayType bool, isJSONB bool) *safeOps {
	indexMethods := make(map[pgdb_v1.MessageOptions_Index_IndexMethod]bool)
	for _, m := range input {
		indexMethods[m] = true
	}
	btree := pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE
	btreeGin := pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE_GIN
	gin := pgdb_v1.MessageOptions_Index_INDEX_METHOD_GIN

	rv := &safeOps{
		Eq:                 safeOpCheck(indexMethods, btree, btreeGin, gin),
		Neq:                safeOpCheck(indexMethods, btree),
		Gt:                 safeOpCheck(indexMethods, btree),
		Gte:                safeOpCheck(indexMethods, btree),
		Lt:                 safeOpCheck(indexMethods, btree),
		Lte:                safeOpCheck(indexMethods, btree),
		In:                 safeOpCheck(indexMethods, btree),
		NotIn:              safeOpCheck(indexMethods, btree),
		IsNull:             safeOpCheck(indexMethods, btree),
		IsNotNull:          safeOpCheck(indexMethods, btree),
		Between:            safeOpCheck(indexMethods, btree),
		NotBetween:         safeOpCheck(indexMethods, btree),
		ArrayOverlap:       safeOpCheck(indexMethods, btreeGin, gin) && isSuportedArrayType,
		ArrayContains:      safeOpCheck(indexMethods, btreeGin, gin) && isSuportedArrayType,
		ArrayIsContained:   safeOpCheck(indexMethods, btreeGin, gin) && isSuportedArrayType,
		ObjectContains:     safeOpCheck(indexMethods, btreeGin, gin) && isJSONB,
		ObjectPathExists:   safeOpCheck(indexMethods, btreeGin, gin) && isJSONB,
		ObjectPath:         safeOpCheck(indexMethods, btreeGin, gin) && isJSONB,
		ObjectKeyExists:    safeOpCheck(indexMethods, btreeGin, gin) && isJSONB,
		ObjectAnyKeyExists: safeOpCheck(indexMethods, btreeGin, gin) && isJSONB,
		ObjectAllKeyExists: safeOpCheck(indexMethods, btreeGin, gin) && isJSONB,
	}
	return rv
}
