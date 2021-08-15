package rhz1

import (
	"context"

	"github.com/drgomesp/rhizom/pkg/p2p"
	"github.com/pkg/errors"
)

// HandleGetBlocks handles an incoming request for blocks message.
func HandleGetBlocks(ctx context.Context, peering Peering, msg Message) (
	p2p.ProtocolType, MessagePacket, error,
) {
	var req MsgGetBlocks
	if err := msg.Decode(&req); err != nil {
		return p2p.NilProtocol, nil, errors.Wrap(err, "message decode failed")
	}

	blocks, err := peering.GetBlocks(ctx, peering, req)
	if err != nil {
		return p2p.NilProtocol, nil, ErrMessageHandleFailed(err)
	}

	return p2p.ProtocolType(MsgTypeBlocks), blocks, nil
}

// HandleBlocks handles an incoming blocks response.
func HandleBlocks(ctx context.Context, peering Peering, msg Message) error {
	var res MsgBlocks
	if err := msg.Decode(&res); err != nil {
		return err
	}

	if err := peering.Blocks(ctx, peering, res); err != nil {
		return ErrMessageHandleFailed(err)
	}

	return nil
}
