package p2p

import (
	"context"
)

// connectBootstrapPeers connects to all bootstrap peers.
func (s *Server) connectBootstrapPeers(ctx context.Context) {
	s.logger.Debug("connecting to bootstrap peers")

	for {
		for _, addr := range s.cfg.BootstrapAddrs {
			p, err := s.connectPeerByAddr(ctx, addr)
			if err != nil {
				s.logger.Error(err)

				return
			}

			s.logger.Debug("connected to bootstrap peer", p.info.ID.ShortString())

			return
		}
	}
}

// bootstrapNetwork the network.
func (s *Server) bootstrapNetwork(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		if err := s.dht.Bootstrap(ctx); err != nil {
			s.logger.Error("failed to bootstrap network ", err)

			return
		}

		s.logger.Debug("bootstrapped network")

		return
	}
}
