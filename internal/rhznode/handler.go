package rhznode

import (
	"github.com/drgomesp/rhizom/internal/rhz"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.SugaredLogger
}

func NewHandler(logger *zap.SugaredLogger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) ReceiveRequest(peer *rhz.Peer, msg rhz.MessagePacket) (rhz.MessagePacket, error) {
	h.logger.Debugw("handling request", "peer", peer, "msg", msg)

	return &rhz.MsgBlocks{
		IsUpdated: false,
		Chain: []struct{ Header struct{ Index uint64 } }{
			{Header: struct{ Index uint64 }{Index: 99999999}},
		},
	}, nil
}

func (h *Handler) ReceiveResponse(peer *rhz.Peer, msg rhz.MessagePacket) error {
	h.logger.Debugw("handling response", "peer", peer, "msg", msg)

	return nil
}
