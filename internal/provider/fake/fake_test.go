package fake

import (
	"context"
	"errors"
	"testing"

	"github.com/Ploos-AS/Ploos-VPS/internal/provider"
)

func TestLifecycle(t *testing.T) {
	ctx := context.Background()
	p := New()
	req := provider.CreateRequest{ID: "vm1", Node: "lab", Image: "debian", VCPU: 1, RAMMiB: 512}
	vm, err := p.Create(ctx, req)
	if err != nil || vm.State != "stopped" {
		t.Fatalf("create: %v %v", vm, err)
	}
	retry, err := p.Create(ctx, req)
	if err != nil || retry != vm {
		t.Fatalf("retry: %v %v", retry, err)
	}
	changed := req
	changed.VCPU++
	if _, err := p.Create(ctx, changed); !errors.Is(err, ErrConflict) {
		t.Fatalf("conflict: %v", err)
	}
	if err := p.Start(ctx, req.ID); err != nil {
		t.Fatal(err)
	}
	vm, err = p.Get(ctx, req.ID)
	if err != nil || vm.State != "running" {
		t.Fatalf("start: %v %v", vm, err)
	}
	if err := p.Stop(ctx, req.ID); err != nil {
		t.Fatal(err)
	}
	vm, err = p.Get(ctx, req.ID)
	if err != nil || vm.State != "stopped" {
		t.Fatalf("stop: %v %v", vm, err)
	}
	for range 2 {
		if err := p.Delete(ctx, req.ID); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := p.Get(ctx, req.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("deleted: %v", err)
	}
}

func TestInvalidAndCancelled(t *testing.T) {
	p := New()
	if _, err := p.Create(context.Background(), provider.CreateRequest{}); !errors.Is(err, ErrInvalid) {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := p.Create(ctx, provider.CreateRequest{}); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
