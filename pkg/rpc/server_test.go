package rpc

import (
	"context"
	"testing"

	"github.com/drgomesp/rhizom/proto/gen/message"
	"github.com/drgomesp/rhizom/proto/gen/service"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/nettest"
	"google.golang.org/grpc"
)

func Test_NewListener(t *testing.T)      { testHelper_NewListener(t) }
func Test_NewStreamService(t *testing.T) { testHelper_NewListener(t) }
func Test_NewServer(t *testing.T)        { testHelper_NewServer(t, "test", testHelper_NewStreamService(t)) }

func TestServer_Name(t *testing.T) {
	name := "test"
	s := testHelper_NewServer(t, name, testHelper_NewStreamService(t))
	assert.Equal(t, name, s.Name())
}

func TestServer_Info(t *testing.T) {
	s := testHelper_NewServer(t, "test", testHelper_NewStreamService(t))
	assert.IsType(t, map[string]grpc.ServiceInfo{}, s.Info())
}

func TestServer_Start(t *testing.T) {
	// nettest local network listener
	netT := testHelper_localListener(t)
	defer func() { assert.NoError(t, netT.Close()) }()

	// run Server
	s := testHelper_NewServer(t, "test", testHelper_NewStreamService(t))
	go func() {
		defer s.Stop()
		s.Start(netT)
	}()

	// local net test
	addr := netT.Addr()
	assert.True(t, nettest.TestableAddress(addr.Network(), addr.String()))
	assert.True(t, nettest.TestableNetwork(netT.Addr().Network()))

	// grpc dial
	conn, err := grpc.Dial(addr.String(), grpc.WithInsecure())
	assert.NoError(t, err)

	// client stream
	stream, err := service.NewNodeClient(conn).GetBlock(context.Background())
	assert.NoError(t, err)
	defer stream.CloseSend()

	// send request
	req := &message.GetBlockRequest{Have: 0, Want: 1}
	err = stream.Send(req)
	assert.NoError(t, err)

	// receive response
	resp, err := stream.Recv()
	assert.NoError(t, err)
	assert.Equal(t, req.Want, resp.Block.Index)
}
