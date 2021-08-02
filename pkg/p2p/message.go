package p2p

import (
	"reflect"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/ugorji/go/codec"
)

type MsgType string

type Message struct {
	Type    MsgType
	Payload interface{}
}

func (m *Message) Decode(v interface{}) error {
	var ch codec.CborHandle
	ch.MapType = reflect.TypeOf(map[string]interface{}(nil))
	h := &ch

	msg := v.(*pubsub.Message)
	m.Type = MsgType(*msg.Topic)
	m.Payload = new(interface{})

	dec := codec.NewDecoderBytes(msg.Data, h)

	return dec.Decode(&m.Payload)
}
