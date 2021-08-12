package rpc

import (
	"google.golang.org/grpc"
)

// API type uses the gRPC and the generated resources from protobuff
// to receive requests and communicate through the internet.
type API struct {
	name string
	grpc *grpc.Server
}

// NewServer instantiates a new API object.
func NewServer(name string, grcpServer *grpc.Server) *API {
	return &API{
		name: name,
		grpc: grcpServer,
	}
}

// Name of this server instance.
func (s *API) Name() string { return s.name }

// Start this server instance.
func (s *API) Start(listener Listener) error { return s.grpc.Serve(listener) }

// Info about this server instance.
func (s *API) Info() map[string]grpc.ServiceInfo { return s.grpc.GetServiceInfo() }

// Stop this server open connections.
func (s *API) Stop() { s.grpc.Stop() }
