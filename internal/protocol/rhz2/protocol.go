package rhz2

import (
	"context"

	"github.com/drgomesp/rhizom/pkg/p2p"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

const ProtocolID = p2p.ProtocolType("/rhz2/1.0.0")

const (
	MsgTypeRequest = iota
	MsgTypeResponse
)

const (
	MsgTypeGetBlocksRequest  = p2p.MsgType("GetBlocksRequest")
	MsgTypeGetBlocksResponse = p2p.MsgType("GetBlocksResponse")
)

type (
	requestHandlerFunc  func(p2p.MsgDecoder) (p2p.ProtocolType, proto.Message, error)
	responseHandlerFunc func(p2p.MsgDecoder) error
)

var requestHandlers = map[p2p.MsgType]requestHandlerFunc{
	MsgTypeGetBlocksRequest: HandleGetBlocksRequest,
}

var responseHandlers = map[p2p.MsgType]responseHandlerFunc{
	MsgTypeGetBlocksResponse: HandleGetBlocksResponse,
}

func ProtocolHandlerFunc(msgType int) p2p.StreamHandlerFunc {
	return func(ctx context.Context, rw p2p.MsgReadWriter) (p2p.ProtocolType, interface{}, error) {
		msg, err := rw.ReadMsg(ctx)
		if err != nil {
			return p2p.NilProtocol, nil, errors.Wrap(err, "message read failed")
		}

		switch msgType {
		case MsgTypeRequest:
			if handlerFunc := requestHandlers[msg.Type]; handlerFunc != nil {
				rpid, res, err := handlerFunc(msg)
				if err != nil {
					return p2p.NilProtocol, nil, err
				}

				return rpid, res, nil
			}

		case MsgTypeResponse:
			if handlerFunc := responseHandlers[msg.Type]; handlerFunc != nil {
				if err := handlerFunc(msg); err != nil {
					return p2p.NilProtocol, nil, err
				}

				return p2p.NilProtocol, nil, nil
			}
		}

		return p2p.NilProtocol, nil, errors.Errorf("unsupported message type '%s'", msg.Type)
	}
}
