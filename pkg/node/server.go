package node

import "context"

// Server defines a server that can be registered to the node.
type Server interface {
	Name() string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
