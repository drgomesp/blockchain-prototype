package rhznode

import (
	"context"
	"time"

	"github.com/drgomesp/rhizom/internal/rhz"
	"github.com/drgomesp/rhizom/pkg/node"
	"github.com/drgomesp/rhizom/pkg/p2p"
	"github.com/drgomesp/rhizom/pkg/rpc"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var bootstrapAddrs = []string{
	"/dns4/bootstrapper-1.rhz.network/tcp/4001/ipfs/Qmf8Lt1FiQnG7tLrQbhwvUXzBMYsj6KicNdKiD1F2rSRW5",
	"/dns4/bootstrapper-2.rhz.network/tcp/4001/ipfs/QmcRoi1mQ7eb7xPDhWZjGL8rivAUHwCv1FMiLw7FGSZvFL",
}

const (
	rhzPrefix = "/rhz/"
	net       = "default_2b678c95-27d5-4f09-bf38-a62be2c5339b"
	// net               = "default_b73c9ec0-b225-4f1a-b137-7c13ea24c387".
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
}

func NewFullNode(ctx context.Context, logger *zap.SugaredLogger) (*FullNode, error) {
	p2pServer, err := p2p.NewServer(ctx, logger, p2p.Config{
		MaxPeers:       p2pServerMaxPeers,
		PingTimeout:    p2pServerPingTimeout,
		BootstrapAddrs: bootstrapAddrs,
		Topics:         topics,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize p2p server")
	}

	rpcServer, err := rpc.NewServer()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize rpc server")
	}

	n, err := node.New(node.Config{
		Type: node.TypeFull,
		Name: "rhz_node",
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed initialize node")
	}

	n.RegisterAPIs(APIs()...)
	n.RegisterServers(p2pServer, rpcServer)

	fullNode := &FullNode{
		node:   n,
		logger: logger,
	}

	p2pServer.RegisterProtocols(fullNode.Protocols(rhz.API(nil))...)

	return fullNode, nil
}

func (n *FullNode) Start(ctx context.Context) error {
	for _, srv := range n.node.Servers() {
		if err := srv.Start(ctx); err != nil {
			break
		}

		n.logger.Infof("%s server started", srv.Name())
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

func APIs() []*rpc.API {
	return []*rpc.API{}
}

func (n *FullNode) Protocols(api rhz.API) []p2p.Protocol {
	return []p2p.Protocol{
		//{
		//	Name: ProtocolRequestBlocks,
		//	Run:  MessageHandler(api),
		//},
		//{
		//	Name: ProtocolResponseBlocks,
		//	Run:  MessageHandler(api),
		//},
		//{
		//	Name: ProtocolRequestDelegates,
		//	Run:  MessageHandler(api),
		//},
		//{
		//	Name: ProtocolResponseDelegates,
		//	Run:  MessageHandler(api),
		//},
	}
}

//
//func MessageHandler(api rhz.API) rhz.MsgHandlerFunc {
//	return func(p *p2p.Peer, msg *interface{}) error {
//		panic("TODO")
//	}
//}
