package rhznode

import (
	"context"

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

func (h *Handler) GetBlocks(ctx context.Context, peer *rhz.Peer, msg rhz.MsgGetBlocks) (rhz.MsgBlocks, error) {
	return rhz.MsgBlocks{
		IsUpdated: true,
		Chain: []struct{ Header struct{ Index uint64 } }{
			{Header: struct{ Index uint64 }{Index: msg.IndexNeed * 3}},
		},
	}, nil
}

func (h *Handler) Blocks(ctx context.Context, peer *rhz.Peer, msg rhz.MsgBlocks) error {
	h.logger.Infow("received response", "msg", msg)

	return nil
}
