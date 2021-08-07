package p2p

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/ugorji/go/codec"
)

type MessageHandler func(req MsgDecoder) (res MsgDecoder, err error)

type MsgType string

type MsgDecoder interface {
	Decode(v interface{}) error
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

var NilMessage = Message{}

type Message struct {
	Type    MsgType
	Payload io.Reader
}

func (m Message) String() string {
	return fmt.Sprintf("<%s %s>\n", m.Type, m.Payload)
}

func (m *Message) Decode(v interface{}) error {
	var ch codec.CborHandle
	h := &ch
	data, err := ioutil.ReadAll(m.Payload)
	if err != nil {
		return err
	}

	dec := codec.NewDecoderBytes(data, h)

	return dec.Decode(&v)
}

func Send(ctx context.Context, rw MsgReadWriter, t MsgType, payload interface{}) error {
	var ch codec.CborHandle
	h := &ch

	var data []byte
	enc := codec.NewEncoderBytes(&data, h)

	if err := enc.Encode(payload); err != nil {
		return err
	}

	return rw.WriteMsg(ctx, &Message{
		Type:    t,
		Payload: bytes.NewReader(data),
	})
}
