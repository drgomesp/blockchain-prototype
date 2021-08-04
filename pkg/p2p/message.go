package p2p

import (
	"context"
	"io"

	"github.com/ugorji/go/codec"
)

type MessageHandler func(req MsgDecoder) (res MsgDecoder, err error)

type MsgType string

type MsgDecoder interface {
	Decode(v interface{}) error
}

type Message struct {
	Type    MsgType
	Payload io.Reader
}

func (m *Message) Decode(v interface{}) error {
	var ch codec.CborHandle
	h := &ch

	dec := codec.NewDecoder(m.Payload, h)

	return dec.Decode(&v)
}

// MsgWriter of messages.
type MsgWriter interface {
	WriteMsg(context.Context, *Message) error
}

// MsgReader of messages.
type MsgReader interface {
	ReadMsg(context.Context) (*Message, error)
}

// MsgReadWriter sends and receives messages.
type MsgReadWriter interface {
	MsgWriter
	MsgReader
}
