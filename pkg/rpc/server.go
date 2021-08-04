package rpc

import (
<<<<<<< Updated upstream
	"github.com/drgomesp/rhizom/proto/gen/service"
	"google.golang.org/grpc"
)

// Server type uses the gRPC and the generated resources from protobuff
// to receive requests and communicate through the internet.
type Server struct {
	name         string
	blockService *BlockService
	grpc         *grpc.Server
=======
<<<<<<< Updated upstream
	"context"
)

type Server struct{}

func NewServer() (*Server, error) {
	return &Server{}, nil
=======
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
>>>>>>> Stashed changes
>>>>>>> Stashed changes
}

// NewServer instantiates a new Server object.
func NewServer(name string, blockService *BlockService) *Server {
	return &Server{
		name:         name,
		blockService: blockService,
		grpc:         grpc.NewServer(),
	}
}

<<<<<<< Updated upstream
// Name of this server instance.
func (s *Server) Name() string { return s.name }
=======
<<<<<<< Updated upstream
func (s *Server) Start(_ context.Context) error {
	return nil
}
=======
// Start this server instance.
func (s *Server) Start(listener Listener) error { return s.grpc.Serve(listener) }
>>>>>>> Stashed changes
>>>>>>> Stashed changes

// Info about this server instance.
func (s *Server) Info() map[string]grpc.ServiceInfo { return s.grpc.GetServiceInfo() }

// Start this server instance.
func (s *Server) Start(listener Listener) error {
	service.RegisterNodeServer(s.grpc, s.blockService)
	return s.grpc.Serve(listener)
}

// Stop this server open connections.
func (s *Server) Stop() { s.grpc.Stop() }
