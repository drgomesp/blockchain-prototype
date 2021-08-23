package rhz2

import (
	"context"

	pb "github.com/drgomesp/rhizom/internal/protocol/rhz2/pb"
	"github.com/drgomesp/rhizom/internal/rhz"
	"github.com/drgomesp/rhizom/pkg/block"
	"github.com/drgomesp/rhizom/pkg/p2p"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

func HandleGetBlocksRequest(ctx context.Context, peering rhz.Peering, msg p2p.MsgDecoder) (
	proto.Message, p2p.ProtocolType, proto.Message, error,
) {
	req := new(pb.GetBlocks_Request)
	if err := msg.Decode(req); err != nil {
		return nil, p2p.NilProtocol, nil, errors.Wrap(err, "message decode failed")
	}

	blocks, err := peering.GetBlocks(ctx, req.Index)
	if err != nil {
		return req, p2p.NilProtocol, nil, errors.Wrap(err, "failed to get blocks")
	}

	return req, p2p.ProtocolType(MsgTypeGetBlocksResponse), &pb.GetBlocks_Response{
		Blocks: wrapBlocks(blocks),
	}, nil
}

func wrapBlocks(blocks []*block.Block) []*pb.Block {
	wrappedBlocks := make([]*pb.Block, 0)
	for _, blk := range blocks {
		wrappedBlocks = append(wrappedBlocks, &pb.Block{
			Header: &pb.Block_Header{
				Index: blk.Header.Index,
			},
		})
	}

	return wrappedBlocks
}

func HandleGetBlocksResponse(_ context.Context, _ rhz.Peering, msg p2p.MsgDecoder) (
	proto.Message,
	error,
) {
	res := new(pb.GetBlocks_Response)
	if err := msg.Decode(res); err != nil {
		return nil, errors.Wrap(err, "message decode failed")
	}

	return res, nil
}
