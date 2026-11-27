package fake

import (
	"context"
	"errors"
	"sync"

	"github.com/Ploos-AS/Ploos-VPS/internal/provider"
)

var (
	ErrNotFound = errors.New("instance not found")
	ErrConflict = errors.New("instance request conflict")
	ErrInvalid  = errors.New("invalid instance request")
)

type entry struct {
	request  provider.CreateRequest
	instance provider.Instance
}

type Provider struct {
	mu        sync.Mutex
	instances map[string]entry
}

var _ provider.Lifecycle = (*Provider)(nil)

func New() *Provider {
	return &Provider{instances: make(map[string]entry)}
}

func (p *Provider) Create(ctx context.Context, req provider.CreateRequest) (provider.Instance, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return provider.Instance{}, err
	}
	if req.ID == "" || req.Node == "" || req.Image == "" || req.VCPU <= 0 || req.RAMMiB <= 0 {
		return provider.Instance{}, ErrInvalid
	}
	if old, ok := p.instances[req.ID]; ok {
		if old.request != req {
			return provider.Instance{}, ErrConflict
		}
		return old.instance, nil
	}
	vm := provider.Instance{ID: req.ID, Node: req.Node, State: "stopped"}
	p.instances[req.ID] = entry{request: req, instance: vm}
	return vm, nil
}

func (p *Provider) Get(ctx context.Context, id string) (provider.Instance, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return provider.Instance{}, err
	}
	e, ok := p.instances[id]
	if !ok {
		return provider.Instance{}, ErrNotFound
	}
	return e.instance, nil
}

func (p *Provider) Start(ctx context.Context, id string) error {
	return p.setState(ctx, id, "running")
}

func (p *Provider) Stop(ctx context.Context, id string) error {
	return p.setState(ctx, id, "stopped")
}

func (p *Provider) setState(ctx context.Context, id, state string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	e, ok := p.instances[id]
	if !ok {
		return ErrNotFound
	}
	e.instance.State = state
	p.instances[id] = e
	return nil
}

func (p *Provider) Delete(ctx context.Context, id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	delete(p.instances, id)
	return nil
}
