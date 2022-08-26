package rhznode

import (
	"context"
	"time"

	config "github.com/ipfs/go-ipfs-config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/drgomesp/acervo/internal"
	"github.com/drgomesp/acervo/internal/protocol/rhz1"
	"github.com/drgomesp/acervo/internal/protocol/rhz2"
	"github.com/drgomesp/acervo/internal/rhz"
	"github.com/drgomesp/acervo/pkg/node"
	"github.com/drgomesp/acervo/pkg/p2p"
)

// serviceTag is an identifier for the discovery service.
const serviceTag = "acervo"

const (
	TopicBlocks       = "/rhz/blk/" + internal.NetworkName
	TopicProducers    = "/rhz/prc/" + internal.NetworkName
	TopicTransactions = "/rhz/tx/" + internal.NetworkName
	TopicRequestSync  = "/rhz/blkchain/req/" + internal.NetworkName

	p2pServerMaxPeers    = 5
	p2pServerPingTimeout = time.Second * 5
)

// FullNode implements a full node type in the acervo network.
type FullNode struct {
	*node.Node
	peering   rhz.Peering
	broadcast rhz.Broadcast
	p2pServer *p2p.Server
}

func NewFullNode() (*node.Node, error) {
	n, err := node.New(node.Config{
		Type: node.TypeFull,
		Name: "full_node",
		P2P: p2p.Config{
			NetworkName:    internal.NetworkName,
			ServiceTag:     serviceTag,
			MaxPeers:       p2pServerMaxPeers,
			PingTimeout:    p2pServerPingTimeout,
			BootstrapAddrs: config.DefaultBootstrapAddresses,
			Topics: []string{
				TopicBlocks,
				TopicProducers,
				TopicTransactions,
				TopicRequestSync,
			},
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed initialize node")
	}

	fullNode := &FullNode{
		p2pServer: n.Server(),
		peering:   rhz.NewPeeringService(),
		// broadcast: peeringService,
	}

	n.RegisterAPIs(nil)
	n.RegisterProtocols(fullNode.Protocols(fullNode.peering)...)
	n.RegisterServices(fullNode)

	return n, nil
}

func (n *FullNode) Start(ctx context.Context) (err error) {
	log.Info().Msg("starting full node")

	for {
		select {
		case <-ctx.Done():
			return n.Stop(ctx)
		default:
			{
				factor := 1
				var msgType p2p.MsgType

				//msgType = rhz2.MsgTypeGetBlocksRequest
				//if err := p2p.Send(
				//	ctx,
				//	n.p2pServer,
				//	msgType,
				//	&pb.GetBlocks_Request{
				//		Index: 55,
				//	},
				//); err != nil {
				//	if errors.Is(err, p2p.ErrNoPeersFound) {
				//		continue
				//	}
				//
				//	n.logger.Errorw(err.Error(), "protocol", msgType)
				//}

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

					log.Error().Err(err).Send()
				}

				time.Sleep(time.Second * time.Duration(factor))
			}
		}
	}
}

func (n *FullNode) Stop(_ context.Context) error {
	log.Info().Msg("stopping full node")

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