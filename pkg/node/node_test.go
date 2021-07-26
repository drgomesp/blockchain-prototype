package node_test

import (
	"testing"

	"github.com/rhizomplatform/rhizom/pkg/node"
	"github.com/rhizomplatform/rhizom/pkg/rpc"
	"github.com/stretchr/testify/assert"
)

// Test node constructor.
// nolint:paralleltest
func TestNodeNew(t *testing.T) {
	n, err := node.New(&node.Config{
		Type: node.TypeFull,
		Name: "test node",
	})

	assert.NoError(t, err)
	assert.NotNilf(t, n, "node should not be nil")
}

// Test node register apis.
func TestNode_RegisterAPIs(t *testing.T) {
	n, _ := node.New(&node.Config{
		Type: node.TypeFull,
		Name: "test node",
	})

	a1, a2 := &rpc.API{}, &rpc.API{}
	n.RegisterAPIs(a1, a2)

	assert.Containsf(t, n.APIs(), a1, "test node should contain a1 api")
	assert.Containsf(t, n.APIs(), a2, "test node should contain a2 api")
}

// Test node register server.
func TestNode_RegisterServer(t *testing.T) {
	n, _ := node.New(&node.Config{
		Type: node.TypeFull,
		Name: "test node",
	})

	s1, s2 := new(fakeServer), new(fakeServer)
	n.RegisterServers(s1, s2)

	assert.Containsf(t, n.Servers(), s1, "test node should contain s1 server")
	assert.Containsf(t, n.Servers(), s2, "test node should contain s2 server")
}

type fakeServer struct {
	Registered bool
}

func (s *fakeServer) Start() error {
	s.Registered = true
	return nil
}

func (s *fakeServer) Stop() error {
	s.Registered = false
	return nil
}
