package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

// connectBootstrapPeers connects to all bootstrap peers.
func (s *Server) connectBootstrapPeers(ctx context.Context) {
	s.logger.Debug("connecting to bootstrap peers")

	for {
		for _, addr := range s.cfg.BootstrapAddrs {
			peerAddr, err := multiaddr.NewMultiaddr(addr)
			if err != nil {
				s.logger.Error("failed to initialize multiaddr", err)

				continue
			}

			peerInfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
			if err != nil {
				s.logger.Error("failed to load addr info from multiaddr", err)

				continue
			}

			if err := s.dht.Host().Connect(ctx, *peerInfo); err != nil {
				s.logger.Error("failed to connect to bootstrap peer", err)

				continue
			}

			s.logger.Debug("connected to bootstrap peer", peerInfo.ID.ShortString())

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
