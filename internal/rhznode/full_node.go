package rhznode

import (
	"context"
	"time"

	"github.com/drgomesp/rhizom/pkg/node"
	"github.com/drgomesp/rhizom/pkg/p2p"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var bootstrapAddrs = []string{
	"/dns4/bootstrapper-1.rhz.network/tcp/4001/ipfs/Qmf8Lt1FiQnG7tLrQbhwvUXzBMYsj6KicNdKiD1F2rSRW5",
	"/dns4/bootstrapper-2.rhz.network/tcp/4001/ipfs/QmcRoi1mQ7eb7xPDhWZjGL8rivAUHwCv1FMiLw7FGSZvFL",
}

const (
	rhzPrefix         = "/rhz/"
	net               = "default_2b678c95-27d5-4f09-bf38-a62be2c5339b"
	TopicBlocks       = "/rhz/blk/" + net
	TopicProducers    = "/rhz/prc/" + net
	TopicTransactions = "/rhz/tx/" + net
	TopicRequestSync  = "/rhz/blkchain/req/" + net

	ProtocolRequestBlocks     = rhzPrefix + "blocks/req/" + net
	ProtocolResponseBlocks    = rhzPrefix + "blocks/resp/" + net
	ProtocolRequestDelegates  = rhzPrefix + "delegates/req/" + net
	ProtocolResponseDelegates = rhzPrefix + "delegates/resp/" + net

	p2pServerMaxPeers    = 5
	p2pServerPingTimeout = time.Second * 5
)

var topics = []string{
	TopicBlocks,
	TopicProducers,
	TopicTransactions,
	TopicRequestSync,
}

type FullNode struct {
	node   *node.Node
	logger *zap.SugaredLogger

	p2pServer *p2p.Server
}

func NewFullNode(logger *zap.SugaredLogger) (*FullNode, error) {
	n, err := node.New(node.Config{
		Type: node.TypeFull,
		Name: "rhz_node",
		P2P: p2p.Config{
			MaxPeers:       p2pServerMaxPeers,
			PingTimeout:    p2pServerPingTimeout,
			BootstrapAddrs: bootstrapAddrs,
			Topics:         topics,
		},
	}, node.WithLogger(logger))
	if err != nil {
		return nil, errors.Wrap(err, "failed initialize node")
	}

	fullNode := &FullNode{
		node:      n,
		logger:    logger,
		p2pServer: n.Server(),
	}

	n.RegisterAPIs(nil)
	n.RegisterProtocols(fullNode.Protocols()...)
	n.RegisterServices(fullNode)

	return fullNode, nil
}

func (n *FullNode) Start(ctx context.Context) error {
	n.logger.Infof("starting full node")

	if err := n.p2pServer.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start p2p server")
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
	n.logger.Infow("stopping full node")

	return nil
}

func (n *FullNode) Name() string {
	return "full_node"
}

func (n *FullNode) Protocols() []p2p.Protocol {
	return []p2p.Protocol{}
}
