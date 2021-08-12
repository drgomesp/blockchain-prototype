package p2p

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/pkg/errors"
)

// ProtocolType defines the type for a user-defined sub-protocol.
type ProtocolType string

var NilProtocol = ProtocolType("")

// StreamHandlerFunc defines the sub-protocol handler function to handle incoming or outgoing messages.
type StreamHandlerFunc func(context.Context, MsgReadWriter) (ProtocolType, interface{}, error)

// Protocol defines a sub-protocol for communication in the network.
type Protocol struct {
	// ID is the unique identifier of the protocol (three-letter word).
	ID string

	// Run is called in a separate go routine for every p2p stream connection opened.
	Run StreamHandlerFunc
}

// protoRW is a read/write pipe used internally for protocol messaging.
type protoRW struct {
	pid                         peer.ID
	host                        host.Host
	readProtocol, writeProtocol ProtocolType
	read                        io.Reader
	write                       io.Writer
}

// ReadMsg ...
func (p *protoRW) ReadMsg(ctx context.Context) (*Message, error) {
	data, err := ioutil.ReadAll(p.read)
	if err != nil {
		return nil, err
	}

	return &Message{
		Type:    MsgType(p.readProtocol),
		Payload: bytes.NewReader(data),
	}, nil
}

// WriteMsg ...
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
		return errors.Wrap(err, "message write failed")
	}

	p.write = out

	return nil
}
