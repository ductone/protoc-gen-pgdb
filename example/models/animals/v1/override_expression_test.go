package v1

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
)

// The Pet model declares a functional index via override_expression:
//
//	override_expression: "tenant_id, ((\"pb$profile\" ->> 'primary'))"
//
// This asserts the generator threads that expression through to a valid
// CREATE INDEX (verbatim body, partial predicate preserved), without a running
// database.
func TestIndexOverrideExpression(t *testing.T) {
	stmts, err := pgdb_v1.IndexSchema(&Pet{}, pgdb_v1.DialectV13)
	require.NoError(t, err)

	var got string
	for _, s := range stmts {
		if strings.Contains(s, "profile_primary_expr") {
			got = s
			break
		}
	}
	require.NotEmpty(t, got, "expected a CREATE INDEX for profile_primary_expr; got: %v", stmts)

	require.Contains(t, got, "USING\n  BTREE")
	// the raw expression body is emitted verbatim (not a quoted column list)
	require.Contains(t, got, `tenant_id, (("pb$profile" ->> 'primary'))`)
	// partial predicate still applies alongside the override
	require.Contains(t, got, "WHERE pb$deleted_at IS NULL")
}
