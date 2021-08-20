package rpc

import (
	"context"
	"testing"

	"github.com/drgomesp/rhizom/proto/gen/message"
	"github.com/drgomesp/rhizom/proto/gen/stream"
	th "github.com/drgomesp/rhizom/test/testhelper"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/nettest"
	"google.golang.org/grpc"
)

func Test_NewBlockStream(t *testing.T) {
	s := NewBlockStream()
	assert.NotNil(t, s)
}

func TestBlockStream(t *testing.T) {
	// nettest local network listener
	netT := th.LocalListener(t)
	defer func() { assert.NoError(t, netT.Close()) }()

	// server setup
	setup := grpc.NewServer()
	stream.RegisterBlockServer(setup, NewBlockStream())
	s := NewServer("test", setup)

	// run Server
	go func() {
		defer s.Stop()
		s.Start(netT)
	}()

	// test local net
	addr := netT.Addr()
	assert.True(t, nettest.TestableAddress(addr.Network(), addr.String()))
	assert.True(t, nettest.TestableNetwork(netT.Addr().Network()))

	// grpc dial
	conn, err := grpc.Dial(addr.String(), grpc.WithInsecure())
	th.FailOnError(t, err)
	defer func() { assert.NoError(t, conn.Close()) }()

	// client stream
	stream, err := stream.NewBlockClient(conn).GetBlock(context.Background())
	th.FailOnError(t, err)
	defer func() { assert.NoError(t, stream.CloseSend()) }()

	// send request
	req := &message.RequestStreamGetBlock{IndexWant: 1}
	err = stream.Send(req)
	th.FailOnError(t, err)

	// receive response
	resp, err := stream.Recv()
	th.FailOnError(t, err)
	assert.Equal(t, req.IndexWant, resp.Block.Index)
}
