package p2p

import (
	"bytes"
	"context"
	"io/ioutil"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

type ProtocolType string

var NilProtocol = ProtocolType("")

type StreamHandlerFunc func(context.Context, network.Stream) (ProtocolType, MsgDecoder, error)

// Protocol defines a sub-protocol for communication in the network.
type Protocol struct {
	// ID is the unique identifier of the protocol (three-letter word).
	ID string

	// Run ...
	Run func(context.Context, network.Stream) (ProtocolType, MsgDecoder, error)
}

type streamRW struct {
	in  network.Stream
	out network.Stream
}

func (s *streamRW) WriteMsg(ctx context.Context, msg *Message) error {
	var ch codec.CborHandle
	h := &ch

	var data []byte

	enc := codec.NewEncoderBytes(&data, h)
	if err := enc.Encode(msg); err != nil {
		return errors.Wrap(err, "message encode failed")
	}

	if _, err := s.out.Write(data); err != nil {
		return err
	}

	return nil
}

func (s *streamRW) ReadMsg(ctx context.Context) (*Message, error) {
	data, err := ioutil.ReadAll(s.in)
	if err != nil {
		return nil, err
	}

	return &Message{
		Type:    MsgType(s.in.Protocol()),
		Payload: bytes.NewReader(data),
	}, nil
}
