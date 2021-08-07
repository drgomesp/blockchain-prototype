package rhz

import "github.com/drgomesp/rhizom/pkg/p2p"

// Message used for direct p2p communication.
type Message interface {
	// Decode ...
	Decode(v interface{}) error
}

type MessagePacket interface {
	Type() p2p.MsgType
}

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

func (g MsgGetBlocks) Type() p2p.MsgType {
	return MsgTypeGetBlocks
}

func (g MsgBlocks) Type() p2p.MsgType {
	return MsgTypeBlocks
}
