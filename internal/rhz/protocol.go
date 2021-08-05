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

func ProtocolHandlerFunc(msgType uint, api API) func(ctx context.Context, rw p2p.MsgReadWriter) error {
	return func(ctx context.Context, rw p2p.MsgReadWriter) error {
		peer := NewPeer(rw)

		switch msgType {
		case MsgTypeRequest:
			return HandleRequest(ctx, api, peer)
		case MsgTypeResponse:
			return HandleResponse(ctx, api, peer)
		}

		return errors.New("message handle failed")
	}
}

type (
	RequestMsgHandler  func(API, Message, *Peer) (Message, error)
	ResponseMsgHandler func(API, Message, *Peer) error
)

func HandleRequest(ctx context.Context, api API, peer *Peer) error {
	msg, err := peer.rw.ReadMsg(ctx)
	if err != nil {
		return errors.Wrap(err, "read peer message failed")
	}

	if handlerFunc := requestHandlers[msg.Type]; handlerFunc != nil {
		// TODO: discard the response, but maybe should store locally
		_, err := handlerFunc(api, msg, peer)

		return err
	}

	return ErrRequestTypeNotSupported
}

func HandleResponse(ctx context.Context, api API, peer *Peer) error {
	msg, err := peer.rw.ReadMsg(ctx)
	if err != nil {
		return errors.Wrap(err, "read peer message failed")
	}

	if handlerFunc := responseHandlers[msg.Type]; handlerFunc != nil {
		return handlerFunc(api, msg, peer)
	}

	return ErrResponseTypeNotSupported
}
