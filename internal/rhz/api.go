package rhz

import (
	"context"
)

// Streaming defines the bi-directional protocol messages that can be exchanged by peers.
type Streaming interface {
	// GetBlocks requests the peer for some blocks.
	GetBlocks(context.Context, *Peer, MsgGetBlocks) (MsgBlocks, error)
	// Blocks receives blocks from the peer.
	Blocks(context.Context, *Peer, MsgBlocks) error
}

// Broadcast defines the network messages that a peer can send and receive.
type Broadcast interface {
	// OnNewBlock is the listener for new block messages.
	OnNewBlock(context.Context, *Block) error
}
