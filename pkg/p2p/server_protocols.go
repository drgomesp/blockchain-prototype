package p2p

import (
	"bufio"
	"context"
	"io"
	"io/ioutil"
	"math/rand"

	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/pkg/errors"
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

			p := s.setupProtocolConnection(ctx, &peerInfo)
			go s.AddPeer(ctx, p)

			rw := &protoRW{
				pid:          pid,
				host:         s.host,
				readProtocol: ProtocolType(protocolID),
				read:         bufio.NewReader(netStream),
				write:        bufio.NewWriter(netStream),
			}

			rpid, msg, err := handler(ctx, rw)
			if err != nil {
				s.logger.Error(err)

				return
			}

			rw.writeProtocol = rpid

			// early return if we are handling a response, which needs no communicating back
			if rpid == NilProtocol {
				return
			}

			if err = Send(ctx, rw, MsgType(rpid), msg.(proto.Message)); err != nil {
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

// WriteMsg ...
func (s *Server) WriteMsg(ctx context.Context, msg *Message) (err error) {
	var peerFound peer.ID

	for topicName := range s.topics {
		if peerFound, err = s.findPeerByTopic(topicName); err != nil {
			if !errors.Is(err, ErrNoPeersFound) {
				return errors.Wrap(err, "find peer by topic failed")
			}
		}
	}

	if peerFound == "" {
		return ErrNoPeersFound
	}

	return stream(ctx, s.dht.Host(), peerFound, protocol.ID(msg.Type), msg.Payload)
}

// ReadMsg ..
func (s *Server) ReadMsg(ctx context.Context) (*Message, error) {
	// TODO
	panic("implement me")
}

// findPeerByTopic finds a pseudo-random peer by a given topic, excluding myself.
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

	// TODO: maybe change this pseudo-random to a more proper random (crypto/rand).
	return peers[rand.Intn(len(peers))], nil
}

// stream opens a new stream and sends a message to a given peer through some protocol.
func stream(ctx context.Context, host host.Host, pid peer.ID, protoID protocol.ID, msg io.Reader) error {
	out, err := host.NewStream(ctx, pid, protoID)
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
