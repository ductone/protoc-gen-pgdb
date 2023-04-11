package pgdb

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
	pgs "github.com/lyft/protoc-gen-star"
)

const (
	// If you use an identifier longer than 63 bytes, postgrs will truncate it or error dependingon the context:
	// NOTICE:  identifier "this_constraint_is_going_to_be_longer_than_sixty_three_characters_id_idx" will be truncated to "this_constraint_is_going_to_be_longer_than_sixty_three_characte"
	//
	// docs: https://www.postgresql.org/docs/current/sql-syntax-lexical.html#SQL-SYNTAX-IDENTIFIERS
	//
	// The system uses no more than NAMEDATALEN-1 bytes of an identifier; longer names can be written in commands, but they will be
	// truncated. By default, NAMEDATALEN is 64 so the maximum identifier length is 63 bytes.
	//
	// we subtract 15 more just to have more buffer for prefixes that are calculated at runtime (eg, nested fields).
	postgresNameLen = (63 - 15)
)

func getTableName(m pgs.Message) (string, error) {
	fqnHash := sha256String(m.FullyQualifiedName())

	pkgName := m.Package().ProtoName().LowerSnakeCase().String()
	msgName := m.Name().LowerSnakeCase().String()
	proposed := strings.Join([]string{"pb", msgName, pkgName}, "_")
	shortHash := fqnHash[0:8]
	// shorten to <63 with enough room to append short hash
	proposed = proposed[0:min(postgresNameLen-(len(shortHash)+1), len(proposed))]
	proposed = strings.ToLower(strings.TrimSuffix(proposed, "_")) + "_" + shortHash
	return proposed, nil
}

func getNestedName(f pgs.Field) string {
	return strconv.FormatInt(int64(*f.Descriptor().Number), 10) + "$"
}

func getColumnName(f pgs.Field) (string, error) {
	buf := strings.Builder{}
	// TOOD(pquerna): figure out prefix?
	_, _ = buf.WriteString(f.Name().LowerSnakeCase().String())

	// if the field name is too long, convert to number version
	if len(buf.String()) < postgresNameLen {
		return buf.String(), nil
	}
	buf.Reset()

	_, _ = buf.WriteString(strconv.FormatInt(int64(*f.Descriptor().Number), 10))
	if len(buf.String()) < postgresNameLen {
		return buf.String(), nil
	}
	panic(fmt.Errorf("pgdb: getColumnName: can't find short enough name for %v", f.FullyQualifiedName()))
}

func getColumnOneOfName(f pgs.OneOf) (string, error) {
	buf := strings.Builder{}
	_, _ = buf.WriteString(f.Name().LowerSnakeCase().String())
	_, _ = buf.WriteString("_oneof")

	// if the field name is too long, convert to number version
	if len(buf.String()) < postgresNameLen {
		return buf.String(), nil
	}

	// oneofs don't have numeric indexes.  we could hash the oneof name, but lets see if this is actually ever needed.
	panic(fmt.Errorf("pgdb: getColumnOneOfName: can't find short enough name for %v", f.FullyQualifiedName()))
}

func getIndexName(m pgs.Message, name string) (string, error) {
	fqnHash := sha256String(m.FullyQualifiedName() + "$" + name)

	pkgName := m.Package().ProtoName().LowerSnakeCase().String()
	msgName := m.Name().LowerSnakeCase().String()
	proposed := strings.Join([]string{name, msgName, pkgName}, "_")
	shortHash := fqnHash[0:8]
	// shorten to <63 with enough room to append short hash
	proposed = proposed[0:min(postgresNameLen-(len(shortHash)+1), len(proposed))]
	proposed = strings.ToLower(strings.TrimSuffix(proposed, "_")) + "_" + shortHash
	return proposed, nil
}

func sha256String(input string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var initCachedConnInfo sync.Once
var cachedConnInfo *pgtype.Map

func pgDataTypeForName(input string) *pgtype.Type {
	initCachedConnInfo.Do(func() {
		cachedConnInfo = pgtype.NewMap()
	})
	rv, ok := cachedConnInfo.TypeForName(input)
	if !ok {
		panic("faild to find postgres type for '" + input + "'")
	}
	return rv
}
