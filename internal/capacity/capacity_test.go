package capacity

import (
	"errors"
	"testing"
)

func TestSelect(t *testing.T) {
	a := Node{ID: "a", Enabled: true, Total: Resources{8, 8192, 100}, Reserved: Resources{2, 2048, 20}, Used: Resources{2, 2048, 40}}
	b := a
	b.ID = "b"
	b.Used.RAMMiB = 0
	disabled := b
	disabled.ID = "disabled"
	disabled.Enabled = false
	req := Resources{2, 4096, 40}
	got, err := Select([]Node{a, disabled, b}, req)
	if err != nil || got.ID != "b" {
		t.Fatalf("selection: %v %v", got, err)
	}
	if _, err := Select([]Node{a, b}, Resources{2, 4096, 41}); !errors.Is(err, ErrNoCapacity) {
		t.Fatalf("expected disk rejection: %v", err)
	}
	if a.Fits(Resources{0, 1, 1}) {
		t.Fatal("accepted invalid request")
	}
	a.Used.RAMMiB = -1
	if a.Fits(req) {
		t.Fatal("accepted invalid inventory")
	}
}

func TestTieAndOverCapacity(t *testing.T) {
	a := Node{ID: "a", Enabled: true, Total: Resources{4, 4096, 40}}
	b := a
	b.ID = "b"
	got, err := Select([]Node{b, a}, Resources{1, 512, 10})
	if err != nil || got.ID != "a" {
		t.Fatalf("tie: %v %v", got, err)
	}
	a.Reserved.RAMMiB = 8192
	if a.Available().RAMMiB != 0 {
		t.Fatal("negative available capacity")
	}
}
