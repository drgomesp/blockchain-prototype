package p2p

import (
	"context"
	"io"
	"reflect"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pb "github.com/libp2p/go-libp2p-pubsub/pb"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

type MessageHandler func(req MsgDecoder) (res MsgDecoder, err error)

type MsgType string

type MsgDecoder interface {
	Decode(v io.Reader) error
	Encode() (*pubsub.Message, error)
}

type Message struct {
	Type    MsgType
	Payload map[string]interface{}
}

func (m *Message) Decode(v io.Reader) error {
	var ch codec.CborHandle
	ch.MapType = reflect.TypeOf(map[string]interface{}(nil))
	h := &ch

	dec := codec.NewDecoder(v, h)

	return dec.Decode(&m.Payload)
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

func SendMsg(ctx context.Context, sender MsgWriter, msgType MsgType, v interface{}) error {
	var payload map[string]interface{}
	if err := mapstructure.Decode(v, &payload); err != nil {
		return err
	}

	msg := &Message{
		Type:    msgType,
		Payload: payload,
	}

	if err := sender.WriteMsg(ctx, msg); err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

func (s *Server) WriteMsg(ctx context.Context, msg *Message) error {
	// todo: make sure to get available peer and stream request
	topicName := string(msg.Type)
	_, topic, err := s.subscribe(ctx, topicName)
	if err != nil {
		return err
	}

	m, err := msg.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode message")
	}

	s.logger.Infow("message sent", "msg", msg, "topic", topicName)
	return topic.Publish(ctx, m.Data)
}
