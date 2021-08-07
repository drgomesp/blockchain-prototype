package rhz

import "context"

// API is a bi-directional protocol for peer message exchange.
type API interface {
	// GetBlocks requests the peer for some blocks.
	GetBlocks(context.Context, *Peer, MsgGetBlocks) (MsgBlocks, error)
	// Blocks receives blocks from the peer.
	Blocks(context.Context, *Peer, MsgBlocks) error
}
