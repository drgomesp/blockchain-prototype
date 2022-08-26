package rhz1

import (
	"context"

	"github.com/pkg/errors"

	"github.com/drgomesp/acervo/internal/rhz"
	"github.com/drgomesp/acervo/pkg/p2p"
)

// HandleGetBlocks handles an incoming request for blocks message.
func HandleGetBlocks(ctx context.Context, peering rhz.Peering, msg Message) (
	p2p.ProtocolType, MessagePacket, error,
) {
	var req MsgGetBlocks
	if err := msg.Decode(&req); err != nil {
		return p2p.NilProtocol, nil, errors.Wrap(err, "message decode failed")
	}

	blocks, err := peering.GetBlocks(ctx, req.IndexNeed)
	if err != nil {
		return p2p.NilProtocol, nil, ErrMessageHandleFailed(err)
	}

	return p2p.ProtocolType(MsgTypeBlocks), MsgBlocks{
		IsUpdated: true,
		Chain:     blocks,
	}, nil
}

// HandleBlocks handles an incoming blocks response.
func HandleBlocks(ctx context.Context, peering rhz.Peering, msg Message) error {
	var res MsgBlocks
	if err := msg.Decode(&res); err != nil {
		return err
	}

	if err := peering.Blocks(ctx, res.Chain); err != nil {
		return ErrMessageHandleFailed(err)
	}

	return nil
}