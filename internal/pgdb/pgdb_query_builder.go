package pgdb

import (
	"fmt"
	"io"
	"strconv"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	"github.com/ductone/protoc-gen-pgdb/pgdb/v1/xpq"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

type qbContext struct {
	ReceiverType        string
	DbType              string
	QueryType           string
	QueryFields         []*safeFieldContext
	UnsafeType          string
	UnsafeFields        []*fieldContext
	NestedFields        []string
	ColumnType          string
	ColumnFields        []*fieldContext
	NestedQueryBuilders []*nestedQueryBuilderContext
	// NestedQueryFields contains all fields from nested query builders that need SafeOperators types generated
	NestedQueryFields   []*safeFieldContext
}

// nestedQueryBuilderContext represents a nested message that should have its own query builder
// for chaining support (e.g., ticketFields.TicketType().Grant().Source().IsExtension().Eq(true))
type nestedQueryBuilderContext struct {
	// TypeName is the name of the nested query builder type (e.g., "TicketTicketTypeQueryBuilder")
	TypeName string
	// GoName is the accessor method name (e.g., "TicketType")
	GoName string
	// Prefix is the local column prefix for this nested message (e.g., "1$")
	Prefix string
	// FullPrefix is the accumulated column prefix from root (e.g., "13$1$8$")
	FullPrefix string
	// GoNamePrefix is the accumulated Go name prefix to strip from field names (e.g., "ZooShopAnything")
	GoNamePrefix string
	// ParentTypeName is the parent query builder type that will have the accessor method
	ParentTypeName string
	// Children are nested query builders within this one
	Children []*nestedQueryBuilderContext
	// Fields are the safe query fields accessible from this nested query builder
	Fields []*nestedSafeFieldContext
}

// nestedSafeFieldContext wraps safeFieldContext with a short field name for nested query builders
type nestedSafeFieldContext struct {
	*safeFieldContext
	// ShortGoName is the field name without the nested prefix (e.g., "Sfixed64" instead of "ZooShopAnythingSfixed64")
	ShortGoName string
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
	Eq         bool
	Neq        bool
	Gt         bool
	Gte        bool
	Lt         bool
	Lte        bool
	IsNotEmpty bool

	// exp.Inable
	In    bool
	NotIn bool

	// For Inet types, special case a IP Prefix Range matcher (using BETWEEN for btree purposes)
	InNetworkPrefix bool

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
	ArrayNotOverlap  bool
	ArrayContains    bool
	ArrayNotContains bool
	ArrayIsContained bool
	// ArrayEqual       bool  -- covered by equal

	// never safe, or at least we can't understand it yet:
	// exp.Likeable
	// exp.Distinctable
	// exp.Bitwiseable

	// Bit Vector Ops
	Distance bool
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
	safeFields := module.getSafeFields(ctx, m, ix)
	nestedQueryBuilders := module.getNestedQueryBuilders(ctx, m, ix, safeFields)

	// Build a set of GoNames used by nested query builders at the root level
	nestedBuilderNames := make(map[string]bool)
	for _, nqb := range nestedQueryBuilders {
		nestedBuilderNames[nqb.GoName] = true
	}

	// Filter safeFields to exclude ones whose GoName conflicts with nested query builders
	// These fields will be accessible via the nested query builder instead
	filteredSafeFields := make([]*safeFieldContext, 0, len(safeFields))
	for _, sf := range safeFields {
		if !nestedBuilderNames[sf.Field.GoName] {
			filteredSafeFields = append(filteredSafeFields, sf)
		}
	}

	// Collect all unique nested query fields that need SafeOperators types generated
	nestedQueryFields := collectNestedQueryFields(nestedQueryBuilders, filteredSafeFields)

	return &qbContext{
		ReceiverType:        ctx.Name(m).String(),
		DbType:              getDbType(ctx, m),
		QueryType:           getQueryType(ctx, m),
		QueryFields:         filteredSafeFields,
		NestedFields:        getNestedFieldNames(msgFields),
		UnsafeType:          getUnsafeType(ctx, m),
		UnsafeFields:        msgFields,
		ColumnType:          getColumnType(ctx, m),
		ColumnFields:        msgFields,
		NestedQueryBuilders: nestedQueryBuilders,
		NestedQueryFields:   nestedQueryFields,
	}
}

// collectNestedQueryFields recursively collects all unique safeFieldContexts from nested query builders
// that are not already in the main safeFields list (i.e., fields that need their SafeOperators types generated)
func collectNestedQueryFields(builders []*nestedQueryBuilderContext, existingSafeFields []*safeFieldContext) []*safeFieldContext {
	// Build a set of existing OpsTypeNames to avoid duplicates
	existingOps := make(map[string]bool)
	for _, f := range existingSafeFields {
		existingOps[f.OpsTypeName] = true
	}

	rv := make([]*safeFieldContext, 0)
	seen := make(map[string]bool)

	var collect func(builders []*nestedQueryBuilderContext)
	collect = func(builders []*nestedQueryBuilderContext) {
		for _, b := range builders {
			for _, f := range b.Fields {
				if !existingOps[f.OpsTypeName] && !seen[f.OpsTypeName] {
					seen[f.OpsTypeName] = true
					rv = append(rv, f.safeFieldContext)
				}
			}
			collect(b.Children)
		}
	}
	collect(builders)
	return rv
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

// getNestedQueryBuilders builds a tree of nested query builder contexts based on the safe fields.
// It groups fields by their nested message path and creates query builder types for each level.
func (module *Module) getNestedQueryBuilders(ctx pgsgo.Context, m pgs.Message, ix *importTracker, safeFields []*safeFieldContext) []*nestedQueryBuilderContext {
	msgFields := module.getMessageFields(ctx, m, ix, "m.self")
	parentTypeName := getQueryType(ctx, m)
	msgName := ctx.Name(m).String()

	// Get ALL fields (indexed and non-indexed) for generating query methods
	allFields := module.getMessageFieldsDeep(ctx, m, ix, "m.self", "", "")

	// Build a map of nested fields and their prefixes
	nestedFieldMap := make(map[string]*nestedFieldInfo)
	for _, f := range msgFields {
		if !f.Nested || f.Field == nil {
			continue
		}
		fieldNum := *f.Field.Descriptor().Number
		prefix := strconv.FormatInt(int64(fieldNum), 10) + "$"
		nestedFieldMap[f.GoName] = &nestedFieldInfo{
			goName:       f.GoName,
			prefix:       prefix,
			field:        f.Field,
			embeddedMsg:  f.Field.Type().Embed(),
		}
	}

	// Build the nested query builders from nested fields
	rv := make([]*nestedQueryBuilderContext, 0)
	for _, nf := range nestedFieldMap {
		if nf.embeddedMsg == nil {
			continue
		}
		builderTypeName := msgName + nf.goName + "QueryBuilder"
		goNamePrefix := nf.goName
		nqb := &nestedQueryBuilderContext{
			TypeName:       builderTypeName,
			GoName:         nf.goName,
			Prefix:         nf.prefix,
			FullPrefix:     nf.prefix, // For top-level, full prefix equals local prefix
			GoNamePrefix:   goNamePrefix,
			ParentTypeName: parentTypeName,
			Fields:         module.getAllNestedFieldsWithShortName(ctx, allFields, nf.prefix, goNamePrefix, msgName, ix),
			Children:       module.getNestedQueryBuildersRecursive(ctx, nf.embeddedMsg, ix, builderTypeName, msgName+nf.goName, nf.prefix, goNamePrefix, allFields),
		}
		rv = append(rv, nqb)
	}
	return rv
}

type nestedFieldInfo struct {
	goName      string
	prefix      string
	field       pgs.Field
	embeddedMsg pgs.Message
}

// getNestedQueryBuildersRecursive builds nested query builders for deeply nested messages
func (module *Module) getNestedQueryBuildersRecursive(
	ctx pgsgo.Context,
	m pgs.Message,
	ix *importTracker,
	parentTypeName string,
	namePrefix string,
	colPrefix string,
	goNamePrefix string,
	allFields []*fieldContext,
) []*nestedQueryBuilderContext {
	rv := make([]*nestedQueryBuilderContext, 0)

	for _, field := range m.Fields() {
		if field.Type().ProtoType() != pgs.MessageT {
			continue
		}
		embeddedMsg := field.Type().Embed()
		if embeddedMsg == nil {
			continue
		}

		// Check if this is a nested message (not a WKT like Timestamp)
		ext := pgdb_v1.FieldOptions{}
		_, err := field.Extension(pgdb_v1.E_Options, &ext)
		if err != nil {
			continue
		}
		if ext.GetMessageBehavior() == pgdb_v1.FieldOptions_MESSAGE_BEHAVIOR_OMIT ||
			ext.GetMessageBehavior() == pgdb_v1.FieldOptions_MESSAGE_BEHAVIOR_JSONB ||
			ext.GetMessageBehavior() == pgdb_v1.FieldOptions_MESSAGE_BEHAVIOR_VECTOR {
			continue
		}

		// Skip well-known types that are not expanded as nested
		typeName := field.Descriptor().GetTypeName()
		if isWellKnownType(typeName) {
			continue
		}

		fieldNum := *field.Descriptor().Number
		fieldPrefix := strconv.FormatInt(int64(fieldNum), 10) + "$"
		fullPrefix := colPrefix + fieldPrefix

		goName := ctx.Name(field).String()
		builderTypeName := namePrefix + goName + "QueryBuilder"
		newGoNamePrefix := goNamePrefix + goName

		// Extract the root message name from namePrefix (e.g., "AttractionsPet" -> "Attractions")
		rootMsgName := namePrefix[:len(namePrefix)-len(goNamePrefix)]

		nqb := &nestedQueryBuilderContext{
			TypeName:       builderTypeName,
			GoName:         goName,
			Prefix:         fieldPrefix,
			FullPrefix:     fullPrefix,
			GoNamePrefix:   newGoNamePrefix,
			ParentTypeName: parentTypeName,
			Fields:         module.getAllNestedFieldsWithShortName(ctx, allFields, fullPrefix, newGoNamePrefix, rootMsgName, ix),
			Children:       module.getNestedQueryBuildersRecursive(ctx, embeddedMsg, ix, builderTypeName, namePrefix+goName, fullPrefix, newGoNamePrefix, allFields),
		}
		rv = append(rv, nqb)
	}
	return rv
}

// isWellKnownType checks if the type is a Google protobuf well-known type
func isWellKnownType(typeName string) bool {
	switch typeName {
	case ".google.protobuf.Any",
		".google.protobuf.Timestamp",
		".google.protobuf.Duration",
		".google.protobuf.Struct",
		".google.protobuf.BoolValue",
		".google.protobuf.StringValue",
		".google.protobuf.Int32Value",
		".google.protobuf.Int64Value",
		".google.protobuf.UInt32Value",
		".google.protobuf.UInt64Value",
		".google.protobuf.FloatValue",
		".google.protobuf.DoubleValue",
		".google.protobuf.BytesValue":
		return true
	default:
		return false
	}
}

// filterSafeFieldsByPrefix returns safe fields that start with the given prefix
func filterSafeFieldsByPrefix(fields []*safeFieldContext, prefix string) []*safeFieldContext {
	rv := make([]*safeFieldContext, 0)
	for _, f := range fields {
		if len(f.ColName) > len(prefix) && f.ColName[:len(prefix)] == prefix {
			rv = append(rv, f)
		}
	}
	return rv
}

// filterSafeFieldsByPrefixWithShortName returns safe fields that start with the given prefix
// and computes the short Go name by stripping the goNamePrefix
func filterSafeFieldsByPrefixWithShortName(fields []*safeFieldContext, colPrefix string, goNamePrefix string) []*nestedSafeFieldContext {
	rv := make([]*nestedSafeFieldContext, 0)
	for _, f := range fields {
		if len(f.ColName) > len(colPrefix) && f.ColName[:len(colPrefix)] == colPrefix {
			shortName := f.Field.GoName
			// Strip the goNamePrefix to get just the field name
			if len(shortName) > len(goNamePrefix) && shortName[:len(goNamePrefix)] == goNamePrefix {
				shortName = shortName[len(goNamePrefix):]
			}
			rv = append(rv, &nestedSafeFieldContext{
				safeFieldContext: f,
				ShortGoName:      shortName,
			})
		}
	}
	return rv
}

// getAllNestedFieldsWithShortName returns ALL fields (indexed or not) that match the given prefix,
// creating safeFieldContext with basic operations for non-indexed fields.
// This enables querying any nested field without requiring an explicit index.
func (module *Module) getAllNestedFieldsWithShortName(
	ctx pgsgo.Context,
	allFields []*fieldContext,
	colPrefix string,
	goNamePrefix string,
	msgName string,
	ix *importTracker,
) []*nestedSafeFieldContext {
	rv := make([]*nestedSafeFieldContext, 0)

	for _, f := range allFields {
		// Skip if no DB field info
		if f.DB == nil {
			continue
		}

		// Check if this field matches the prefix (is a child of the nested message)
		if len(f.DBFieldNameDeep) <= len(colPrefix) {
			continue
		}
		if f.DBFieldNameDeep[:len(colPrefix)] != colPrefix {
			continue
		}

		// Skip nested message fields - they get their own query builders
		if f.Nested {
			continue
		}

		// Get the remaining column name after the prefix
		remainingColName := f.DBFieldNameDeep[len(colPrefix):]
		// Skip if this field belongs to a deeper nested message (contains another $)
		if containsNestedPrefix(remainingColName) {
			continue
		}

		inputType, err := f.Convert.GoType()
		if err != nil {
			continue
		}

		// Skip fields with invalid or unsupported types (e.g., wrapper types)
		if inputType == "invalid" || inputType == "" {
			continue
		}

		// Create basic ops for all fields (works without index, but may not be optimized)
		ops := basicSafeOps(f.DB.Type, inputType)

		// Get the actual field name from the proto field definition
		// The f.GoName can be unreliable for deeply nested fields due to the message expansion logic
		var shortName string
		if f.Field != nil {
			shortName = ctx.Name(f.Field).String()
		} else if f.GoName != "" {
			// For virtual fields, use GoName
			shortName = f.GoName
			// Try to strip the prefix if it matches
			if len(shortName) > len(goNamePrefix) && shortName[:len(goNamePrefix)] == goNamePrefix {
				shortName = shortName[len(goNamePrefix):]
			}
		} else {
			continue // Skip if we can't determine a name
		}

		// Construct a unique GoName for the OpsTypeName that includes the full path
		fullGoName := goNamePrefix + shortName

		// Track imports for special types
		isJSONB := f.DB.Type == "jsonb"
		isInet := f.DB.Type == "inet"
		isArray := len(f.DB.Type) > 0 && f.DB.Type[0] == '_'
		_, isSupportedArrayType := xpq.SupportedArrayGoTypes[inputType]

		ix.JSON = ix.JSON || isJSONB
		ix.XPQ = ix.XPQ || ops.ObjectAllKeyExists || ops.ObjectAnyKeyExists || (isArray && isSupportedArrayType) || isInet
		ix.NetIP = ix.NetIP || isInet

		rv = append(rv, &nestedSafeFieldContext{
			safeFieldContext: &safeFieldContext{
				InputType:   inputType,
				OpsTypeName: msgName + fullGoName + "SafeOperators",
				Field:       f,
				ColName:     f.DBFieldNameDeep,
				Op:          ops,
			},
			ShortGoName: shortName,
		})
	}
	return rv
}

// containsNestedPrefix checks if a column name contains a nested prefix (e.g., "5$field" means it's in a nested message)
func containsNestedPrefix(colName string) bool {
	for i, c := range colName {
		if c == '$' && i > 0 {
			return true
		}
	}
	return false
}

// basicSafeOps returns a basic set of safe operations for fields without an index.
// These operations work but may not be optimized by an index.
func basicSafeOps(dbType string, inputType string) *safeOps {
	isArray := len(dbType) > 0 && dbType[0] == '_'
	isJSONB := dbType == "jsonb"
	isText := dbType == "text" || dbType == "varchar"
	isInet := dbType == "inet"

	// Check if the input type is supported for array operations
	_, isSupportedArrayType := xpq.SupportedArrayGoTypes[inputType]
	enableArrayOps := isArray && isSupportedArrayType

	return &safeOps{
		Eq:                 true,
		Neq:                true,
		Gt:                 true,
		Gte:                true,
		Lt:                 true,
		Lte:                true,
		In:                 true,
		NotIn:              true,
		IsNull:             true,
		IsNotNull:          true,
		Between:            true,
		NotBetween:         true,
		IsNotEmpty:         isText,
		InNetworkPrefix:    isInet,
		ArrayOverlap:       enableArrayOps,
		ArrayNotOverlap:    enableArrayOps,
		ArrayContains:      enableArrayOps,
		ArrayNotContains:   enableArrayOps,
		ArrayIsContained:   enableArrayOps,
		ObjectContains:     isJSONB,
		ObjectPathExists:   isJSONB,
		ObjectPath:         isJSONB,
		ObjectKeyExists:    isJSONB,
		ObjectAnyKeyExists: isJSONB,
		ObjectAllKeyExists: isJSONB,
	}
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
		isInet := false
		isText := false
		isBits := false

		if f.DB != nil {
			isArray = f.DB.Type[0] == '_'
			isJSONB = f.DB.Type == "jsonb"
			isText = f.DB.Type == "text" || f.DB.Type == "varchar"
			isInet = f.DB.Type == "inet"
			isBits = f.DB.Type == "bit"
		}
		_, isSupportedArrayType := xpq.SupportedArrayGoTypes[inputType]
		ops := safeOpsForIndexTypes(methods, isArray && isSupportedArrayType, isJSONB, isText, isInet, isBits)

		ix.JSON = ix.JSON || isJSONB
		ix.XPQ = ix.XPQ || ops.ObjectAllKeyExists || ops.ObjectAnyKeyExists || (isArray && isSupportedArrayType) || isInet
		ix.NetIP = ix.NetIP || isInet

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

func safeOpsForIndexTypes(input []pgdb_v1.MessageOptions_Index_IndexMethod, isSuportedArrayType bool, isJSONB bool, isText bool, isInet bool, isBits bool) *safeOps {
	indexMethods := make(map[pgdb_v1.MessageOptions_Index_IndexMethod]bool)
	for _, m := range input {
		indexMethods[m] = true
	}
	btree := pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE
	btreeGin := pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE_GIN
	gin := pgdb_v1.MessageOptions_Index_INDEX_METHOD_GIN
	hnswCosine := pgdb_v1.MessageOptions_Index_INDEX_METHOD_HNSW_COSINE

	rv := &safeOps{
		Eq: safeOpCheck(indexMethods, btree, btreeGin, gin),
		// not acutally safe!
		// Neq:                safeOpCheck(indexMethods, btree),
		IsNotEmpty:      safeOpCheck(indexMethods, btree) && isText,
		Gt:              safeOpCheck(indexMethods, btree),
		Gte:             safeOpCheck(indexMethods, btree),
		Lt:              safeOpCheck(indexMethods, btree),
		Lte:             safeOpCheck(indexMethods, btree),
		In:              safeOpCheck(indexMethods, btree),
		InNetworkPrefix: safeOpCheck(indexMethods, btree) && isInet,
		NotIn:           safeOpCheck(indexMethods, btree),
		IsNull:          safeOpCheck(indexMethods, btree),
		IsNotNull:       safeOpCheck(indexMethods, btree),
		Between:         safeOpCheck(indexMethods, btree),
		NotBetween:      safeOpCheck(indexMethods, btree),
		ArrayOverlap:    safeOpCheck(indexMethods, btreeGin, gin) && isSuportedArrayType,
		// This is a bit of a misnomer. It's usually unsafe, but we want to include it if ArrayOverlap exists
		ArrayNotOverlap: safeOpCheck(indexMethods, btreeGin, gin) && isSuportedArrayType,
		ArrayContains:   safeOpCheck(indexMethods, btreeGin, gin) && isSuportedArrayType,
		// This is a bit of a misnomer. It's usually unsafe, but we want to include it if ArrayContains exists
		ArrayNotContains:   safeOpCheck(indexMethods, btreeGin, gin) && isSuportedArrayType,
		ArrayIsContained:   safeOpCheck(indexMethods, btreeGin, gin) && isSuportedArrayType,
		ObjectContains:     safeOpCheck(indexMethods, btreeGin, gin) && isJSONB,
		ObjectPathExists:   safeOpCheck(indexMethods, btreeGin, gin) && isJSONB,
		ObjectPath:         safeOpCheck(indexMethods, btreeGin, gin) && isJSONB,
		ObjectKeyExists:    safeOpCheck(indexMethods, btreeGin, gin) && isJSONB,
		ObjectAnyKeyExists: safeOpCheck(indexMethods, btreeGin, gin) && isJSONB,
		ObjectAllKeyExists: safeOpCheck(indexMethods, btreeGin, gin) && isJSONB,
		Distance:           safeOpCheck(indexMethods, hnswCosine) && isBits,
	}
	return rv
}
