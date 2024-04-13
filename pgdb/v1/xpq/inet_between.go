package xpq

import (
	"net/netip"

	"github.com/gaissmai/extnetip"
)

func NetworkRange(p netip.Prefix) (netip.Addr, netip.Addr) {
	return extnetip.Range(p)
}
