package rhz

import (
	"context"

	"github.com/drgomesp/rhizom/pkg/p2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/pkg/errors"
)

const (
	MsgTypeRequest = iota
	MsgTypeResponse
)

const (
	// MsgTypeGetBlocks represents a request for blocks.
	MsgTypeGetBlocks = p2p.MsgType(ProtocolRequestBlocks)
	// MsgTypeBlocks represents a response to a request for blocks.
	MsgTypeBlocks = p2p.MsgType(ProtocolResponseBlocks)
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

func ProtocolHandlerFunc(msgType uint, api API) p2p.StreamHandlerFunc {
	return func(ctx context.Context, stream network.Stream) (p2p.ProtocolType, interface{}, error) {
		switch msgType {
		case MsgTypeRequest:
			return HandleRequest(ctx, api, stream)
		case MsgTypeResponse:
			return p2p.NilProtocol, nil, HandleResponse(ctx, api, stream)
		}

		return p2p.NilProtocol, nil, errors.New("not implemented")
	}
}

type (
	RequestMsgHandler  func(context.Context, API, Message, *Peer) (p2p.ProtocolType, MessagePacket, error)
	ResponseMsgHandler func(context.Context, API, Message, *Peer) error
)

func HandleRequest(ctx context.Context, api API, peer network.Stream) (p2p.ProtocolType, interface{}, error) {
	msg := &p2p.Message{
		Type:    p2p.MsgType(peer.Protocol()),
		Payload: peer,
	}

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
func HandleResponse(ctx context.Context, api API, peer network.Stream) error {
	msg := &p2p.Message{
		Type:    p2p.MsgType(peer.Protocol()),
		Payload: peer,
	}

	if handlerFunc := responseHandlers[msg.Type]; handlerFunc != nil {
		return handlerFunc(ctx, api, msg, nil)
	}

	return ErrResponseTypeNotSupported
}
