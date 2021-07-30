package p2p

import (
	"github.com/libp2p/go-libp2p-core/peer"
)

// Peer is a remote peer in a p2p network.
type Peer struct {
	info peer.AddrInfo
}

func (p *Peer) Info() peer.AddrInfo {
	return p.info
}

func (p *Peer) String() string {
	return p.info.ID.ShortString()
}

// NewPeer creates a new peer from a given address info.
func NewPeer(peerInfo peer.AddrInfo) (*Peer, error) {
	p := &Peer{
		info: peerInfo,
	}

	return p, nil
}
