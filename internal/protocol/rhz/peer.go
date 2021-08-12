package rhz

import "github.com/drgomesp/rhizom/pkg/p2p"

type Peer struct {
	rw p2p.MsgReadWriter
}

func NewPeer(rw p2p.MsgReadWriter) *Peer {
	return &Peer{rw}
}
