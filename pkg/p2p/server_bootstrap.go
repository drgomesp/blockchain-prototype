package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
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

func (s *Server) openStreamWithPeer(ctx context.Context, peerID peer.ID, pid protocol.ID) (network.Stream, error) {
	str, err := s.dht.Host().NewStream(ctx, peerID, pid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open stream with peer")
	}

	return str, nil
}

func (s *Server) connectPeerByAddr(ctx context.Context, addr string) (*Peer, error) {
	peerAddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize multiaddr")
	}

	peerInfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load addr info from multiaddr")
	}

	return s.connectPeer(ctx, peerInfo, nil)
}

func (s *Server) connectPeer(ctx context.Context, peerInfo *peer.AddrInfo, stream network.Stream) (*Peer, error) {
	if err := s.dht.Host().Connect(ctx, *peerInfo); err != nil {
		return nil, errors.Wrap(err, "failed to connect to bootstrap peer")
	}

	return &Peer{
		info:   peerInfo,
		conn:   &connection{transport: &streamTransport{stream: stream}},
		pubSub: s.pubSub,
	}, nil
}
