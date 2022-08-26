package rhz1

import "github.com/drgomesp/acervo/pkg/p2p"

type Peer struct {
	rw p2p.MsgReadWriter
}

func NewPeer(rw p2p.MsgReadWriter) *Peer {
	return &Peer{rw}
}