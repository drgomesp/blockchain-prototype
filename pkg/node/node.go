package node

import "github.com/rhizomplatform/rhizom/pkg/rpc"

// Node implements a multi-protocol node in the Rhizom network.
type Node struct {
	cfg *Config

	servers []Server
	apis    []*rpc.API
}

func New(cfg *Config) (*Node, error) {
	node := &Node{
		cfg:     cfg,
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

func (n *Node) Servers() []Server {
	return n.servers
}

func (n *Node) APIs() []*rpc.API {
	return n.apis
}
