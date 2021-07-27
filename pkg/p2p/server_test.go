package p2p_test

import (
	"context"
	"testing"

	"github.com/rhizomplatform/rhizom/pkg/p2p"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestServer_Start(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t).Sugar()
	ctx := context.Background()
	srv, err := p2p.NewServer(ctx, logger, p2p.Config{MaxPeers: 3})

	assert.NoError(t, err)
	assert.NotNil(t, srv)
}
