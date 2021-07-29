package p2p

import (
	"context"
	"time"

	p2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	router "github.com/libp2p/go-libp2p-core/routing"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	secio "github.com/libp2p/go-libp2p-secio"
	yamux "github.com/libp2p/go-libp2p-yamux"
	"github.com/libp2p/go-tcp-transport"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const networkStatePeriod = 5 * time.Second

// peerChannels manages channels where peers are sent through.
type peerChannels struct {
	discovered chan *Peer // discovered peers found through Kademlia DHT
	connected  chan *Peer // connected peers in the network
}

const ServerName = "p2p.server"

// Server manages p2p connections.
type Server struct {
	cfg    Config
	logger *zap.SugaredLogger
	node   Node
	dht    *kaddht.IpfsDHT

	running        bool              // running controls the run loop
	quit           chan bool         // quit channel to receive the stop signal
	peerChan       peerChannels      // peerChan manages channel-sent peers
	peersConnected map[peer.ID]*Peer // peersConnected holds recently connected Peer nodes
}

// NewServer initializes a p2p Server from a given Config capable of managing a network.
func NewServer(ctx context.Context, logger *zap.SugaredLogger, config Config) (*Server, error) {
	srv := &Server{
		cfg:     config,
		logger:  logger,
		node:    nil,
		dht:     new(kaddht.IpfsDHT),
		running: false,
		quit:    make(chan bool),
		peerChan: peerChannels{
			discovered: make(chan *Peer),
			connected:  make(chan *Peer),
		},
		peersConnected: make(map[peer.ID]*Peer),
	}

	if err := srv.setupLocalHost(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to initialize node")
	}

	if err := srv.setupDiscovery(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to initialize discovery")
	}

	return srv, nil
}

// HandlePeerFound receives a discovered peer.
func (s *Server) HandlePeerFound(peerInfo peer.AddrInfo) {
	p, err := NewPeer(peerInfo)
	if err != nil {
		s.logger.Error("failed to initialize peer: ", err)
	}

	s.peerChan.discovered <- p
}

func (s *Server) Name() string {
	return ServerName
}

func (s *Server) Start(ctx context.Context) error {
	go s.discover(ctx)
	go s.ping(ctx)
	go s.run(ctx)

	s.running = true

	s.connectBootstrapPeers(ctx)
	s.bootstrapNetwork(ctx)
	s.setupPeerSubscriptions(ctx)

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.running = false

	return nil
}

// setupLocalHost sets up the local p2p node.
func (s *Server) setupLocalHost(ctx context.Context) error {
	h, err := p2p.New(
		ctx,
		p2p.ChainOptions(
			p2p.Transport(tcp.NewTCPTransport),
		),
		p2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/0",
		),
		p2p.ChainOptions(
			p2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		),
		p2p.Security(secio.ID, secio.New),
		p2p.Routing(func(h host.Host) (router.PeerRouting, error) {
			var err error
			s.dht, err = kaddht.New(ctx, h)
			if err != nil {
				return nil, errors.Wrap(err, "failed to initialize dht")
			}

			return s.dht, nil
		}),
	)
	if err != nil {
		return errors.Wrap(err, "local node setup failed")
	}

	s.node = h

	return nil
}

// ping connected peers regularly.
func (s *Server) ping(ctx context.Context) {
	for {
		for _, p := range s.peersConnected {
			if err := s.node.Connect(ctx, p.Info); err != nil {
				s.RemovePeer(p)

				break
			}

			s.logger.Debug("peer connection check ", p)
		}

		time.Sleep(s.cfg.PingTimeout)
	}
}

// run the main server loop.
func (s *Server) run(_ context.Context) {
running:
	for s.running {
		select {
		case <-s.quit:
			break running
		case p := <-s.peerChan.connected:
			{
				s.peersConnected[p.Info.ID] = p
				s.logger.Info("peer added ", p)
			}
		case <-time.After(networkStatePeriod):
			{
				s.logger.Debugw("online", "connected", len(s.peersConnected))
			}
		}
	}
}

// AddPeer adds a peer to the network.
func (s *Server) AddPeer(ctx context.Context, peer *Peer) {
	for {
		var err error

		_, isConnected := s.peersConnected[peer.Info.ID]
		if isConnected {
			return
		}

		if err = s.dht.Host().Connect(ctx, peer.Info); err != nil {
			s.logger.Warnw("couldn't connect to peer", "err", err)
			continue
		}

		var p *Peer
		if p, err = NewPeer(peer.Info); err != nil {
			s.logger.Error("failed to initialize peer: ", err)
			continue
		}

		s.peerChan.connected <- p
		break
	}
}

// RemovePeer removes a peer from the network.
func (s *Server) RemovePeer(p *Peer) {
	delete(s.peersConnected, p.Info.ID)
	s.logger.Info("peer removed ", p)
}
