package node

import "context"

// Server defines a server that can be registered to the node.
type Server interface {
	// Name of the server.
	Name() string

	// Start the server.
	Start(ctx context.Context) error

	// Stop the server.
	Stop(ctx context.Context) error
}
