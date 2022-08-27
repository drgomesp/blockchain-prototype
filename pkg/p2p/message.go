package p2p

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/gogo/protobuf/proto"
	protolegacy "github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// MsgType defines a type identifier for messages.
type MsgType string

// MsgDecoder defines a decoder for messages.
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

// Message represents an encoded message that can be exchanged in the network.
type Message struct {
	Type    MsgType
	Payload io.Reader
}

func (m Message) String() string {
	return fmt.Sprintf("<%s %s>\n", m.Type, m.Payload)
}

func (m *Message) Decode(v interface{}) error {
	//var ch codec.CborHandle
	//h := &ch
	//
	//data, err := ioutil.ReadAll(m.Payload)
	//if err != nil {
	//	return errors.Wrap(err, "message payload read failed")
	//}
	//
	//dec := codec.NewDecoderBytes(data, h)
	//if err = dec.Decode(&v); err != nil {
	//	return errors.Wrap(err, "message decode failed")
	//}

	data, err := io.ReadAll(m.Payload)
	if err != nil {
		return errors.Wrap(err, "message read failed")
	}

	if err := proto.Unmarshal(data, protolegacy.MessageV1(v.(proto.Message))); err != nil {
		return err
	}

	return nil
}

// Send an encoded message through the read/writer pipe.
func Send(ctx context.Context, rw MsgReadWriter, t MsgType, msg proto.Message) error {
	// TODO: abstract encoding somewhere, but definitely not here.
	//var ch codec.CborHandle
	//h := &ch
	//
	//var out []byte
	//enc := codec.NewEncoderBytes(&out, h)
	//
	//if err := enc.Encode(msg); err != nil {
	//	return errors.Wrap(err, "message encode failed")
	//}

	log.Trace().Str("t", string(t)).Interface("msg", msg).Msgf("p2p.Send")
	out, err := proto.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "message encode failed")
	}

	if err := rw.WriteMsg(ctx, &Message{
		Type:    t,
		Payload: bytes.NewReader(out),
	}); err != nil {
		return errors.Wrap(err, "message write failed")
	}

	return nil
}
