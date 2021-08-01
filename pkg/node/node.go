package node

import (
	"github.com/drgomesp/rhizom/pkg/p2p"
	"github.com/drgomesp/rhizom/pkg/rpc"
)

// Node implements a multi-protocol node in the Rhizom network.
type Node struct {
	config Config

	servers   []Server
	apis      []*rpc.API
	protocols []p2p.Protocol
}

func New(config Config) (*Node, error) {
	node := &Node{
		config:  config,
		servers: make([]Server, 0),
		apis:    make([]*rpc.API, 0),
	}

	return node, nil
}

func (n *Node) RegisterAPIs(apis ...*rpc.API) {
	n.apis = append(n.apis, apis...)
}

func (n *Node) RegisterServers(servers ...Server) {
	n.servers = append(n.servers, servers...)
}

func (n *Node) RegisterProtocols(protocols ...p2p.Protocol) {
	n.protocols = append(n.protocols, protocols...)
}

func (n *Node) Config() Config {
	return n.config
}

func (n *Node) Servers() []Server {
	return n.servers
}

func (n *Node) APIs() []*rpc.API {
	return n.apis
}
