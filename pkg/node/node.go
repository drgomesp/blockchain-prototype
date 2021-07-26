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
		cfg: cfg,
	}

	return node, nil
}

func (n *Node) RegisterAPIs(apis ...*rpc.API) {
	n.apis = append(n.apis, apis...)
}

func (n *Node) RegisterServers(servers ...Server) {
	n.servers = append(n.servers, servers...)
}
