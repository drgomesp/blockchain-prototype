package rhz

import (
	"context"

	"github.com/drgomesp/rhizom/internal"
	"github.com/drgomesp/rhizom/pkg/p2p"
	"github.com/pkg/errors"
)

const (
	MsgTypeRequest = iota
	MsgTypeResponse
)

// Message used for direct p2p communication.
type Message interface {
	// Decode ...
	Decode(v interface{}) error
}

// MessagePacket defines the message packet type which carries the message data.
// Protocol-specific messages must implement this interface in order to be supported.
type MessagePacket interface {
	Type() p2p.MsgType
}

// Request/response messages.
const (
	// MsgTypeGetBlocks represents a request for blocks.
	MsgTypeGetBlocks = p2p.MsgType("/rhz/blocks/req/" + internal.NetworkName)
	// MsgTypeBlocks represents a response to a request for blocks.
	MsgTypeBlocks = p2p.MsgType("/rhz/blocks/resp/" + internal.NetworkName)
)

// Async broadcast messages.
const (
	// MsgTypeNewBlock represents a new block.
	MsgTypeNewBlock = p2p.MsgType("NewBlock")
)

// Protocol-specific errors.
var (
	ErrMessageHandleFailed = func(e error) error {
		return errors.Wrap(e, "message handle failed")
	}
	ErrUnsupportedMessageType = func(msgType p2p.MsgType) error {
		return errors.Wrap(
			errors.New("unsupported message type"), string(msgType),
		)
	}
	ErrRequestTypeNotSupported = func(t p2p.MsgType) error {
		return errors.Errorf("request type '%s' not supported by protocol handler", t)
	}
	ErrResponseTypeNotSupported = func(t p2p.MsgType) error {
		return errors.Errorf("response type '%s' not supported by protocol handler", t)
	}
)

type (
	// requestHandlerFunc defines the function type for handling protocol request messages.
	requestHandlerFunc func(context.Context, Peering, Message) (p2p.ProtocolType, MessagePacket, error)
	// responseHandlerFunc defines the function type for handling protocol response messages.
	responseHandlerFunc func(context.Context, Peering, Message) error
)

var (
	requestHandlers = map[p2p.MsgType]requestHandlerFunc{
		MsgTypeGetBlocks: HandleGetBlocks,
	}
	responseHandlers = map[p2p.MsgType]responseHandlerFunc{
		MsgTypeBlocks: HandleBlocks,
	}
)

// ProtocolHandlerFunc is the actual protocol handler implementation required by p2p.Protocol.
func ProtocolHandlerFunc(msgType uint, peering Peering) p2p.StreamHandlerFunc {
	return func(ctx context.Context, rw p2p.MsgReadWriter) (p2p.ProtocolType, interface{}, error) {
		_ = NewPeer(rw) // TODO: see if there's a use for peer still.

		msg, err := rw.ReadMsg(ctx)
		if err != nil {
			return p2p.NilProtocol, nil, errors.Wrap(err, "message read failed")
		}

		switch msgType {
		case MsgTypeRequest:
			return handleRequest(ctx, peering, msg)
		case MsgTypeResponse:
			return p2p.NilProtocol, nil, handleResponse(ctx, peering, msg)
		}

		return p2p.NilProtocol, nil, ErrUnsupportedMessageType(msg.Type)
	}
}

// handleRequest handles request type messages.
func handleRequest(ctx context.Context, peering Peering, msg *p2p.Message) (p2p.ProtocolType, interface{}, error) {
	if handlerFunc := requestHandlers[msg.Type]; handlerFunc != nil {
		pid, res, err := handlerFunc(ctx, peering, msg)
		if err != nil {
			return p2p.NilProtocol, nil, err
		}

		return pid, res, nil
	}

	return p2p.NilProtocol, nil, ErrRequestTypeNotSupported(msg.Type)
}

// handleResponse handles response type messages.
func handleResponse(ctx context.Context, peering Peering, msg *p2p.Message) error {
	if handlerFunc := responseHandlers[msg.Type]; handlerFunc != nil {
		return handlerFunc(ctx, peering, msg)
	}

	return ErrResponseTypeNotSupported(msg.Type)
}
