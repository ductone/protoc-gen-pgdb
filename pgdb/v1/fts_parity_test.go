package v1

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/stretchr/testify/require"
)

// TestFTSCorpusParity pins the tsvector output for a shared corpus. The same
// corpus + expected outputs are committed in the Rust protoc-gen-pgdb-rs repo
// (crates/pgdb/tests/fixtures), whose FTS port must produce byte-identical
// tsvectors so full-text search works across both languages.
func TestFTSCorpusParity(t *testing.T) {
	corpusF, err := os.Open("testdata/fts_corpus.jsonl")
	require.NoError(t, err)
	defer corpusF.Close()
	expectedF, err := os.Open("testdata/fts_expected.jsonl")
	require.NoError(t, err)
	defer expectedF.Close()

	type inputDoc struct {
		Type   int32  `json:"type"`
		Weight int32  `json:"weight"`
		Value  string `json:"value"`
	}
	type expected struct {
		TSVector string `json:"tsvector"`
		Empty    bool   `json:"empty"`
	}

	corpus := bufio.NewScanner(corpusF)
	corpus.Buffer(make([]byte, 1024*1024), 1024*1024)
	exps := bufio.NewScanner(expectedF)
	exps.Buffer(make([]byte, 4*1024*1024), 4*1024*1024)

	line := 0
	for corpus.Scan() {
		line++
		require.True(t, exps.Scan(), "expected file shorter than corpus")

		var docs []inputDoc
		require.NoError(t, json.Unmarshal(corpus.Bytes(), &docs))
		var want expected
		require.NoError(t, json.Unmarshal(exps.Bytes(), &want))

		sc := make([]*SearchContent, 0, len(docs))
		for _, d := range docs {
			sc = append(sc, &SearchContent{
				Type:   FieldOptions_FullTextType(d.Type),
				Weight: FieldOptions_FullTextWeight(d.Weight),
				Value:  d.Value,
			})
		}
		e := FullTextSearchVectors(sc)
		ds := goqu.Dialect("postgres").From("t").Where(exp.NewLiteralExpression("?", e)).Prepared(true)
		_, params, err := ds.ToSQL()
		require.NoError(t, err)

		if want.Empty {
			require.Empty(t, params, "line %d: expected empty tsvector literal", line)
			continue
		}
		require.Len(t, params, 1, "line %d", line)
		got, _ := params[0].(string)
		require.Equal(t, want.TSVector, got, "line %d: tsvector mismatch", line)
	}
}
