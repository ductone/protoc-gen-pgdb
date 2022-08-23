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

func tokenStringer(tokens []*jargon.Token) []string {
	return slice.Convert(tokens, func(token *jargon.Token) string {
		return token.String()
	})
}
