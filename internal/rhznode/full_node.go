package rhznode

import (
	"context"
	"time"

	"github.com/drgomesp/rhizom/internal"
	"github.com/drgomesp/rhizom/internal/protocol/rhz1"
	"github.com/drgomesp/rhizom/internal/protocol/rhz2"
	pb "github.com/drgomesp/rhizom/internal/protocol/rhz2/pb"
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
	TopicBlocks       = "/rhz/blk/" + internal.NetworkName
	TopicProducers    = "/rhz/prc/" + internal.NetworkName
	TopicTransactions = "/rhz/tx/" + internal.NetworkName
	TopicRequestSync  = "/rhz/blkchain/req/" + internal.NetworkName

	p2pServerMaxPeers    = 5
	p2pServerPingTimeout = time.Second * 5
)

// FullNode implements a full node type in the Rhizom network.
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

	fullNode := &FullNode{
		logger:    logger,
		p2pServer: n.Server(),
		peering:   rhz.NewPeeringService(logger),
		// broadcast: peeringService,
	}

	// TODO: this needs to be injected somehow.
	rhz2.Logger = logger

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
				factor := 1
				var msgType p2p.MsgType

				msgType = rhz2.MsgTypeGetBlocksRequest
				if err := p2p.Send(
					ctx,
					n.p2pServer,
					msgType,
					&pb.GetBlocks_Request{
						Index: 55,
					},
				); err != nil {
					if errors.Is(err, p2p.ErrNoPeersFound) {
						continue
					}

					n.logger.Errorw(err.Error(), "protocol", msgType)
				}

				msgType = rhz1.MsgTypeGetBlocks
				if err := p2p.Send(
					ctx,
					n.p2pServer,
					msgType,
					rhz1.MsgGetBlocks{
						IndexHave: 0,
						IndexNeed: 10,
					},
				); err != nil {
					if errors.Is(err, p2p.ErrNoPeersFound) {
						continue
					}

					n.logger.Errorw(err.Error(), "protocol", msgType)
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
			ID:  string(rhz1.MsgTypeGetBlocks),
			Run: rhz1.ProtocolHandlerFunc(rhz1.MsgTypeRequest, backend),
		},
		{
			ID:  string(rhz1.MsgTypeBlocks),
			Run: rhz1.ProtocolHandlerFunc(rhz1.MsgTypeResponse, backend),
		},
		{
			ID:  string(rhz2.MsgTypeGetBlocksRequest),
			Run: rhz2.ProtocolHandlerFunc(rhz2.MsgTypeRequest, backend),
		},
		{
			ID:  string(rhz2.MsgTypeGetBlocksResponse),
			Run: rhz2.ProtocolHandlerFunc(rhz2.MsgTypeResponse, backend),
		},
	}
}
