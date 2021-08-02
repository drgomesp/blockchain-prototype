package node_test

import (
	"context"
	"testing"

	"github.com/drgomesp/rhizom/pkg/node"
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

// Test node register services.
func TestNode_RegisterServices(t *testing.T) {
	t.Parallel()

	n, _ := node.New(node.Config{
		Type: node.TypeFull,
		Name: "test node",
	})

	s1, s2 := new(fakeService), new(fakeService)
	n.RegisterServices(s1, s2)

	assert.Containsf(t, n.Services(), s1, "test node should contain s1 service")
	assert.Containsf(t, n.Services(), s2, "test node should contain s2 service")
}

type fakeService struct{}

func (s *fakeService) Name() string {
	return "fake"
}

func (s *fakeService) Start(ctx context.Context) error {
	return nil
}

func (s *fakeService) Stop(ctx context.Context) error {
	return nil
}
