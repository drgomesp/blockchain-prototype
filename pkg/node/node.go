package node

// Node implements a multi-protocol node in the Rhizom network.
type Node struct {
	cfg *Config
}

func New(cfg *Config) (*Node, error) {
	node := &Node{
		cfg: cfg,
	}

	return node, nil
}
