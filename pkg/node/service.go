package node

import "context"

// Service defines a service that can be managed by the node as part of its lifecycle.
type Service interface {
	// Name of the service.
	Name() string
	// Start the service.
	Start(ctx context.Context) error
	// Stop the service.
	Stop(ctx context.Context) error
}
