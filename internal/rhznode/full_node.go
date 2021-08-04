package rhznode

import (
	"context"
	"time"

	"github.com/drgomesp/rhizom/internal/rhz"
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
	rhzPrefix = "/rhz/"
	devNet    = "default_2b678c95-27d5-4f09-bf38-a62be2c5339b"
	testNet   = "rhz_testnet_e19d2c16-8c39-4f0f-8c88-2427e37c12bb"
	net       = devNet

	TopicBlocks       = rhzPrefix + "blk/" + net
	TopicProducers    = rhzPrefix + "prc/" + net
	TopicTransactions = rhzPrefix + "tx/" + net
	TopicRequestSync  = rhzPrefix + "blkchain/req/" + net

	ProtocolRequestBlocks     = rhzPrefix + "blocks/req/" + net
	ProtocolResponseBlocks    = rhzPrefix + "blocks/resp/" + net
	ProtocolRequestDelegates  = rhzPrefix + "delegates/req/" + net
	ProtocolResponseDelegates = rhzPrefix + "delegates/resp/" + net

	p2pServerMaxPeers    = 5
	p2pServerPingTimeout = time.Second * 5
)

var topics = []string{
	TopicBlocks,
	// TopicProducers,
	// TopicTransactions,
	// TopicRequestSync,
	ProtocolRequestBlocks,
	ProtocolResponseBlocks,
}

type FullNode struct {
	node         *node.Node
	logger       *zap.SugaredLogger
	peerExchange rhz.PeerExchange

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
		node:         n,
		logger:       logger,
		p2pServer:    n.Server(),
		peerExchange: rhz.NewHandler(logger),
	}

	n.RegisterAPIs(nil)
	n.RegisterProtocols(fullNode.Protocols(fullNode.peerExchange)...)
	n.RegisterServices(fullNode)

	return fullNode, nil
}

func (n *FullNode) Start(ctx context.Context) error {
	var err error
	n.logger.Infof("starting full node")

	if err = n.p2pServer.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start p2p server")
	}

	n.p2pServer.RegisterProtocols(n.Protocols(n.peerExchange)...)

	for {
		select {
		case <-ctx.Done():
			return n.Stop(ctx)
		default:
			{
				time.Sleep(time.Second * 5)

				req := rhz.GetBlocksRequest{IndexHave: 0, IndexNeed: 10}
				if err := n.p2pServer.StreamMsg(
					ctx,
					ProtocolRequestBlocks,
					req,
				); err != nil {
					n.logger.Error(err)
				}

				n.logger.Debugw("request sent", "req", req)
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

func (n *FullNode) Protocols(backend rhz.PeerExchange) []p2p.Protocol {
	return []p2p.Protocol{
		{
			ID:  ProtocolRequestBlocks,
			Run: n.requestHandler(backend),
		},
		{
			ID:  ProtocolResponseBlocks,
			Run: n.responseHandler(backend),
		},
		{
			ID:  "rhz_test",
			Run: n.requestHandler(backend),
		},
	}
}

func (n *FullNode) requestHandler(backend rhz.PeerExchange) func(ctx context.Context, rw p2p.MsgReadWriter) error {
	return func(ctx context.Context, rw p2p.MsgReadWriter) error {
		peer := rhz.NewPeer(rw)
		return rhz.HandleRequestMsg(ctx, backend, peer)
	}
}

func (n *FullNode) responseHandler(backend rhz.PeerExchange) func(ctx context.Context, rw p2p.MsgReadWriter) error {
	return func(ctx context.Context, rw p2p.MsgReadWriter) error {
		peer := rhz.NewPeer(rw)
		return rhz.HandleResponseMsg(ctx, backend, peer)
	}
}
