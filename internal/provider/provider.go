package provider

import "context"

type Node struct {
	Name   string
	Status string
}

type Inventory interface {
	ListNodes(ctx context.Context) ([]Node, error)
}
