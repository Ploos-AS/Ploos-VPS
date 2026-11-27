package ipam

import (
	"errors"
	"net/netip"
	"testing"
)

func TestNext(t *testing.T) {
	for _, tc := range []struct {
		prefix string
		used   string
		want   string
	}{
		{"192.0.2.0/29", "192.0.2.1", "192.0.2.2"},
		{"2001:db8::/64", "2001:db8::1", "2001:db8::2"},
	} {
		got, err := Next(netip.MustParsePrefix(tc.prefix), map[netip.Addr]bool{netip.MustParseAddr(tc.used): true})
		if err != nil || got.String() != tc.want {
			t.Fatalf("got %v, %v; want %s", got, err, tc.want)
		}
	}
}

func TestExhaustion(t *testing.T) {
	p := netip.MustParsePrefix("192.0.2.0/30")
	used := map[netip.Addr]bool{
		netip.MustParseAddr("192.0.2.1"): true,
		netip.MustParseAddr("192.0.2.2"): true,
	}
	if _, err := Next(p, used); !errors.Is(err, ErrExhausted) {
		t.Fatalf("expected exhaustion: %v", err)
	}
	if _, err := Next(netip.Prefix{}, nil); err == nil {
		t.Fatal("expected invalid prefix error")
	}
}
