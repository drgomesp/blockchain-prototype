package protocol

import (
	"bytes"
	"context"

	"github.com/drgomesp/rhizom/pkg/p2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/ugorji/go/codec"
)

const (
	Ping = "/rhz/ping/1.0.0"
)

func PingHandler(ctx context.Context, stream network.Stream) (p2p.ProtocolType, interface{}, error) {
	var ch codec.CborHandle
	h := &ch

	var data []byte

	enc := codec.NewEncoderBytes(&data, h)
	if err := enc.Encode([]byte("PONG")); err != nil {
		return p2p.NilProtocol, nil, err
	}

	msg := &p2p.Message{
		Type:    Ping,
		Payload: bytes.NewReader(data),
	}

	return Ping, msg, nil
}
