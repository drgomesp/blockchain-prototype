package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/rs/zerolog/log"
)

// HandlePeerFound receives a discovered peer.
func (s *Server) HandlePeerFound(peerInfo peer.AddrInfo) {
	p, err := NewPeer(&peerInfo, s.pubSub)
	if err != nil {
		log.Error().Err(err)

		return
	}

	s.peerChan.discovered <- p
}

// setupDiscovery sets up the peer discovery mechanism.
func (s *Server) setupDiscovery() error {
	disc := mdns.NewMdnsService(s.host, s.cfg.ServiceTag, s)
	s.disc = disc
	return nil
}

// discover for incoming discovered peers.
func (s *Server) discover(ctx context.Context) {
listening:
	for {
		select {
		case <-ctx.Done():
			{
				s.quit <- true

				break listening
			}
		case p := <-s.peerChan.discovered:
			{
				if s.PeerDiscovered(p.info) {
					continue
				}

				s.peersDiscovered[p.String()] = p

				if !s.PeerConnected(p) {
					log.Debug().Msgf("peer discovered ", p)
					s.AddPeer(ctx, p)
				}
			}
		}
	}
}

// PeerDiscovered checks if the peer is discovered by the network.
func (s *Server) PeerDiscovered(peerInfo *peer.AddrInfo) bool {
	_, is := s.peersDiscovered[peerInfo.String()]

	return is
}