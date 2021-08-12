package rhz

import (
	"github.com/drgomesp/rhizom/pkg/block"
	"github.com/drgomesp/rhizom/pkg/p2p"
)

type MsgGetBlocks struct {
	IndexHave uint64
	IndexNeed uint64
}

type MsgBlocks struct {
	IsUpdated bool
	Chain     []struct {
		Header struct {
			Index uint64
		}
	}
}

type MsgNewBlock struct {
	Block block.Block
}

func (g MsgGetBlocks) Type() p2p.MsgType {
	return MsgTypeGetBlocks
}

func (g MsgBlocks) Type() p2p.MsgType {
	return MsgTypeBlocks
}

func (m MsgNewBlock) Type() p2p.MsgType {
	return MsgTypeNewBlock
}
