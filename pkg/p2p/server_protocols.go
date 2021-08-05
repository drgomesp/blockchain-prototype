package p2p

import (
	"context"
	"math/rand"
	"reflect"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

type msgHandlerFunc func(context.Context, MsgReadWriter) error

func (s *Server) registerProtocols(ctx context.Context) {
	streaming := func(protocolID protocol.ID, handler msgHandlerFunc) network.StreamHandler {
		return func(netStream network.Stream) {
			defer func() {
				_ = netStream.Close()
			}()

			pid := netStream.Conn().RemotePeer()
			peerInfo := s.host.Peerstore().PeerInfo(pid)

			p, err := s.connectPeer(ctx, &peerInfo, netStream)
			if err != nil {
				s.logger.Error("stream open failed", err)

				return
			}

			p.conn.SetProtocolID(protocolID)
			s.AddPeer(ctx, p)

			err = handler(ctx, p.conn.transport)
			if err != nil {
				s.logger.Error(err)
			}
		}
	}

	for _, proto := range s.protocols {
		go s.dht.Host().SetStreamHandler(protocol.ID(proto.ID), streaming(protocol.ID(proto.ID), proto.Run))
	}
}

func (s *Server) findPeerByTopic(topicName string) (peer.ID, error) {
	if len(s.topics) == 0 {
		return "", errors.New("no topic subscriptions")
	}
	topic, ok := s.topics[topicName]
	if !ok {
		return "", errors.New("no peer available")
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
		return "", errors.New("no available peers")
	}
	chosen := peers[rand.Intn(len(peers))]
	return chosen, nil
}

func stream(ctx context.Context, host host.Host, pid peer.ID, protocol protocol.ID, msg []byte) error {
	out, err := host.NewStream(ctx, pid, protocol)
	if err != nil {
		return err
	}

	defer func() {
		_ = out.Close()
	}()

	if _, err := out.Write(msg); err != nil {
		return err
	}

	return nil
}

func (s *Server) StreamMsg(ctx context.Context, msgType MsgType, msg interface{}) (err error) {
	var found peer.ID
	for tn := range s.topics {
		found, err = s.findPeerByTopic(tn)
		if found != "" {
			err = nil
			break
		}
	}

	if err != nil {
		return err
	}

	var ch codec.CborHandle
	ch.MapType = reflect.TypeOf(map[string]interface{}(nil))
	h := &ch

	var data []byte

	enc := codec.NewEncoderBytes(&data, h)
	if err := enc.Encode(msg); err != nil {
		return errors.Wrap(err, "message encode failed")
	}

	if err := stream(ctx, s.dht.Host(), found, protocol.ID(msgType), data); err != nil {
		return err
	}

	return nil
}
