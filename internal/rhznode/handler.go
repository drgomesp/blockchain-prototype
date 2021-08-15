package rhznode

import (
	"context"

	rhz2 "github.com/drgomesp/rhizom/internal/protocol/rhz1"
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

func (h *Handler) GetBlocks(ctx context.Context, peering rhz2.Peering, msg rhz2.MsgGetBlocks) (rhz2.MsgBlocks, error) {
	return rhz2.MsgBlocks{
		IsUpdated: true,
		Chain: []struct{ Header struct{ Index uint64 } }{
			{Header: struct{ Index uint64 }{Index: msg.IndexNeed * 2}},
		},
	}, nil
}

func (h *Handler) Blocks(ctx context.Context, peering rhz2.Peering, msg rhz2.MsgBlocks) error {
	h.logger.Infow("", "response", msg)

	return nil
}

func (h *Handler) OnNewBlock(ctx context.Context, msg rhz2.MsgNewBlock) error {
	h.logger.Infow("OnNewBlock", "msg", msg)

	return nil
}
