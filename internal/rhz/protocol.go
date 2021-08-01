package rhz

import "github.com/drgomesp/rhizom/pkg/rhz"

// MsgType ...
type MsgType uint64

const (
	MsgNewBlock = MsgType(iota)
)

// Message that peers exchange in the network.
type Message interface {
	Type() MsgType
}

func (t MsgType) String() string {
	switch t {
	case MsgNewBlock:
		return "MsgNewBlock"
	default:
		return ""
	}
}

// NewBlock ...
type NewBlock struct {
	Block *rhz.Block
}

func (n NewBlock) Type() MsgType {
	return MsgNewBlock
}
