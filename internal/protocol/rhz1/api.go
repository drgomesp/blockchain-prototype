package rhz1

import (
	"context"
)

// Peering defines the bi-directional protocol messages that can be exchanged by peers.
type Peering interface {
	// GetBlocks requests the peer for some blocks.
	GetBlocks(context.Context, Peering, MsgGetBlocks) (MsgBlocks, error)
	// Blocks receives blocks from the peer.
	Blocks(context.Context, Peering, MsgBlocks) error
}

// Broadcast defines the network messages that a peer can send and receive.
type Broadcast interface {
	// OnNewBlock is the listener for new block messages.
	OnNewBlock(context.Context, MsgNewBlock) error
}
