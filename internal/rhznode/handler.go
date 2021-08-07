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

func (h *Handler) GetBlocks(peer *rhz.Peer, msg rhz.MsgGetBlocks) (rhz.MsgBlocks, error) {
	h.logger.Debugw("handling request", "peer", peer, "msg", msg)

	return rhz.MsgBlocks{
		IsUpdated: true,
		Chain: []struct{ Header struct{ Index uint64 } }{
			{Header: struct{ Index uint64 }{Index: 123456}},
		},
	}, nil
}

func (h *Handler) Blocks(peer *rhz.Peer, msg rhz.MsgBlocks) error {
	h.logger.Debugw("received response", "msg", msg)

	return nil
}
