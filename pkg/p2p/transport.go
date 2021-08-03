package p2p

import (
	"bytes"
	"context"
	"io/ioutil"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/pkg/errors"
)

type streamTransport struct {
	protocolID protocol.ID
	stream     network.Stream
}

func (t *streamTransport) WriteMsg(ctx context.Context, msg *Message) error {
	defer func() {
		_ = t.stream.Close()
	}()

	encoded, err := msg.Encode()
	if err != nil {
		return err
	}

	if _, err := t.stream.Write(encoded.Data); err != nil {
		return err
	}

	return nil
}

func (t *streamTransport) ReadMsg(ctx context.Context) (*Message, error) {
	if t.stream == nil {
		return nil, errors.New("stream is empty")
	}

	data, err := ioutil.ReadAll(t.stream)
	if err != nil {
		return nil, errors.Wrap(err, "stream message read failed")
	}

	var msg Message
	if err := msg.Decode(bytes.NewReader(data)); err != nil {
		return nil, err
	}
	msg.Type = MsgType(t.protocolID)

	return &msg, nil
}

func (t *streamTransport) SetProtocolID(pid protocol.ID) {
	t.protocolID = pid
}
