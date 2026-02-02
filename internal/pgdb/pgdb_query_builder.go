package pgdb

import (
	"fmt"
	"io"
	"sort"
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
	// NestedQueryFields contains all fields from nested query builders that need SafeOperators types generated.
	NestedQueryFields []*safeFieldContext
}

// nestedQueryBuilderContext represents a nested message that should have its own query builder
// for chaining support (e.g., ticketFields.TicketType().Grant().Source().IsExtension().Eq(true)).
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
	// Children are nested query builders within this one (for indexed paths)
	Children []*nestedQueryBuilderContext
	// UnsafeChildren are nested query builders for non-indexed paths
	UnsafeChildren []*nestedQueryBuilderContext
	// Fields are the safe query fields accessible from this nested query builder (indexed fields)
	Fields []*nestedSafeFieldContext
	// UnsafeFields are the non-indexed fields accessible via Unsafe accessors
	UnsafeFields []*nestedSafeFieldContext
	// IsUnsafe indicates this is an unsafe query builder (for non-indexed paths)
	IsUnsafe bool
	// HasIndexedFields indicates this nested message has at least one indexed field
	HasIndexedFields bool
	// HasUnsafeFields indicates this nested message has at least one non-indexed field
	HasUnsafeFields bool
	// SkipTypeDefinition indicates the type struct has already been generated elsewhere,
	// so only accessor methods should be generated. This is used when the same child types
	// are accessible from both safe and unsafe parent query builders.
	SkipTypeDefinition bool
}

// nestedSafeFieldContext wraps safeFieldContext with a short field name for nested query builders.
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

func (module *Module) renderQueryBuilder(ctx pgsgo.Context, w io.Writer, in pgs.File, m pgs.Message, ix *importTracker, generatedOpsTypes map[string]bool, generatedNestedTypes map[string]bool) error {
	c := module.getQueryBuilder(ctx, m, ix, generatedOpsTypes, generatedNestedTypes)
	return templates["query_builder.tmpl"].Execute(w, c)
}

func (module *Module) getQueryBuilder(ctx pgsgo.Context, m pgs.Message, ix *importTracker, generatedOpsTypes map[string]bool, generatedNestedTypes map[string]bool) *qbContext {
	msgFields := module.getMessageFields(ctx, m, ix, "m.self")
	safeFields := module.getSafeFields(ctx, m, ix)
	nestedQueryBuilders := module.getNestedQueryBuilders(ctx, m, ix, safeFields, generatedNestedTypes)

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

	// Collect all unique nested query fields that need SafeOperators types generated.
	// Pass generatedOpsTypes to deduplicate across messages in the same file.
	nestedQueryFields := collectNestedQueryFields(nestedQueryBuilders, filteredSafeFields, generatedOpsTypes)

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
// that are not already in the main safeFields list (i.e., fields that need their SafeOperators types generated).
// generatedOpsTypes tracks types already generated by other messages in the same file to prevent duplicates.
func collectNestedQueryFields(builders []*nestedQueryBuilderContext, existingSafeFields []*safeFieldContext, generatedOpsTypes map[string]bool) []*safeFieldContext {
	// Build a set of existing OpsTypeNames to avoid duplicates within this message
	existingOps := make(map[string]bool)
	for _, f := range existingSafeFields {
		existingOps[f.OpsTypeName] = true
	}

	rv := make([]*safeFieldContext, 0)
	seen := make(map[string]bool)

	var collect func(builders []*nestedQueryBuilderContext)
	collect = func(builders []*nestedQueryBuilderContext) {
		for _, b := range builders {
			// Collect indexed fields
			for _, f := range b.Fields {
				// Skip if already generated by this message, another message, or globally tracked
				if !existingOps[f.OpsTypeName] && !seen[f.OpsTypeName] && !generatedOpsTypes[f.OpsTypeName] {
					seen[f.OpsTypeName] = true
					generatedOpsTypes[f.OpsTypeName] = true
					rv = append(rv, f.safeFieldContext)
				}
			}
			// Collect unsafe (non-indexed) fields as well - they also need SafeOperators types
			for _, f := range b.UnsafeFields {
				if !existingOps[f.OpsTypeName] && !seen[f.OpsTypeName] && !generatedOpsTypes[f.OpsTypeName] {
					seen[f.OpsTypeName] = true
					generatedOpsTypes[f.OpsTypeName] = true
					rv = append(rv, f.safeFieldContext)
				}
			}
			// Recurse into indexed children
			collect(b.Children)
			// Recurse into unsafe children
			collect(b.UnsafeChildren)
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
// generatedNestedTypes tracks all type names that have been generated across all messages in the file
// to prevent duplicates when the same nested type is reachable via different paths.
func (module *Module) getNestedQueryBuilders(ctx pgsgo.Context, m pgs.Message, ix *importTracker, safeFields []*safeFieldContext, generatedNestedTypes map[string]bool) []*nestedQueryBuilderContext {
	msgFields := module.getMessageFields(ctx, m, ix, "m.self")
	parentTypeName := getQueryType(ctx, m)
	msgName := ctx.Name(m).String()

	// Get ALL fields (indexed and non-indexed) for generating query methods
	allFields := module.getMessageFieldsDeep(ctx, m, ix, "m.self", "", "")

	// Build a set of indexed field paths for quick lookup
	indexedFieldPaths := make(map[string]bool)
	for _, sf := range safeFields {
		indexedFieldPaths[sf.ColName] = true
	}

	// Build a map of nested fields and their prefixes
	nestedFieldMap := make(map[string]*nestedFieldInfo)
	for _, f := range msgFields {
		if !f.Nested || f.Field == nil {
			continue
		}
		fieldNum := *f.Field.Descriptor().Number
		prefix := strconv.FormatInt(int64(fieldNum), 10) + "$"
		nestedFieldMap[f.GoName] = &nestedFieldInfo{
			goName:      f.GoName,
			prefix:      prefix,
			field:       f.Field,
			embeddedMsg: f.Field.Type().Embed(),
		}
	}

	// Build the nested query builders from nested fields
	// Sort the keys to ensure deterministic output order
	nestedFieldNames := make([]string, 0, len(nestedFieldMap))
	for name := range nestedFieldMap {
		nestedFieldNames = append(nestedFieldNames, name)
	}
	sort.Strings(nestedFieldNames)

	rv := make([]*nestedQueryBuilderContext, 0)
	for _, name := range nestedFieldNames {
		nf := nestedFieldMap[name]
		if nf.embeddedMsg == nil {
			continue
		}
		builderTypeName := msgName + nf.goName + "QueryBuilder"
		goNamePrefix := nf.goName

		// Get all fields and separate them into indexed and unsafe
		allNestedFields := module.getAllNestedFieldsWithShortName(ctx, allFields, nf.prefix, goNamePrefix, msgName, ix)
		indexedFields, unsafeFields := separateIndexedAndUnsafeFields(allNestedFields, indexedFieldPaths)

		// Recursively build children with index awareness
		children, unsafeChildren := module.getNestedQueryBuildersRecursiveWithIndex(
			ctx, nf.embeddedMsg, ix, builderTypeName, msgName+nf.goName,
			nf.prefix, goNamePrefix, allFields, indexedFieldPaths, generatedNestedTypes,
		)

		nqb := &nestedQueryBuilderContext{
			TypeName:         builderTypeName,
			GoName:           nf.goName,
			Prefix:           nf.prefix,
			FullPrefix:       nf.prefix, // For top-level, full prefix equals local prefix
			GoNamePrefix:     goNamePrefix,
			ParentTypeName:   parentTypeName,
			Fields:           indexedFields,
			UnsafeFields:     unsafeFields,
			Children:         children,
			UnsafeChildren:   unsafeChildren,
			HasIndexedFields: len(indexedFields) > 0 || len(children) > 0,
			HasUnsafeFields:  len(unsafeFields) > 0 || len(unsafeChildren) > 0,
		}
		rv = append(rv, nqb)
	}
	return rv
}

// separateIndexedAndUnsafeFields separates nested fields into indexed (safe) and non-indexed (unsafe) categories.
func separateIndexedAndUnsafeFields(fields []*nestedSafeFieldContext, indexedFieldPaths map[string]bool) ([]*nestedSafeFieldContext, []*nestedSafeFieldContext) {
	indexed := make([]*nestedSafeFieldContext, 0)
	unsafe := make([]*nestedSafeFieldContext, 0)
	for _, f := range fields {
		if indexedFieldPaths[f.ColName] {
			indexed = append(indexed, f)
		} else {
			unsafe = append(unsafe, f)
		}
	}
	return indexed, unsafe
}

type nestedFieldInfo struct {
	goName      string
	prefix      string
	field       pgs.Field
	embeddedMsg pgs.Message
}

// getNestedQueryBuildersRecursiveWithIndex builds nested query builders for deeply nested messages,
// separating children into indexed and unsafe categories.
// generatedNestedTypes tracks all type names that have been generated to prevent duplicates
// when the same type is reachable via different paths.
func (module *Module) getNestedQueryBuildersRecursiveWithIndex(
	ctx pgsgo.Context,
	m pgs.Message,
	ix *importTracker,
	parentTypeName string,
	namePrefix string,
	colPrefix string,
	goNamePrefix string,
	allFields []*fieldContext,
	indexedFieldPaths map[string]bool,
	generatedNestedTypes map[string]bool,
) ([]*nestedQueryBuilderContext, []*nestedQueryBuilderContext) {
	children := make([]*nestedQueryBuilderContext, 0)
	unsafeChildren := make([]*nestedQueryBuilderContext, 0)

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

		// Get all fields and separate them into indexed and unsafe
		allNestedFields := module.getAllNestedFieldsWithShortName(ctx, allFields, fullPrefix, newGoNamePrefix, rootMsgName, ix)
		indexedFields, unsafeFields := separateIndexedAndUnsafeFields(allNestedFields, indexedFieldPaths)

		// Recursively build children
		childChildren, childUnsafeChildren := module.getNestedQueryBuildersRecursiveWithIndex(
			ctx, embeddedMsg, ix, builderTypeName, namePrefix+goName,
			fullPrefix, newGoNamePrefix, allFields, indexedFieldPaths, generatedNestedTypes,
		)

		hasIndexedContent := len(indexedFields) > 0 || len(childChildren) > 0
		hasUnsafeContent := len(unsafeFields) > 0 || len(childUnsafeChildren) > 0

		// Create a query builder for this nested message
		nqb := &nestedQueryBuilderContext{
			TypeName:         builderTypeName,
			GoName:           goName,
			Prefix:           fieldPrefix,
			FullPrefix:       fullPrefix,
			GoNamePrefix:     newGoNamePrefix,
			ParentTypeName:   parentTypeName,
			Fields:           indexedFields,
			UnsafeFields:     unsafeFields,
			Children:         childChildren,
			UnsafeChildren:   childUnsafeChildren,
			HasIndexedFields: hasIndexedContent,
			HasUnsafeFields:  hasUnsafeContent,
		}

		// Add to appropriate list based on whether it has indexed content
		// A nested message goes in the "safe" list if it has ANY indexed content
		// (so users can chain through it to reach indexed fields)
		if hasIndexedContent {
			children = append(children, nqb)
		}
		// Also add to unsafe children if it has unsafe content
		// (so users can access non-indexed fields via UnsafeXxx() accessor)
		if hasUnsafeContent {
			unsafeTypeName := namePrefix + goName + "UnsafeQueryBuilder"

			// Check if this unsafe type was already generated via a different path.
			// This handles the case where the same type is reachable through different
			// intermediate paths (e.g., via oneof members that both contain the same nested type).
			skipUnsafeTypeDef := generatedNestedTypes[unsafeTypeName]
			if !skipUnsafeTypeDef {
				generatedNestedTypes[unsafeTypeName] = true
			}

			// For the children of the unsafe builder:
			// If there's ALSO a safe builder (hasIndexedContent), the safe builder's UnsafeChildren
			// will generate the child types. So the unsafe builder's Children should skip type
			// definitions to avoid duplicates.
			skipChildTypeDef := hasIndexedContent

			// Fix: Copy children with correct ParentTypeName pointing to unsafe type.
			// childUnsafeChildren were built with parentTypeName = builderTypeName (safe type),
			// but when used as children of unsafeNqb, they need ParentTypeName = unsafeTypeName.
			fixedChildren := copyChildrenWithParent(childUnsafeChildren, unsafeTypeName, skipChildTypeDef)

			unsafeNqb := &nestedQueryBuilderContext{
				TypeName:           unsafeTypeName,
				GoName:             goName,
				Prefix:             fieldPrefix,
				FullPrefix:         fullPrefix,
				GoNamePrefix:       newGoNamePrefix,
				ParentTypeName:     parentTypeName,
				Fields:             unsafeFields,
				Children:           fixedChildren,
				IsUnsafe:           true,
				HasUnsafeFields:    true,
				HasIndexedFields:   false,
				SkipTypeDefinition: skipUnsafeTypeDef,
			}
			unsafeChildren = append(unsafeChildren, unsafeNqb)
		}
	}
	return children, unsafeChildren
}

// copyChildrenWithParent creates a deep copy of nested query builder contexts
// with the ParentTypeName updated to the new parent type.
// This is needed because children are built once with a single parentTypeName,
// but may be used in both safe and unsafe query builder contexts which have
// different parent type names.
// If skipTypeDef is true, the copied children will have SkipTypeDefinition=true,
// indicating the type structs are already generated elsewhere.
func copyChildrenWithParent(children []*nestedQueryBuilderContext, newParentTypeName string, skipTypeDef bool) []*nestedQueryBuilderContext {
	if len(children) == 0 {
		return nil
	}
	result := make([]*nestedQueryBuilderContext, len(children))
	for i, child := range children {
		copied := *child // shallow copy
		copied.ParentTypeName = newParentTypeName
		// Only skip type definition if the types are already generated elsewhere
		copied.SkipTypeDefinition = skipTypeDef
		// Recursively copy children with THIS child's TypeName as their parent
		copied.Children = copyChildrenWithParent(child.Children, copied.TypeName, skipTypeDef)
		result[i] = &copied
	}
	return result
}

// isWellKnownType checks if the type is a Google protobuf well-known type.
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
		switch {
		case f.Field != nil:
			shortName = ctx.Name(f.Field).String()
		case f.GoName != "":
			// For virtual fields, use GoName
			shortName = f.GoName
			// Try to strip the prefix if it matches
			if len(shortName) > len(goNamePrefix) && shortName[:len(goNamePrefix)] == goNamePrefix {
				shortName = shortName[len(goNamePrefix):]
			}
		default:
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

// containsNestedPrefix checks if a column name contains a nested prefix (e.g., "5$field" means it's in a nested message).
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
