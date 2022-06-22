package pgdb

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"

	"github.com/jackc/pgtype"
	pgs "github.com/lyft/protoc-gen-star"
)

const (
	// If you use an identifier longer than 63 bytes, postgrs will truncate it or error dependingon the context:
	// NOTICE:  identifier "this_constraint_is_going_to_be_longer_than_sixty_three_characters_id_idx" will be truncated to "this_constraint_is_going_to_be_longer_than_sixty_three_characte"
	//
	// docs: https://www.postgresql.org/docs/current/sql-syntax-lexical.html#SQL-SYNTAX-IDENTIFIERS
	//
	// The system uses no more than NAMEDATALEN-1 bytes of an identifier; longer names can be written in commands, but they will be
	// truncated. By default, NAMEDATALEN is 64 so the maximum identifier length is 63 bytes
	postgresNameLen = 63
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

func getColumnName(f pgs.Field) (string, error) {
	return fmt.Sprintf("pb_%d", *f.Descriptor().Number), nil
}

func sha256String(input string) string {
	h := sha256.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var initCachedConnInfo sync.Once
var cachedConnInfo *pgtype.ConnInfo

func pgDataTypeForName(input string) (*pgtype.DataType, bool) {
	initCachedConnInfo.Do(func() {
		cachedConnInfo = pgtype.NewConnInfo()
	})
	return cachedConnInfo.DataTypeForName(input)
}
