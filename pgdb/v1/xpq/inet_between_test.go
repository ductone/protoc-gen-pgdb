package xpq

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIPNetAddressRange(t *testing.T) {
	cases := []struct {
		cidr  string
		start string
		end   string
	}{
		{
			cidr:  "192.168.1.0/24",
			start: "192.168.1.0",
			end:   "192.168.1.255",
		},
		{
			cidr:  "192.168.1.1/24",
			start: "192.168.1.0",
			end:   "192.168.1.255",
		},
		{
			cidr:  "192.168.1.1/8",
			start: "192.0.0.0",
			end:   "192.255.255.255",
		},
		{
			cidr:  "0.0.0.0/0",
			start: "0.0.0.0",
			end:   "255.255.255.255",
		},
		{
			cidr:  "2001:db8:abcd:12::0/64",
			start: "2001:db8:abcd:12::",
			end:   "2001:db8:abcd:12:ffff:ffff:ffff:ffff",
		},
	}

	for _, c := range cases {
		cdr, err := netip.ParsePrefix(c.cidr)
		require.NoError(t, err)

		start, end := NetworkRange(cdr)
		require.Equal(t, c.start, start.String())
		require.Equal(t, c.end, end.String())
	}
}
