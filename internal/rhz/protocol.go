package rhz

import (
	"context"

	"github.com/drgomesp/rhizom/pkg/p2p"
	"github.com/pkg/errors"
)

const (
	rhzPrefix                 = "/rhz/"
	devNet                    = "default_2b678c95-27d5-4f09-bf38-a62be2c5339b"
	testNet                   = "rhz_testnet_e19d2c16-8c39-4f0f-8c88-2427e37c12bb"
	net                       = testNet
	TopicBlocks               = rhzPrefix + "blk/" + net
	TopicProducers            = rhzPrefix + "prc/" + net
	TopicTransactions         = rhzPrefix + "tx/" + net
	TopicRequestSync          = rhzPrefix + "blkchain/req/" + net
	ProtocolRequestBlocks     = rhzPrefix + "blocks/req/" + net
	ProtocolResponseBlocks    = rhzPrefix + "blocks/resp/" + net
	ProtocolRequestDelegates  = rhzPrefix + "delegates/req/" + net
	ProtocolResponseDelegates = rhzPrefix + "delegates/resp/" + net
)

// request/response messages.
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
)

type (
	RequestMsgHandler  func(PeerExchange, Message, *Peer) (MessagePacket, error)
	ResponseMsgHandler func(PeerExchange, Message, *Peer) error
)

var requestHandlers = map[p2p.MsgType]RequestMsgHandler{
	MsgTypeGetBlocks: HandleGetBlocks,
}

var responseHandlers = map[p2p.MsgType]ResponseMsgHandler{
	MsgTypeBlocks: HandleBlocks,
}

func HandleRequest(ctx context.Context, api PeerExchange, peer *Peer) error {
	msg, err := peer.rw.ReadMsg(ctx)
	if err != nil {
		return errors.Wrap(err, "read peer message failed")
	}

	if handlerFunc := requestHandlers[msg.Type]; handlerFunc != nil {
		_, err := handlerFunc(api, msg, peer) // TODO: discard the response, but maybe should store locally

		return err
	}

	return ErrRequestTypeNotSupported
}

func HandleResponse(ctx context.Context, api PeerExchange, peer *Peer) error {
	msg, err := peer.rw.ReadMsg(ctx)
	if err != nil {
		return errors.Wrap(err, "read peer message failed")
	}

	if handlerFunc := responseHandlers[msg.Type]; handlerFunc != nil {
		return handlerFunc(api, msg, peer)
	}

	return ErrResponseTypeNotSupported
}
