package rhz2

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"

	rhz2 "github.com/drgomesp/acervo/internal/protocol/rhz2/pb"
	"github.com/drgomesp/acervo/internal/rhz"
	"github.com/drgomesp/acervo/pkg/p2p"
)

var Logger *zap.SugaredLogger

const (
	MsgTypeRequest = iota
	MsgTypeResponse
)

const (
	MsgTypeGetBlocksRequest  = p2p.MsgType("GetBlocksRequest")
	MsgTypeGetBlocksResponse = p2p.MsgType("GetBlocksResponse")
)

type (
	requestHandlerFunc func(context.Context, rhz.Peering, p2p.MsgDecoder) (
		proto.Message, p2p.ProtocolType, proto.Message, error,
	)

	responseHandlerFunc func(context.Context, rhz.Peering, p2p.MsgDecoder) (proto.Message, error)
)

var requestHandlers = map[p2p.MsgType]requestHandlerFunc{
	MsgTypeGetBlocksRequest: func(ctx context.Context, peering rhz.Peering, decoder p2p.MsgDecoder) (
		proto.Message, p2p.ProtocolType, proto.Message, error,
	) {
		msg := &rhz2.GetBlocks_Response{
			Updated: true,
			Blocks: []*rhz2.Block{
				{Header: &rhz2.Block_Header{Index: 1}},
				{Header: &rhz2.Block_Header{Index: 2}},
				{Header: &rhz2.Block_Header{Index: 3}},
			},
		}
		spew.Dump(decoder)
		return nil, p2p.ProtocolType(MsgTypeGetBlocksResponse), msg, nil
	},
}

var responseHandlers = map[p2p.MsgType]responseHandlerFunc{
	MsgTypeGetBlocksResponse: func(ctx context.Context, peering rhz.Peering, decoder p2p.MsgDecoder) (proto.Message, error) {
		spew.Dump(decoder)
		return nil, nil
	},
}

func ProtocolHandlerFunc(msgType int, peering rhz.Peering) p2p.StreamHandlerFunc {
	return func(ctx context.Context, rw p2p.MsgReadWriter) (p2p.ProtocolType, proto.Message, error) {
		msg, err := rw.ReadMsg(ctx)
		if err != nil {
			return p2p.NilProtocol, nil, errors.Wrap(err, "message read failed")
		}

		switch msgType {
		case MsgTypeRequest:
			if handlerFunc := requestHandlers[msg.Type]; handlerFunc != nil {
				req, rpid, res, err := handlerFunc(ctx, peering, msg)
				if err != nil {
					return p2p.NilProtocol, nil, err
				}

				log.Debug().Msgf("request received: %+v", req)

				return rpid, res, nil
			}

		case MsgTypeResponse:
			if handlerFunc := responseHandlers[msg.Type]; handlerFunc != nil {
				var res proto.Message

				if res, err = handlerFunc(ctx, peering, msg); err != nil {
					return p2p.NilProtocol, nil, err
				}

				log.Debug().Msgf("response received: %+v", res)

				return p2p.NilProtocol, nil, nil
			}
		}

		return p2p.NilProtocol, nil, errors.Errorf("unsupported message type '%s'", msg.Type)
	}
}