package rhz

import (
	"context"

	"github.com/drgomesp/rhizom/pkg/block"
)

// Peering defines the bi-directional protocol messages that can be exchanged by peers.
type Peering interface {
	// GetBlocks requests the peer for some blocks.
	GetBlocks(ctx context.Context, index uint64) ([]*block.Block, error)

	// Blocks receives blocks from the peer.
	Blocks(ctx context.Context, blocks []*block.Block) error
}

// Broadcast defines the network messages that a peer can send and receive.
type Broadcast interface {
	// OnNewBlock is the listener for new block messages.
	OnNewBlock(context.Context, *block.Block) error
}
