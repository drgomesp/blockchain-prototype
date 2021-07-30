package p2p

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	"github.com/pkg/errors"
)

// serviceTag is an identifier for the discovery service.
const serviceTag = "rhizom"

// HandlePeerFound receives a discovered peer.
func (s *Server) HandlePeerFound(peerInfo peer.AddrInfo) {
	p, err := NewPeer(peerInfo)
	if err != nil {
		s.logger.Error("failed to initialize peer: ", err)
	}

	s.peerChan.discovered <- p
}

// setupDiscovery sets up the peer discovery mechanism.
func (s *Server) setupDiscovery(ctx context.Context) error {
	disc, err := discovery.NewMdnsService(ctx, s.host, time.Second, serviceTag)
	if err != nil {
		return errors.Wrap(err, "failed to initialize disc")
	}

	disc.RegisterNotifee(s)

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
				if s.PeerDiscovered(p.Info()) {
					continue
				}

				s.peersDiscovered[p.Info().ID] = p

				if !s.PeerConnected(p) {
					s.logger.Debug("peer discovered ", p)
					go s.AddPeer(ctx, p)
				}
			}
		}
	}
}
