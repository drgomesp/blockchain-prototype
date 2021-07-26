package rhznode

import (
	"github.com/rhizomplatform/rhizom/pkg/node"
	"github.com/rhizomplatform/rhizom/pkg/rpc"
)

type FullNode struct{}

func New(node *node.Node) (*FullNode, error) {
	rhz := &FullNode{}

	node.RegisterAPIs(rhz.APIs()...)
	node.RegisterServers(rhz.Servers()...)

	return rhz, nil
}

func (n *FullNode) APIs() []*rpc.API {
	return []*rpc.API{}
}

func (n *FullNode) Servers() []node.Server {
	return []node.Server{}
}

func (n *FullNode) Start() error {
	panic("implement me")
}

func (n *FullNode) Stop() error {
	panic("implement me")
}
