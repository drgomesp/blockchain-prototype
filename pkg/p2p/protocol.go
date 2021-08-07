package p2p

import (
	"bytes"
	"context"
	"io/ioutil"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

type ProtocolType string

var NilProtocol = ProtocolType("")

type StreamHandlerFunc func(context.Context, MsgReadWriter) (ProtocolType, interface{}, error)

// Protocol defines a sub-protocol for communication in the network.
type Protocol struct {
	// ID is the unique identifier of the protocol (three-letter word).
	ID string

	// Run ...
	Run StreamHandlerFunc
}

type protoRW struct {
	pid      peer.ID
	host     host.Host
	read     network.Stream
	write    network.Stream
	writePID ProtocolType
}

func (p *protoRW) ReadMsg(ctx context.Context) (*Message, error) {
	data, err := ioutil.ReadAll(p.read)
	if err != nil {
		return nil, err
	}

	return &Message{
		Type:    MsgType(p.read.Protocol()),
		Payload: bytes.NewReader(data),
	}, nil
}

func (p *protoRW) WriteMsg(ctx context.Context, msg *Message) error {
	out, err := p.host.NewStream(ctx, p.pid, protocol.ID(msg.Type))
	if err != nil {
		return err
	}

	defer func() {
		_ = out.Close()
	}()

	data, err := ioutil.ReadAll(msg.Payload)
	if err != nil {
		return err
	}

	if _, err := out.Write(data); err != nil {
		return err
	}

	p.write = out

	return nil
}
