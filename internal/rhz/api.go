package rhz

import "context"

// API defines the bi-directional protocol messages that can be exchanged by peers.
type API interface {
	// GetBlocks requests the peer for some blocks.
	GetBlocks(context.Context, *Peer, MsgGetBlocks) (MsgBlocks, error)
	// Blocks receives blocks from the peer.
	Blocks(context.Context, *Peer, MsgBlocks) error
}
