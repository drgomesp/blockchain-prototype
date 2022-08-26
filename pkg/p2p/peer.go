package p2p

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

type transport interface {
	MsgReadWriter
}

type connection struct {
	transport
}

// Peer is a remote peer in a p2p network.
type Peer struct {
	info   *peer.AddrInfo
	conn   *connection
	pubSub *pubsub.PubSub
}

func (p *Peer) String() string {
	return p.info.ID.ShortString()
}

// NewPeer creates a new peer from a given address info.
func NewPeer(peerInfo *peer.AddrInfo, pubsub *pubsub.PubSub) (*Peer, error) {
	p := &Peer{
		pubSub: pubsub,
		conn:   new(connection),
		info:   peerInfo,
	}

	return p, nil
}