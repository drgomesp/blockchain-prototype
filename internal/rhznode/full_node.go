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

var BootstrapAddrs = []string{
	"/dns4/bootstrapper-1.rhz.network/tcp/4001/ipfs/Qmf8Lt1FiQnG7tLrQbhwvUXzBMYsj6KicNdKiD1F2rSRW5",
	"/dns4/bootstrapper-2.rhz.network/tcp/4001/ipfs/QmcRoi1mQ7eb7xPDhWZjGL8rivAUHwCv1FMiLw7FGSZvFL",
}

type FullNode struct {
	node   *node.Node
	logger *zap.SugaredLogger
}

func NewFullNode(ctx context.Context, logger *zap.SugaredLogger) (*FullNode, error) {
	n, err := node.New(node.Config{
		Type: node.TypeFull,
		Name: "rhz_node",
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed initialize node")
	}

	const maxPeers = 5

	p2pServer, err := p2p.NewServer(ctx, logger, p2p.Config{
		MaxPeers:       maxPeers,
		PingTimeout:    time.Second * 5,
		BootstrapAddrs: BootstrapAddrs,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize p2p server")
	}

	rpcServer, err := rpc.NewServer()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize rpc server")
	}

	rhz := &FullNode{
		node:   n,
		logger: logger,
	}

	n.RegisterAPIs(rhz.APIs()...)
	n.RegisterServers(p2pServer, rpcServer)

	return rhz, nil
}

func (n *FullNode) Start(ctx context.Context) error {
	for _, srv := range n.node.Servers() {
		if err := srv.Start(ctx); err != nil {
			break
		}

		n.logger.With("name", srv.Name()).Info("server started")
	}

	for {
		select {
		case <-ctx.Done():
			return n.Stop(ctx)
		default:
			{
			}
		}
	}
}

func (n *FullNode) Stop(_ context.Context) error {
	n.logger.Infow("stopping node", "name", n.node.Config().Name)

	return nil
}

func (n *FullNode) APIs() []*rpc.API {
	return []*rpc.API{}
}
