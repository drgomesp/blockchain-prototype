package p2p

import (
	"context"
	"io/ioutil"
	"math/rand"
	"reflect"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

func (s *Server) registerProtocols(ctx context.Context) {
	streaming := func(handler ProtoRunFunc) network.StreamHandler {
		return func(stream network.Stream) {
			data, err := ioutil.ReadAll(stream)
			if err != nil {
				s.logger.Error(err)

				return
			}

			if err := handler(data); err != nil {
				s.logger.Error(err)
			}
		}
	}

	for _, proto := range s.protocols {
		go s.dht.Host().SetStreamHandler(protocol.ID(proto.ID), streaming(proto.Run))
	}
}

func (s *Server) StreamMsg(ctx context.Context, msgType MsgType, msg interface{}) error {
	topic, ok := s.topics["/rhz/blk/default_2b678c95-27d5-4f09-bf38-a62be2c5339b"]
	if !ok {
		return errors.New("fuck")
	}

	removeMyself := func(peers []peer.ID) []peer.ID {
		me := s.dht.Host().ID()
		for i, p := range peers {
			if p == me {
				lastIndex := len(peers) - 1
				peers[i], peers = peers[lastIndex], peers[:lastIndex]
			}
		}

		return peers
	}

	peers := removeMyself(topic.ListPeers())
	if len(peers) == 0 {
		return errors.New("no available peers")
	}
	chosen := peers[rand.Intn(len(peers))]

	var ch codec.CborHandle
	ch.MapType = reflect.TypeOf(map[string]interface{}(nil))
	h := &ch

	var data []byte

	enc := codec.NewEncoderBytes(&data, h)
	if err := enc.Encode(msg); err != nil {
		return errors.Wrap(err, "message encode failed")
	}

	if err := stream(ctx, s.dht.Host(), chosen, protocol.ID(msgType), data); err != nil {
		return err
	}

	return nil
}

func stream(ctx context.Context, host host.Host, pid peer.ID, protocol protocol.ID, msg []byte) error {
	out, err := host.NewStream(ctx, pid, protocol)
	if err != nil {
		return err
	}

	defer out.Close()

	if _, err := out.Write(msg); err != nil {
		return err
	}

	return nil
}
