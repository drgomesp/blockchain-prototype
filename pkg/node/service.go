package node

import "context"

// Service defines a service that can be registered as part of the node lifecycle.
type Service interface {
	// Name of the service.
	Name() string
	// Start the service.
	Start(ctx context.Context) error
	// Stop the service.
	Stop(ctx context.Context) error
}
