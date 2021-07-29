package p2p

import (
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
)

// Peer is a remote peer connected to the network.
type Peer struct {
	Info peer.AddrInfo
	node host.Host
}

func (p Peer) String() string {
	return p.Info.ID.ShortString()
}

func NewPeer(peerInfo peer.AddrInfo) (*Peer, error) {
	p := &Peer{
		Info: peerInfo,
	}

	return p, nil
}
