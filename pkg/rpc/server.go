package rpc

import (
	"fmt"
	"io"
	"net"

	"github.com/drgomesp/rhizom/proto/gen/entity"
	"github.com/drgomesp/rhizom/proto/gen/message"
	"github.com/drgomesp/rhizom/proto/gen/service"
	"google.golang.org/grpc"
)

// Server type uses the gRPC and the generated resources from protobuff
// to receive requests and communicate through the internet.
type Server struct {
	name string
	grpc *grpc.Server
	service.NodeServer
}

// NewServer instantiates a new Server object.
func NewServer(name string) *Server {
	return &Server{
		name:       name,
		grpc:       grpc.NewServer(),
		NodeServer: service.UnimplementedNodeServer{},
	}
}

// Name of this server instance.
func (s *Server) Name() string { return s.name }

// Info about this server instance.
func (s *Server) Info() map[string]grpc.ServiceInfo { return s.grpc.GetServiceInfo() }

// Start this server instance.
func (s *Server) Start(listener Listener) error {
	service.RegisterNodeServer(s.grpc, s)
	return s.grpc.Serve(listener)
}

// Stop this server open connections.
func (s *Server) Stop() { s.grpc.Stop() }

// GetBlocks calls the core service GetBlocks method and maps the result to a grpc service response.
func (s *Server) GetBlock(stream service.Node_GetBlockServer) error {
	for {
		var resp *message.GetBlockResponse

		switch req, err := stream.Recv(); err {
		case nil:
			resp = &message.GetBlockResponse{
				Block: &entity.Block{Index: req.Want},
			}

		case io.EOF:
			return nil

		default:
			resp = &message.GetBlockResponse{Err: err.Error()}
		}

		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}

type Listener net.Listener

// NewListener instantiates a new new.NetListener with
// the provided port to be exposed for TCP connections.
func NewListener(port int) Listener {
	tcp, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	return Listener(tcp)
}
