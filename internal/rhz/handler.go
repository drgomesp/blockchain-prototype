package rhz

import "go.uber.org/zap"

type Handler struct {
	logger *zap.SugaredLogger
}

func NewHandler(logger *zap.SugaredLogger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) HandleRequest(peer *Peer, msg MessagePacket) (MessagePacket, error) {
	h.logger.Debugw("handling request", "peer", peer, "msg", msg)

	return nil, nil
}

func (h *Handler) HandleResponse(peer *Peer, msg MessagePacket) error {
	h.logger.Debugw("handling response", "peer", peer, "msg", msg)

	return nil
}
