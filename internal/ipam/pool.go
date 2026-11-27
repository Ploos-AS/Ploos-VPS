package ipam

import (
	"errors"
	"net/netip"
)

var ErrExhausted = errors.New("address pool exhausted")

// Next returns the first unoccupied address within the prefix.
// The first address is reserved; IPv4 broadcast is also excluded.
// Callers must reserve the result atomically before assigning it.
func Next(prefix netip.Prefix, occupied map[netip.Addr]bool) (netip.Addr, error) {
	if !prefix.IsValid() || prefix.Addr().Is4In6() {
		return netip.Addr{}, errors.New("invalid address prefix")
	}
	prefix = prefix.Masked()
	for addr := prefix.Addr().Next(); addr.IsValid() && prefix.Contains(addr); addr = addr.Next() {
		if addr.Is4() && !prefix.Contains(addr.Next()) {
			break
		}
		if !occupied[addr] {
			return addr, nil
		}
	}
	return netip.Addr{}, ErrExhausted
}
