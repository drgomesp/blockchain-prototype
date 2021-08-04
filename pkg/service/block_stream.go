package service

import (
	"io"

	"github.com/drgomesp/rhizom/proto/gen/entity"
	"github.com/drgomesp/rhizom/proto/gen/message"
	"github.com/drgomesp/rhizom/proto/gen/stream"
)

// BlockStream is responsible to stream blocks through requests.
// This type implements one of the interfaces provided from the generated protobuff code.
type BlockStream struct {
	stream.UnimplementedBlockServer
}

// NewBlockStream instantiates a new service type of BlockStream object.
func NewBlockStream() *BlockStream {
	return &BlockStream{}
}

// GetBlocks calls the core service GetBlocks method and maps the result to a grpc service response.
func (b *BlockStream) GetBlock(stream stream.Block_GetBlockServer) error {
	for {
		var resp *message.ResponseStreamGetBlock

		switch req, err := stream.Recv(); err {
		case nil:
			resp = &message.ResponseStreamGetBlock{
				Block: &entity.Block{Index: req.IndexWant}}

		case io.EOF:
			return nil

		default:
			resp = &message.ResponseStreamGetBlock{Err: err.Error()}
		}

		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}
