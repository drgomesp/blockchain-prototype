package p2p

import (
	"bytes"
	"context"
	"io"
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

var ErrNoPeersFound = errors.New("no peers found")

func (s *Server) registerProtocols(ctx context.Context) {
	streamHandler := func(protocolID protocol.ID, handler StreamHandlerFunc) network.StreamHandler {
		return func(netStream network.Stream) {
			defer func() {
				_ = netStream.Close()
			}()

			pid := netStream.Conn().RemotePeer()
			peerInfo := s.host.Peerstore().PeerInfo(pid)

			p, err := s.setupProtocolConnection(ctx, &peerInfo, netStream)
			if err != nil {
				s.logger.Error("stream open failed", err)

				return
			}

			s.AddPeer(ctx, p)

			rpid, msg, err := handler(ctx, netStream)
			if err != nil {
				s.logger.Error(err)
				return
			}

			// early return if we are handling a response, which needs no communicating back
			if rpid == NilProtocol {
				return
			}

			var ch codec.CborHandle
			h := &ch

			var data []byte

			enc := codec.NewEncoderBytes(&data, h)
			if err := enc.Encode(msg); err != nil {
				s.logger.Error(err)
				return
			}

			if err := stream(ctx, s.host, pid, protocol.ID(rpid), bytes.NewReader(data)); err != nil {
				s.logger.Error(err)
				return
			}
		}
	}

	for _, proto := range s.protocols {
		pid := protocol.ID(proto.ID)
		go s.dht.Host().SetStreamHandler(pid, streamHandler(pid, proto.Run))
	}
}

func (s *Server) findPeerByTopic(topicName string) (peer.ID, error) {
	if len(s.topics) == 0 {
		return "", errors.New("no topic subscriptions")
	}
	topic, ok := s.topics[topicName]
	if !ok {
		return "", ErrNoPeersFound
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
		return "", ErrNoPeersFound
	}
	chosen := peers[rand.Intn(len(peers))]
	return chosen, nil
}

func stream(ctx context.Context, host host.Host, pid peer.ID, protocol protocol.ID, msg io.Reader) error {
	out, err := host.NewStream(ctx, pid, protocol)
	if err != nil {
		return err
	}

	defer func() {
		_ = out.Close()
	}()

	data, err := ioutil.ReadAll(msg)
	if err != nil {
		return err
	}

	if _, err := out.Write(data); err != nil {
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

	if err := stream(ctx, s.dht.Host(), found, protocol.ID(msgType), bytes.NewReader(data)); err != nil {
		return err
	}

	return nil
}
