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
	p2pServerMaxPeers    = 5
	p2pServerPingTimeout = time.Second * 5
)

var topics = []string{
	rhz.TopicBlocks,
	rhz.TopicProducers,
	rhz.TopicTransactions,
	rhz.TopicRequestSync,
	rhz.ProtocolRequestBlocks,
	rhz.ProtocolResponseBlocks,
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
		peerExchange: NewHandler(logger),
	}

	n.RegisterAPIs(nil)
	n.RegisterProtocols(fullNode.Protocols(fullNode.peerExchange)...)

	return fullNode, nil
}

func (n *FullNode) Start(ctx context.Context) error {
	var err error
	n.logger.Infof("starting full node")

	if err = n.p2pServer.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start p2p server")
	}

	n.p2pServer.RegisterProtocols(n.Protocols(n.peerExchange)...)

	// factor := 1

	for {
		select {
		case <-ctx.Done():
			return n.Stop(ctx)
		default:
			{
				//time.Sleep(time.Second * time.Duration(factor))
				//
				//req := rhz.MsgGetBlocks{IndexHave: 0, IndexNeed: 10 * uint64(factor)}
				//if err := n.p2pServer.StreamMsg(
				//	ctx,
				//	rhz.ProtocolRequestBlocks,
				//	req,
				//); err != nil {
				//	n.logger.Warn(err)
				//
				//	break
				//}
				//
				//factor++
				//n.logger.Debugw("request sent", "req", req)
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
			ID:  rhz.ProtocolRequestBlocks,
			Run: n.requestHandler(backend),
		},
		{
			ID:  rhz.ProtocolResponseBlocks,
			Run: n.responseHandler(backend),
		},
		{
			ID:  rhz.ProtocolRequestDelegates,
			Run: n.requestHandler(backend),
		},
		{
			ID:  rhz.ProtocolResponseDelegates,
			Run: n.requestHandler(backend),
		},
		//{
		//	ID:  "rhz_test",
		//	Run: n.requestHandler(backend),
		//},
	}
}

func (n *FullNode) requestHandler(backend rhz.PeerExchange) func(ctx context.Context, rw p2p.MsgReadWriter) error {
	return func(ctx context.Context, rw p2p.MsgReadWriter) error {
		peer := rhz.NewPeer(rw)
		return rhz.HandleRequest(ctx, backend, peer)
	}
}

func (n *FullNode) responseHandler(backend rhz.PeerExchange) func(ctx context.Context, rw p2p.MsgReadWriter) error {
	return func(ctx context.Context, rw p2p.MsgReadWriter) error {
		peer := rhz.NewPeer(rw)
		return rhz.HandleResponse(ctx, backend, peer)
	}
}
