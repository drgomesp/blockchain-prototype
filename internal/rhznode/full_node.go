package rhznode

import (
	"time"

	"github.com/rhizomplatform/rhizom/pkg/node"
	"github.com/rhizomplatform/rhizom/pkg/p2p"
	"github.com/rhizomplatform/rhizom/pkg/rpc"
	"go.uber.org/zap"
)

type FullNode struct {
	logger *zap.Logger
}

func NewFullNode(node *node.Node) (*FullNode, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	rhz := &FullNode{
		logger: logger,
	}

	node.RegisterAPIs(rhz.APIs()...)
	node.RegisterServers(rhz.Servers()...)

	return rhz, nil
}

func (n *FullNode) Start() error {
	for {
		n.logger.Info("running")

		time.Sleep(time.Second)
	}
}

func (n *FullNode) Stop() error {
	panic("implement me")
}

func (n *FullNode) APIs() []*rpc.API {
	return []*rpc.API{}
}

func (n *FullNode) Servers() []node.Server {
	return []node.Server{
		p2p.New(),
	}
}
