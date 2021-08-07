package rhz

import (
	"context"

	"github.com/drgomesp/rhizom/pkg/p2p"
	"github.com/pkg/errors"
)

const (
	MsgTypeRequest = iota
	MsgTypeResponse
)

// protocol messages.
const (
	// MsgTypeGetBlocks represents a request for blocks.
	MsgTypeGetBlocks = p2p.MsgType(ProtocolRequestBlocks)
	// MsgTypeBlocks represents a response to a request for blocks.
	MsgTypeBlocks = p2p.MsgType(ProtocolResponseBlocks)
)

// async messages.
const (
	// MsgNewBlock represents a new block.
	MsgNewBlock = p2p.MsgType("NewBlock")
)

var (
	ErrRequestTypeNotSupported  = errors.New("request type not supported by protocol handler")
	ErrResponseTypeNotSupported = errors.New("response type not supported by protocol handler")

	requestHandlers = map[p2p.MsgType]RequestMsgHandler{
		MsgTypeGetBlocks: HandleGetBlocks,
	}
	responseHandlers = map[p2p.MsgType]ResponseMsgHandler{
		MsgTypeBlocks: HandleBlocks,
	}
)

type (
	RequestMsgHandler  func(context.Context, API, Message, *Peer) (p2p.ProtocolType, MessagePacket, error)
	ResponseMsgHandler func(context.Context, API, Message, *Peer) error
)

func ProtocolHandlerFunc(msgType uint, api API) p2p.StreamHandlerFunc {
	return func(ctx context.Context, rw p2p.MsgReadWriter) (p2p.ProtocolType, interface{}, error) {
		msg, err := rw.ReadMsg(ctx)
		if err != nil {
			return p2p.NilProtocol, nil, err
		}

		switch msgType {
		case MsgTypeRequest:
			return HandleRequest(ctx, api, msg)
		case MsgTypeResponse:
			return p2p.NilProtocol, nil, HandleResponse(ctx, api, msg)
		}

		return p2p.NilProtocol, nil, errors.New("not implemented")
	}
}

func HandleRequest(ctx context.Context, api API, msg *p2p.Message) (p2p.ProtocolType, interface{}, error) {
	if handlerFunc := requestHandlers[msg.Type]; handlerFunc != nil {
		pid, res, err := handlerFunc(ctx, api, msg, nil)
		if err != nil {
			return p2p.NilProtocol, nil, err
		}

		return pid, res, nil
	}

	return p2p.NilProtocol, nil, ErrRequestTypeNotSupported
}

// HandleResponse ... TODO: change peer to an actual Peer pointer.
func HandleResponse(ctx context.Context, api API, msg *p2p.Message) error {
	if handlerFunc := responseHandlers[msg.Type]; handlerFunc != nil {
		return handlerFunc(ctx, api, msg, nil)
	}

	return ErrResponseTypeNotSupported
}
