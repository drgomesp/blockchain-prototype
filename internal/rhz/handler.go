package rhz

import "go.uber.org/zap"

type NetworkHandler struct {
	logger *zap.SugaredLogger
}

func NewNetworkHandler(logger *zap.SugaredLogger) *NetworkHandler {
	return &NetworkHandler{
		logger: logger,
	}
}

func (h *NetworkHandler) HandleRequest(peer *Peer, msg MessagePacket) (MessagePacket, error) {
	h.logger.Debugw("handling request", "peer", peer, "msg", msg)

	return nil, nil
}

func (h *NetworkHandler) HandleResponse(peer *Peer, msg MessagePacket) error {
	h.logger.Debugw("handling response", "peer", peer, "msg", msg)

	return nil
}
