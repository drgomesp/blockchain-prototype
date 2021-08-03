package rhz

import "github.com/drgomesp/rhizom/pkg/p2p"

type Peer struct {
	*p2p.Peer
}

func NewPeer() *Peer {
	return &Peer{}
}
