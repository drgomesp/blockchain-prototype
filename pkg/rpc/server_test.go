package rpc

import (
	"testing"

	th "github.com/drgomesp/rhizom/test/testhelper"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/nettest"
	"google.golang.org/grpc"
)

func Test_NewServer(t *testing.T) {
	s := NewServer("test", grpc.NewServer())
	assert.NotNil(t, s)
}

func TestServer_Name(t *testing.T) {
	name := "test"
	s := NewServer("test", grpc.NewServer())
	assert.Equal(t, name, s.Name())
}

func TestServer_Info(t *testing.T) {
	s := NewServer("test", grpc.NewServer())
	assert.IsType(t, map[string]grpc.ServiceInfo{}, s.Info())
}

func TestServer(t *testing.T) {
	// nettest local network listener
	netT := th.LocalListener(t)
	defer func() { assert.NoError(t, netT.Close()) }()

	// start Server
	s := NewServer("test", grpc.NewServer())
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
}
