package node

import "github.com/drgomesp/rhizom/pkg/p2p"

// Config defines the node configuration options.
type Config struct {
	Type Type   // Type of the node.
	Name string // Name of the node.

	P2P p2p.Config
}
