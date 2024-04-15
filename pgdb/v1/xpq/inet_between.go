package xpq

import (
	"net/netip"

	"github.com/gaissmai/extnetip"
)

func NetworkRange(p netip.Prefix) (netip.Addr, netip.Addr) {
	start, end := extnetip.Range(p)
	return start.Unmap(), end.Unmap()
}
