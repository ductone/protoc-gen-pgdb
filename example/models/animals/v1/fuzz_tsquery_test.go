package v1

import (
	"context"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/ductone/protoc-gen-pgdb/internal/pgtest"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	"github.com/stretchr/testify/require"
)

func FuzzFullTextSearchQuery(f *testing.F) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(f, err)
	defer pg.Stop()

	testcases := []string{"Hello, world", " ", "hello:world", `hello\:world`, `snowman: ☃`}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}

	qb := goqu.Dialect("postgres")

	f.Fuzz(func(t *testing.T, input string) {
		qexp := pgdb_v1.FullTextSearchQuery(input)
		query, params, err := qb.Select(qexp).Prepared(true).ToSQL()
		if err != nil {
			t.Errorf("Failed to generate query: input: %q  error: %q ", input, err)
		}
		_, err = pg.DB.Exec(ctx, query, params...)
		require.NoError(t, err, "query failed with input: %q\n\n%q\n\n%q\n\n", input, query, params)
	})
}
