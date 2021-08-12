package p2p

//
//import (
//	"context"
//	"reflect"
//
//	"github.com/libp2p/go-libp2p-core/network"
//	"github.com/libp2p/go-libp2p-core/protocol"
//	"github.com/pkg/errors"
//	"github.com/ugorji/go/codec"
//)
//
//type streamTransport struct {
//	protocolID protocol.ID
//	stream     network.Stream
//}
//
//func (t *streamTransport) WriteMsg(ctx context.Context, msg Message) error {
//	defer func() {
//		_ = t.stream.Close()
//	}()
//
//	var ch codec.CborHandle
//	ch.MapType = reflect.TypeOf(map[string]interface{}(nil))
//	h := &ch
//
//	var data []byte
//
//	enc := codec.NewEncoderBytes(&data, h)
//	if err := enc.Encode(msg.Payload); err != nil {
//		return errors.Wrap(err, "message encode failed")
//	}
//
//	if _, err := t.stream.Write(data); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (t *streamTransport) ReadMsg(ctx context.Context) (Message, error) {
//	if t.stream == nil {
//		return NilMessage, errors.New("stream is empty")
//	}
//
//	var msg Message
//	msg.Payload = t.stream
//	msg.Type = MsgType(t.protocolID)
//
//	return msg, nil
//}
