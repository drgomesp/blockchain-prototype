package rpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/rhizomplatform/rhizom/proto/gen/service"
	"google.golang.org/grpc"
)

// Server type uses the gRPC and the generated resources from protobuff
// to receive requests and communicate through the internet.
type Server struct {
	service.NodeServer
	name string
	grpc *grpc.Server
}

// NewServer instantiates a new Server object.
func NewServer(name string) *Server {
	return &Server{
		NodeServer: service.UnimplementedNodeServer{},
		name:       name,
		grpc:       grpc.NewServer(),
	}
}

// Name of this server instance.
func (s *Server) Name() string { return s.name }

// Start this server instance.
func (s *Server) Start(port uint) error {
	service.RegisterNodeServer(s.grpc, s.NodeServer)
	return s.grpc.Serve(newNetListener(port))
}

// Stop this server with a timeout duration.
// The server will try to stop it gracefully, so, it will wait for all the
// pending RPCs to be finished. However, if the timout is reached, it will make
// a forced stop. If you want that not being forced, give a timeout value of -1.
// But, if you want to immediately stop, give a timout of 0.
func (s *Server) Stop(ctx context.Context, timeout time.Duration) error {
	ctxTO, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	tick := time.NewTicker(timeout)
	defer tick.Stop()

	go func() {
		s.grpc.GracefulStop()
		ctxTO.Done()
	}()

	select {
	case <-ctxTO.Done():
	case <-tick.C:
		s.grpc.Stop()
	}

	return ctx.Err()
}

// newNetListener instantiates a new new.NetListener with
// the provided port to be exposed for TCP connections.
func newNetListener(port uint) net.Listener {
	tcp, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	return tcp
}
