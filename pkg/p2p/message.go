package p2p

import (
	"context"
	"reflect"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pb "github.com/libp2p/go-libp2p-pubsub/pb"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

type MsgType string

type Message struct {
	Type    MsgType
	Payload interface{}
}

func (m *Message) Encode() (*pubsub.Message, error) {
	var ch codec.CborHandle
	ch.MapType = reflect.TypeOf(map[string]interface{}(nil))
	h := &ch

	var data []byte

	enc := codec.NewEncoderBytes(&data, h)
	if err := enc.Encode(m.Payload); err != nil {
		return nil, errors.Wrap(err, "message encode failed")
	}

	topicName := string(m.Type)
	msg := &pubsub.Message{
		Message: &pb.Message{
			Data:  data,
			Topic: &topicName,
		},
	}

	return msg, nil
}

func (m *Message) Decode(msg *pubsub.Message) error {
	var ch codec.CborHandle
	ch.MapType = reflect.TypeOf(map[string]interface{}(nil))
	h := &ch

	dec := codec.NewDecoderBytes(msg.Data, h)

	m.Type = MsgType(*msg.Topic)

	return dec.Decode(&m.Payload)
}

// Sender of messages.
type Sender interface {
	// SendMsg ...
	SendMsg(context.Context, *Message) error
}

// Receiver of messages.
type Receiver interface {
	// ReceiveMsg ...
	ReceiveMsg(context.Context) (*Message, error)
}

// MsgPipe sends and receives messages.
type MsgPipe interface {
	Sender
	Receiver
}

func SendMsg(ctx context.Context, sender Sender, msgType MsgType, payload interface{}) error {
	msg := &Message{
		Type:    msgType,
		Payload: payload,
	}

	if err := sender.SendMsg(ctx, msg); err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}
