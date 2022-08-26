package p2p

import (
	"context"

	"github.com/rs/zerolog/log"
)

// connectBootstrapPeers connects to all bootstrap peers.
func (s *Server) connectBootstrapPeers(ctx context.Context) {
	log.Debug().Msg("connecting to bootstrap peers")

	for {
		for _, addr := range s.cfg.BootstrapAddrs {
			_, err := s.connectPeerByAddr(ctx, addr)
			if err != nil {
				log.Error().Err(err).Send()
				return
			}

			log.Debug().Msgf("connected to bootstrap peer: %s", addr)

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
			log.Error().Msgf("failed to bootstrap network ", err)

			return
		}

		log.Debug().Msgf("bootstrapped network")

		return
	}
}