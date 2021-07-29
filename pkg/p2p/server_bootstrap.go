package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

// connectBootstrapPeers connects to all bootstrap peers.
func (s *Server) connectBootstrapPeers(ctx context.Context) {
	s.logger.Info("connecting to bootstrap peers")

	connected := make([]*Peer, len(s.cfg.BootstrapAddrs))

bootstrap:
	for {
		for i, addr := range s.cfg.BootstrapAddrs {
			peerAddr, err := multiaddr.NewMultiaddr(addr)
			if err != nil {
				s.logger.Error("failed to initialize multiaddr: ", err)

				continue
			}

			peerInfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
			if err != nil {
				s.logger.Error("failed to load addr info from multiaddr: ", err)

				continue
			}

			var p *Peer
			if p, err = NewPeer(*peerInfo); err != nil {
				s.logger.Error("failed to initialize peer: ", err)
			}

			s.AddPeer(ctx, p)
			connected[i] = p
		}

		for _, peerInfo := range connected {
			if _, ok := s.peersConnected[peerInfo.Info.ID]; !ok {
				continue bootstrap
			}
		}

		break bootstrap
	}

	s.logger.Info("connected to bootstrap peers")
}

// bootstrapNetwork the network.
func (s *Server) bootstrapNetwork(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		if err := s.dht.Bootstrap(ctx); err != nil {
			s.logger.Error("failed to bootstrap network ", err)
		}

		s.logger.Info("bootstrapped network")

		return
	}
}
