package rhz2

import (
	"context"

	"github.com/drgomesp/rhizom/internal/rhz"
	"github.com/drgomesp/rhizom/pkg/p2p"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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
	MsgTypeGetBlocksRequest: HandleGetBlocksRequest,
}

var responseHandlers = map[p2p.MsgType]responseHandlerFunc{
	MsgTypeGetBlocksResponse: HandleGetBlocksResponse,
}

func ProtocolHandlerFunc(msgType int, peering rhz.Peering) p2p.StreamHandlerFunc {
	return func(ctx context.Context, rw p2p.MsgReadWriter) (p2p.ProtocolType, interface{}, error) {
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

				Logger.Debugf("request received: %+v", req)

				return rpid, res, nil
			}

		case MsgTypeResponse:
			if handlerFunc := responseHandlers[msg.Type]; handlerFunc != nil {
				var res proto.Message

				if res, err = handlerFunc(ctx, peering, msg); err != nil {
					return p2p.NilProtocol, nil, err
				}

				Logger.Debugf("response received: %+v", res)

				return p2p.NilProtocol, nil, nil
			}
		}

		return p2p.NilProtocol, nil, errors.Errorf("unsupported message type '%s'", msg.Type)
	}
}
