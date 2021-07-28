package p2p

import (
	"github.com/libp2p/go-libp2p-core/peer"
)

// HandlePeerFound Receive a peer info in an channel.
func (s *Server) HandlePeerFound(peerInfo peer.AddrInfo) {
	s.peerChan.discovered <- peerInfo
}
