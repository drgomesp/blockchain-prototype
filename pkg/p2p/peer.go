package p2p

import (
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Peer is a remote peer in a p2p network.
type Peer struct {
	pubSub *pubsub.PubSub
	info   peer.AddrInfo
}

func (p *Peer) Info() peer.AddrInfo {
	return p.info
}

func (p *Peer) String() string {
	return p.info.ID.ShortString()
}

// NewPeer creates a new peer from a given address info.
func NewPeer(pubsub *pubsub.PubSub, peerInfo peer.AddrInfo) (*Peer, error) {
	p := &Peer{
		pubSub: pubsub,
		info:   peerInfo,
	}

	return p, nil
}

func (p *Peer) ReadMsg() (Message, error) {
	panic("implement me")
}

func (p *Peer) WriteMsg(msg Message) error {
	panic("implement me")
}
