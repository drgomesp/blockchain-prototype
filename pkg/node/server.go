package node

// Server defines a server that can be registered to the node.
type Server interface {
	Start() error
	Stop() error
}
