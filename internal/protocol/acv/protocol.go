package acv

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"

	acv "github.com/drgomesp/acervo/internal/protocol/acv/pb"
	"github.com/drgomesp/acervo/internal/rhz"
	"github.com/drgomesp/acervo/pkg/p2p"
)

var Logger *zap.SugaredLogger

const (
	MsgTypeRequest = iota
	MsgTypeResponse
)

const (
	MsgGetFeedsReq = p2p.MsgType("GetFeedsRequest")
	MsgGetFeedsRes = p2p.MsgType("GetFeedsResponse")
)

type (
	requestHandlerFunc  func(context.Context, rhz.Peering, p2p.MsgDecoder) (p2p.ProtocolType, proto.Message, error)
	responseHandlerFunc func(context.Context, rhz.Peering, p2p.MsgDecoder) error
)

var requestHandlers = map[p2p.MsgType]requestHandlerFunc{
	MsgGetFeedsReq: func(ctx context.Context, peering rhz.Peering, msg p2p.MsgDecoder) (
		p2p.ProtocolType, proto.Message, error,
	) {
		var req acv.GetFeeds_Request
		if err := msg.Decode(&req); err != nil {
			return p2p.NilProtocol, nil, err
			log.Error().Err(err).Send()
		}

		res := &acv.GetFeeds_Response{
			Urls: []string{
				"https://drgomesp.dev/rss",
				"https://lobste.rs/rss",
			},
		}

		log.Trace().Interface("res", res).Msgf("sending response")
		return p2p.ProtocolType(MsgGetFeedsRes), res, nil
	},
}

var responseHandlers = map[p2p.MsgType]responseHandlerFunc{
	MsgGetFeedsRes: func(ctx context.Context, peering rhz.Peering, msg p2p.MsgDecoder) error {
		var res acv.GetFeeds_Response
		if err := msg.Decode(&res); err != nil {
			log.Error().Err(err).Send()
			return err
		}

		log.Trace().Interface("res", res).Msgf("receiving response")

		return nil
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
				rpid, res, err := handlerFunc(ctx, peering, msg)
				if err != nil {
					return p2p.NilProtocol, nil, err
				}

				return rpid, res, nil
			}

		case MsgTypeResponse:
			if handlerFunc := responseHandlers[msg.Type]; handlerFunc != nil {
				if err = handlerFunc(ctx, peering, msg); err != nil {
					return p2p.NilProtocol, nil, err
				}

				return p2p.NilProtocol, nil, nil
			}
		}

		return p2p.NilProtocol, nil, errors.Errorf("unsupported message type '%s'", msg.Type)
	}
}
