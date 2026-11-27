package provider

import "context"

type CreateRequest struct {
	ID     string
	Node   string
	Image  string
	VCPU   int
	RAMMiB int
}

type Instance struct {
	ID    string
	Node  string
	State string
}

type Lifecycle interface {
	Create(context.Context, CreateRequest) (Instance, error)
	Get(context.Context, string) (Instance, error)
	Start(context.Context, string) error
	Stop(context.Context, string) error
	Delete(context.Context, string) error
}
