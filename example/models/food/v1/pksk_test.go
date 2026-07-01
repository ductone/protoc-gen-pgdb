package v1

import (
	"testing"

	"github.com/stretchr/testify/require"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
)

// TestPKSKMatchesRecord asserts the generated cheap PKSK() helper produces
// exactly what Postgres stores: Record()["pb$pk"] | Record()["pb$sk"] (which is
// what the pb$pksk generated column, `pb$pk || '|' || pb$sk`, evaluates to).
// Covers composite partition keys, sk_const sort keys, and created-at / KSUID
// partitioned tables, across both dialects (the key columns are
// dialect-independent, so both must agree).
func TestPKSKMatchesRecord(t *testing.T) {
	cases := []struct {
		name string
		msg  pgdb_v1.DBReflectMessage
	}{
		{
			name: "composite_pk_with_sk_const",
			msg:  Pasta_builder{TenantId: "t1", Id: "p1"}.Build(),
		},
		{
			name: "pasta_ingredient",
			msg:  PastaIngredient_builder{TenantId: "t1", Id: "pi1", PastaId: "p1", IngredientId: "i1"}.Build(),
		},
		{
			name: "sauce_ingredient",
			msg:  SauceIngredient_builder{TenantId: "t2", Id: "s2", SourceAddr: "1.2.3.4"}.Build(),
		},
		{
			name: "created_at_partitioned",
			msg:  GarlicIngredient_builder{TenantId: "t3", Id: "g3"}.Build(),
		},
		{
			name: "ksuid_partitioned",
			msg:  CheeseIngredient_builder{TenantId: "t4", Id: "c4", EventId: "evt4"}.Build(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, dialect := range []pgdb_v1.Dialect{pgdb_v1.DialectV13, pgdb_v1.DialectV17} {
				m := tc.msg.DBReflect(dialect)
				record, err := m.Record()
				require.NoError(t, err)

				pk, ok := record["pb$pk"].(string)
				require.True(t, ok, "pb$pk should be a string")
				sk, ok := record["pb$sk"].(string)
				require.True(t, ok, "pb$sk should be a string")

				pk2, ok := m.(pgdb_v1.PrimaryKeyer)
				require.True(t, ok, "top-level message should implement PrimaryKeyer")
				require.Equal(t, pk+"|"+sk, pk2.PKSK(),
					"PKSK() must equal Record pb$pk|pb$sk for dialect %s", dialect)
			}
		})
	}
}

// TestPKSKNestedOnly confirms nested-only messages (which have no primary key)
// do not implement PrimaryKeyer — the PKSK method is not generated for them.
func TestPKSKNestedOnly(t *testing.T) {
	m := PastaIngredient_ModelEmbedding_builder{}.Build().DBReflect(pgdb_v1.DialectV13)
	_, ok := m.(pgdb_v1.PrimaryKeyer)
	require.False(t, ok, "nested-only message must not implement PrimaryKeyer")
}
