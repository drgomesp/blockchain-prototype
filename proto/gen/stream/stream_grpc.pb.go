// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package stream

import (
	context "context"
	message "github.com/drgomesp/rhizom/proto/gen/message"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// BlockClient is the client API for Block service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BlockClient interface {
	GetBlock(ctx context.Context, opts ...grpc.CallOption) (Block_GetBlockClient, error)
}

type blockClient struct {
	cc grpc.ClientConnInterface
}

func NewBlockClient(cc grpc.ClientConnInterface) BlockClient {
	return &blockClient{cc}
}

func (c *blockClient) GetBlock(ctx context.Context, opts ...grpc.CallOption) (Block_GetBlockClient, error) {
	stream, err := c.cc.NewStream(ctx, &Block_ServiceDesc.Streams[0], "/stream.Block/GetBlock", opts...)
	if err != nil {
		return nil, err
	}
	x := &blockGetBlockClient{stream}
	return x, nil
}

type Block_GetBlockClient interface {
	Send(*message.GetBlockRequest) error
	Recv() (*message.GetBlockResponse, error)
	grpc.ClientStream
}

type blockGetBlockClient struct {
	grpc.ClientStream
}

func (x *blockGetBlockClient) Send(m *message.GetBlockRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *blockGetBlockClient) Recv() (*message.GetBlockResponse, error) {
	m := new(message.GetBlockResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BlockServer is the server API for Block service.
// All implementations must embed UnimplementedBlockServer
// for forward compatibility
type BlockServer interface {
	GetBlock(Block_GetBlockServer) error
	mustEmbedUnimplementedBlockServer()
}

// UnimplementedBlockServer must be embedded to have forward compatible implementations.
type UnimplementedBlockServer struct {
}

func (UnimplementedBlockServer) GetBlock(Block_GetBlockServer) error {
	return status.Errorf(codes.Unimplemented, "method GetBlock not implemented")
}
func (UnimplementedBlockServer) mustEmbedUnimplementedBlockServer() {}

// UnsafeBlockServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BlockServer will
// result in compilation errors.
type UnsafeBlockServer interface {
	mustEmbedUnimplementedBlockServer()
}

func RegisterBlockServer(s grpc.ServiceRegistrar, srv BlockServer) {
	s.RegisterService(&Block_ServiceDesc, srv)
}

func _Block_GetBlock_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BlockServer).GetBlock(&blockGetBlockServer{stream})
}

type Block_GetBlockServer interface {
	Send(*message.GetBlockResponse) error
	Recv() (*message.GetBlockRequest, error)
	grpc.ServerStream
}

type blockGetBlockServer struct {
	grpc.ServerStream
}

func (x *blockGetBlockServer) Send(m *message.GetBlockResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *blockGetBlockServer) Recv() (*message.GetBlockRequest, error) {
	m := new(message.GetBlockRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Block_ServiceDesc is the grpc.ServiceDesc for Block service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Block_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "stream.Block",
	HandlerType: (*BlockServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetBlock",
			Handler:       _Block_GetBlock_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/stream.proto",
}
