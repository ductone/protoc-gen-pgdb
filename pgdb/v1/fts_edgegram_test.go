package v1

import (
	"testing"

	"github.com/clipperhouse/jargon"
	"github.com/ductone/protoc-gen-pgdb/internal/slice"
	"github.com/stretchr/testify/require"
)

func TestNgramStream(t *testing.T) {
	ns3 := edgegramStream(3)
	result, err := jargon.TokenizeString("github").Filter(ns3).Words().ToSlice()
	require.NoError(t, err)
	require.Equal(t, []string{"g", "gi", "git", "gith", "githu", "github", "hub", "ith", "thu"}, tokenStringer(result))

	result, err = jargon.TokenizeString("super sweet").Filter(ns3).Words().ToSlice()
	require.NoError(t, err)
	require.Equal(t, []string{"per", "s", "su", "sup", "supe", "super", "upe", "eet", "s", "sw", "swe", "swee", "sweet", "wee"}, tokenStringer(result))
}

func TestEdgegram_Dedup_ASCIIRepeats(t *testing.T) {
	// Repeated characters produce duplicate windows; ensure dedup + lexicographic order.
	ns2 := edgegramStream(2)
	result, err := jargon.TokenizeString("aaaa").Filter(ns2).Words().ToSlice()
	require.NoError(t, err)
	require.Equal(t, []string{"a", "aa", "aaa", "aaaa"}, tokenStringer(result))
	// Sanity: strictly non-decreasing lexicographic order.
	got := tokenStringer(result)
	for i := 1; i < len(got); i++ {
		if got[i-1] > got[i] {
			t.Fatalf("output not sorted: %v", got)
		}
	}
	// Sanity: uniqueness.
	seen := map[string]struct{}{}
	for _, s := range got {
		if _, ok := seen[s]; ok {
			t.Fatalf("duplicate gram: %q in %v", s, got)
		}
		seen[s] = struct{}{}
	}
}

func TestEdgegram_UnicodeCombining(t *testing.T) {
	// "cafe\u0301" = cafe + combining acute; grapheme cluster at the end.
	ns3 := edgegramStream(3)
	result, err := jargon.TokenizeString("cafe\u0301").Filter(ns3).Words().ToSlice()
	require.NoError(t, err)
	// Expected grams: prefixes + windows (deduped), sorted lexicographically.
	// Graphemes: c a f e\u0301
	// Prefixes: c, ca, caf, cafe\u0301
	// Windows (n=3): caf, afe\u0301 -> caf duplicates; union sorted: afe\u0301, c, ca, caf, cafe\u0301
	require.Equal(t, []string{"afe\u0301", "c", "ca", "caf", "cafe\u0301"}, tokenStringer(result))
}

func TestEdgegram_PunctAndSpacePassthrough(t *testing.T) {
	// Punctuation and spaces should be passed through unchanged between word grams.
	ns3 := edgegramStream(3)
	ts := jargon.TokenizeString("go, go!").Filter(ns3)
	var out []string
	for {
		tk, err := ts.Next()
		require.NoError(t, err)
		if tk == nil {
			break
		}
		out = append(out, tk.String())
	}
	// Should contain the punctuation and space tokens as-is.
	require.Contains(t, out, ",")
	require.Contains(t, out, " ")
	require.Contains(t, out, "!")
}

func tokenStringer(tokens []*jargon.Token) []string {
	return slice.Convert(tokens, func(token *jargon.Token) string {
		return token.String()
	})
}
