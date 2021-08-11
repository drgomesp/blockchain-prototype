package node

import (
	"context"

	"github.com/drgomesp/rhizom/pkg/p2p"
	"github.com/drgomesp/rhizom/pkg/rpc"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Node implements a multi-protocol node in network.
type Node struct {
	config    Config
	logger    *zap.SugaredLogger
	server    *p2p.Server
	services  []Service
	apis      []*rpc.API
	protocols []p2p.Protocol
}

type Option func(*Node)

func WithLogger(l *zap.SugaredLogger) Option {
	return func(node *Node) {
		node.logger = l
	}
}

func New(config Config, opt ...Option) (n *Node, err error) {
	n = &Node{
		config: config,
		apis:   make([]*rpc.API, 0),
	}

	for _, option := range opt {
		option(n)
	}

	if n.server, err = p2p.NewServer(config.P2P, p2p.WithLogger(n.logger)); err != nil {
		return nil, errors.Wrap(err, "failed to create p2p server")
	}

	return n, nil
}

func (n *Node) RegisterAPIs(apis ...*rpc.API) {
	n.apis = append(n.apis, apis...)
}

func (n *Node) RegisterServices(services ...Service) {
	n.services = append(n.services, services...)
}

func (n *Node) RegisterProtocols(protocols ...p2p.Protocol) {
	n.protocols = protocols
}

func (n *Node) Services() []Service {
	return n.services
}

func (n *Node) Server() *p2p.Server {
	return n.server
}

func (n *Node) Name() string {
	panic("implement me")
}

func (n *Node) Start(ctx context.Context) error {
	if err := n.server.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start server")
	}

	n.server.RegisterProtocols(n.protocols...)

	for _, service := range n.services {
		if err := service.Start(ctx); err != nil {
			return errors.Wrap(err, "failed to start service")
		}
	}

	return nil
}

func (n *Node) Stop(ctx context.Context) error {
	panic("implement me")
}
