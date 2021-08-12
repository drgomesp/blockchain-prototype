package rhznode

import (
	"context"
	"time"

	"github.com/drgomesp/rhizom/internal"
	"github.com/drgomesp/rhizom/internal/protocol/rhz"
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
	TopicBlocks       = "/rhz/blk/" + internal.NetworkName
	TopicProducers    = "/rhz/prc/" + internal.NetworkName
	TopicTransactions = "/rhz/tx/" + internal.NetworkName
	TopicRequestSync  = "/rhz/blkchain/req/" + internal.NetworkName

	p2pServerMaxPeers    = 5
	p2pServerPingTimeout = time.Second * 5
)

// FullNode implements a full node type in the Rhizom NetworkName.
type FullNode struct {
	*node.Node
	logger    *zap.SugaredLogger
	peering   rhz.Peering
	broadcast rhz.Broadcast
	p2pServer *p2p.Server
}

func NewFullNode(logger *zap.SugaredLogger) (*node.Node, error) {
	n, err := node.New(node.Config{
		Type: node.TypeFull,
		Name: "full_node",
		P2P: p2p.Config{
			NetworkName:    internal.NetworkName,
			ServiceTag:     serviceTag,
			MaxPeers:       p2pServerMaxPeers,
			PingTimeout:    p2pServerPingTimeout,
			BootstrapAddrs: bootstrapAddrs,
			Topics: []string{
				TopicBlocks,
				TopicProducers,
				TopicTransactions,
				TopicRequestSync,
			},
		},
	}, node.WithLogger(logger))
	if err != nil {
		return nil, errors.Wrap(err, "failed initialize node")
	}

	backend := NewHandler(logger)

	fullNode := &FullNode{
		logger:    logger,
		p2pServer: n.Server(),
		peering:   backend,
		broadcast: backend,
	}

	n.RegisterAPIs(nil)
	n.RegisterProtocols(fullNode.Protocols(fullNode.peering)...)
	n.RegisterServices(fullNode)

	return n, nil
}

func (n *FullNode) Start(ctx context.Context) (err error) {
	n.logger.Info("starting full node")

	for {
		select {
		case <-ctx.Done():
			return n.Stop(ctx)
		default:
			{
				factor := 5

				p := rhz.MsgTypeGetBlocks
				if err := p2p.Send(
					ctx,
					n.p2pServer,
					p,
					rhz.MsgGetBlocks{IndexHave: 0, IndexNeed: 5},
				); err != nil {
					if errors.Is(err, p2p.ErrNoPeersFound) {
						continue
					}

					n.logger.Errorw(err.Error(), "protocol", p)
				}

				time.Sleep(time.Second * time.Duration(factor))
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
			ID:  string(rhz.MsgTypeGetBlocks),
			Run: rhz.ProtocolHandlerFunc(rhz.MsgTypeRequest, backend),
		},
		{
			ID:  string(rhz.MsgTypeBlocks),
			Run: rhz.ProtocolHandlerFunc(rhz.MsgTypeResponse, backend),
		},
	}
}
