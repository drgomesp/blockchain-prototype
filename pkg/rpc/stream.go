package rpc

import (
	"fmt"
	"io"

	"github.com/drgomesp/rhizom/proto/gen/entity"
	"github.com/drgomesp/rhizom/proto/gen/message"
	"github.com/drgomesp/rhizom/proto/gen/service"
)

type BlockService struct {
	service.UnimplementedNodeServer
}

func NewStreamService() *BlockService {
	return &BlockService{}
}

// GetBlocks calls the core service GetBlocks method and maps the result to a grpc service response.
func (b *BlockService) GetBlock(stream service.Node_GetBlockServer) error {
	for i := 0; ; i++ {
		fmt.Printf("\rreq: %d", i)

		var resp *message.GetBlockResponse
		req, err := stream.Recv()

		switch err {
		case nil:
			resp = &message.GetBlockResponse{
				Block: &entity.Block{Index: req.Want},
			}

		case io.EOF:
			return nil

		default:
			resp = &message.GetBlockResponse{Err: err.Error()}
		}

		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}
