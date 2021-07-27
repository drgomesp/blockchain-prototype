package rhznode

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rhizomplatform/rhizom/pkg/node"
	"github.com/rhizomplatform/rhizom/pkg/p2p"
	"github.com/rhizomplatform/rhizom/pkg/rpc"
	"go.uber.org/zap"
)

type FullNode struct {
	node   *node.Node
	logger *zap.SugaredLogger
}

func NewFullNode(logger *zap.SugaredLogger, node *node.Node) (*FullNode, error) {
	rhz := &FullNode{
		node:   node,
		logger: logger,
	}

	const maxPeers = 5

	srv, err := p2p.NewServer(
		context.Background(),
		logger, p2p.Config{MaxPeers: maxPeers},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize p2p server")
	}

	node.RegisterAPIs(rhz.APIs()...)
	node.RegisterServers(srv)

	return rhz, nil
}

func (n *FullNode) Start(ctx context.Context) error {
	var started []node.Server

	for _, srv := range n.node.Servers() {
		if err := srv.Start(ctx); err != nil {
			break
		}

		started = append(started, srv)
	}

	for {
		n.logger.With("servers", started).Info("running")
		time.Sleep(time.Second)
	}
}

func (n *FullNode) Stop(ctx context.Context) error {
	panic("implement me")
}

func (n *FullNode) APIs() []*rpc.API {
	return []*rpc.API{}
}
