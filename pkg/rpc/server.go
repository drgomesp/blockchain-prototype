package rpc

import (
	"github.com/drgomesp/rhizom/proto/gen/service"
	"google.golang.org/grpc"
)

// Server type uses the gRPC and the generated resources from protobuff
// to receive requests and communicate through the internet.
type Server struct {
	name         string
	blockService *BlockService
	grpc         *grpc.Server
}

// NewServer instantiates a new Server object.
func NewServer(name string, blockService *BlockService) *Server {
	return &Server{
		name:         name,
		blockService: blockService,
		grpc:         grpc.NewServer(),
	}
}

// Name of this server instance.
func (s *Server) Name() string { return s.name }

// Info about this server instance.
func (s *Server) Info() map[string]grpc.ServiceInfo { return s.grpc.GetServiceInfo() }

// Start this server instance.
func (s *Server) Start(listener Listener) error {
	service.RegisterNodeServer(s.grpc, s.blockService)
	return s.grpc.Serve(listener)
}

// Stop this server open connections.
func (s *Server) Stop() { s.grpc.Stop() }
