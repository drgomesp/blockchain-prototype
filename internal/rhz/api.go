package rhz

// API is a bi-directional protocol for peer message exchange.
type API interface {
	// GetBlocks requests the peer for some blocks.
	GetBlocks(*Peer, *MsgGetBlocks) (*MsgBlocks, error)
	// Blocks receives blocks from the peer.
	Blocks(*Peer, *MsgBlocks) error
}
