package rpc

import (
	"google.golang.org/grpc"
)

// Server type uses the gRPC and the generated resources from protobuff
// to receive requests and communicate through the internet.
type Server struct {
	name string
	grpc *grpc.Server
}

// NewServer instantiates a new Server object.
func NewServer(name string, grcpServer *grpc.Server) *Server {
	return &Server{
		name: name,
		grpc: grcpServer,
	}
}

// Name of this server instance.
func (s *Server) Name() string { return s.name }

// Start this server instance.
func (s *Server) Start(listener Listener) error { return s.grpc.Serve(listener) }

// Info about this server instance.
func (s *Server) Info() map[string]grpc.ServiceInfo { return s.grpc.GetServiceInfo() }

// Stop this server open connections.
func (s *Server) Stop() { s.grpc.Stop() }
