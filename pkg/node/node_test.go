package node_test

import (
	"context"
	"testing"

	"github.com/drgomesp/rhizom/pkg/node"
	"github.com/drgomesp/rhizom/pkg/rpc"
	"github.com/stretchr/testify/assert"
)

// Test node constructor.
// nolint:paralleltest
func TestNodeNew(t *testing.T) {
	n, err := node.New(node.Config{
		Type: node.TypeFull,
		Name: "test node",
	})

	assert.NoError(t, err)
	assert.NotNilf(t, n, "node should not be nil")
}

// Test node register apis.
func TestNode_RegisterAPIs(t *testing.T) {
	t.Parallel()

	n, _ := node.New(node.Config{
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
	t.Parallel()

	n, _ := node.New(node.Config{
		Type: node.TypeFull,
		Name: "test node",
	})

	s1, s2 := new(fakeServer), new(fakeServer)
	n.RegisterProtocols(s1, s2)

	assert.Containsf(t, n.Servers(), s1, "test node should contain s1 server")
	assert.Containsf(t, n.Servers(), s2, "test node should contain s2 server")
}

type fakeServer struct{}

func (s *fakeServer) Name() string {
	return "fake"
}

func (s *fakeServer) Start(ctx context.Context) error {
	return nil
}

func (s *fakeServer) Stop(ctx context.Context) error {
	return nil
}
