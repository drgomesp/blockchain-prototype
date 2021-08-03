package rhz

import (
	"context"

	"github.com/drgomesp/rhizom/pkg/p2p"
)

type MsgHandlerFunc func(msg interface{}) error

var handlers = map[p2p.MsgType]MsgHandlerFunc{
	MsgGetBlocksRequest:  HandleGetBlocksRequest,
	MsgGetBlocksResponse: HandleGetBlocksResponse,
}

func HandleMessage(ctx context.Context, rw p2p.MsgReadWriter) error {
	msg, err := rw.ReadMsg(ctx)
	if err != nil {
		return err
	}

	if handlerFunc := handlers[msg.Type]; handlerFunc != nil {
		return handlerFunc(msg)
	}

	return nil
}
