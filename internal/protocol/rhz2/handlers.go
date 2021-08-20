package rhz2

import (
	"log"

	pb "github.com/drgomesp/rhizom/internal/protocol/rhz2/pb"
	"github.com/drgomesp/rhizom/pkg/p2p"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

func HandleGetBlocksRequest(msg p2p.MsgDecoder) (p2p.ProtocolType, proto.Message, error) {
	var req pb.GetBlocks_Request
	if err := msg.Decode(&req); err != nil {
		return p2p.NilProtocol, nil, errors.Wrap(err, "message decode failed")
	}

	return p2p.ProtocolType(MsgTypeGetBlocksResponse), &pb.GetBlocks_Response{
		Blocks: []*pb.Block{
			{Header: &pb.Block_Header{Index: req.Index * 2}},
		},
	}, nil
}

func HandleGetBlocksResponse(msg p2p.MsgDecoder) error {
	var res pb.GetBlocks_Response
	if err := msg.Decode(&res); err != nil {
		return errors.Wrap(err, "message decode failed")
	}

	log.Printf("%+v", res)

	return nil
}
