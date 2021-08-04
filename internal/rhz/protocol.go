package rhz

import (
	"context"
	"log"

	"github.com/drgomesp/rhizom/pkg/p2p"
)

// request/response messages.
const (
	// MsgGetBlocksRequest represents a request for blocks.
	MsgGetBlocksRequest = p2p.MsgType("/rhz/blocks/req/default_2b678c95-27d5-4f09-bf38-a62be2c5339b")
	// MsgGetBlocksResponse represents a response to a request for blocks.
	MsgGetBlocksResponse = p2p.MsgType("/rhz/blocks/resp/default_2b678c95-27d5-4f09-bf38-a62be2c5339b")
)

// async messages.
const (
	// MsgNewBlock represents a new block.
	MsgNewBlock = p2p.MsgType("NewBlock")
)

type (
	RequestMessageHandler func(PeerExchange, Message, *Peer) (MessagePacket, error)
	ResponseMsgHandler    func(PeerExchange, Message, *Peer) error
)

var requestHandlers = map[p2p.MsgType]RequestMessageHandler{
	MsgGetBlocksRequest: HandleGetBlocksRequest,
}

var responseHandlers = map[p2p.MsgType]ResponseMsgHandler{
	MsgGetBlocksResponse: HandleGetBlocksResponse,
}

func HandleRequestMsg(ctx context.Context, api PeerExchange, peer *Peer) error {
	msg, err := peer.rw.ReadMsg(ctx)
	if err != nil {
		return err
	}

	if handlerFunc := requestHandlers[msg.Type]; handlerFunc != nil {
		res, err := handlerFunc(api, msg, peer)
		if err != nil {
			return err
		}

		log.Println(res)
	}

	return nil
}

func HandleResponseMsg(ctx context.Context, api PeerExchange, peer *Peer) error {
	msg, err := peer.rw.ReadMsg(ctx)
	if err != nil {
		return err
	}

	if handlerFunc := responseHandlers[msg.Type]; handlerFunc != nil {
		err := handlerFunc(api, msg, peer)
		if err != nil {
			return err
		}
	}

	return nil
}
