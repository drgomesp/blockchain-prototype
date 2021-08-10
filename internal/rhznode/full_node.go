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

// serviceTag is an identifier for the discovery service.
const serviceTag = "rhizom"

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

// FullNode implements a full node type in the Rhizom network.
type FullNode struct {
	node      *node.Node
	logger    *zap.SugaredLogger
	peering   rhz.Peering
	broadcast rhz.Broadcast
	p2pServer *p2p.Server
}

func NewFullNode(logger *zap.SugaredLogger) (*FullNode, error) {
	n, err := node.New(node.Config{
		Type: node.TypeFull,
		Name: "rhz_node",
		P2P: p2p.Config{
			ServiceTag:     serviceTag,
			MaxPeers:       p2pServerMaxPeers,
			PingTimeout:    p2pServerPingTimeout,
			BootstrapAddrs: bootstrapAddrs,
			Topics:         topics,
		},
	}, node.WithLogger(logger))
	if err != nil {
		return nil, errors.Wrap(err, "failed initialize node")
	}

	backend := NewHandler(logger)

	fullNode := &FullNode{
		node:      n,
		logger:    logger,
		p2pServer: n.Server(),
		peering:   backend,
		broadcast: backend,
	}

	n.RegisterAPIs(nil)
	n.RegisterProtocols(fullNode.Protocols(fullNode.peering)...)

	return fullNode, nil
}

func (n *FullNode) Start(ctx context.Context) (err error) {
	n.logger.Info("starting full node")

	if err = n.p2pServer.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start p2p server")
	}

	n.p2pServer.RegisterProtocols(n.Protocols(n.peering)...)

	for {
		select {
		case <-ctx.Done():
			return n.Stop(ctx)
		default:
			{
				factor := 5
				time.Sleep(time.Second * time.Duration(factor))

				if err := p2p.Send(
					ctx,
					n.p2pServer,
					rhz.MsgTypeGetBlocks,
					rhz.MsgGetBlocks{IndexHave: 0, IndexNeed: 1000},
				); err != nil {
					if errors.Is(err, p2p.ErrNoPeersFound) {
						continue
					}

					n.logger.Error(err)
				}
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

func (n *FullNode) Protocols(backend rhz.Peering) []p2p.Protocol {
	return []p2p.Protocol{
		{
			ID:  rhz.ProtocolRequestBlocks,
			Run: rhz.ProtocolHandlerFunc(rhz.MsgTypeRequest, backend),
		},
		{
			ID:  rhz.ProtocolResponseBlocks,
			Run: rhz.ProtocolHandlerFunc(rhz.MsgTypeResponse, backend),
		},
		{
			ID:  rhz.ProtocolRequestDelegates,
			Run: rhz.ProtocolHandlerFunc(rhz.MsgTypeRequest, backend),
		},
		{
			ID:  rhz.ProtocolResponseDelegates,
			Run: rhz.ProtocolHandlerFunc(rhz.MsgTypeResponse, backend),
		},
	}
}
