package node_test

import (
	"testing"

	"github.com/rhizomplatform/rhizom/pkg/node"
	"github.com/stretchr/testify/assert"
)

// nolint:paralleltest
// Test node constructor.
func TestNodeNew(t *testing.T) {
	n, err := node.New(&node.Config{
		Type: node.TypeFull,
		Name: "test node",
	})

	assert.NoError(t, err)
	assert.NotNilf(t, n, "node should not be nil")
}
