package xpq

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBitsValue(t *testing.T) {
	cases := []struct {
		name string
		in   BitsValue
		want any
	}{
		{
			name: "nil renders as NULL",
			in:   nil,
			want: nil,
		},
		{
			name: "zero length renders as NULL",
			in:   BitsValue{},
			want: nil,
		},
		{
			name: "single byte renders bits msb first",
			in:   BitsValue{0b10100000},
			want: "10100000",
		},
		{
			name: "two bytes render every bit",
			in:   BitsValue{0x00, 0xFF},
			want: "0000000011111111",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := c.in.Value()
			require.NoError(t, err)
			require.Equal(t, c.want, got)
		})
	}
}

// A sized bit column is the only thing BitsValue feeds, so an empty operand must
// not reach pgvector as B” — that throws "different bit lengths N and 0".
func TestBitsValueEmptyIsNotEmptyBitString(t *testing.T) {
	got, err := BitsValue{}.Value()
	require.NoError(t, err)
	require.Nil(t, got)

	full, err := BitsValue(strings.Repeat("\x00", 512)).Value()
	require.NoError(t, err)
	require.Len(t, full, 4096)
}
