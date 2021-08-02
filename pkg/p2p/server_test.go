package p2p_test

import (
	"testing"

	"github.com/drgomesp/rhizom/pkg/p2p"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestServer_Start(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t).Sugar()
	srv, err := p2p.NewServer(p2p.Config{MaxPeers: 3}, p2p.WithLogger(logger))

	assert.NoError(t, err)
	assert.NotNil(t, srv)
}
