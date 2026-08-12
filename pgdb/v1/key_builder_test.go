package v1

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeyBuilder_StripsNul(t *testing.T) {
	nul := string(rune(0))

	tests := []struct {
		name string
		in   []string
		want string
	}{
		{"clean", []string{"tenant", ":", "user"}, "tenant:user"},
		{"nul mid-part", []string{"us" + nul + "er"}, "user"},
		{"nul at edges", []string{nul + "user" + nul}, "user"},
		{"multiple nuls", []string{"a" + nul + nul + "b" + nul + "c"}, "abc"},
		{"nul-only part", []string{nul, nul + nul}, ""},
		{"empty part", []string{""}, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var sb KeyBuilder
			for _, part := range tc.in {
				n, err := sb.WriteString(part)
				require.NoError(t, err)
				require.Equal(t, len(part)-countNul(part), n)
			}
			require.Equal(t, tc.want, sb.String())

			sb.Reset()
			require.Equal(t, "", sb.String())
		})
	}
}

func countNul(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] == 0 {
			n++
		}
	}
	return n
}
