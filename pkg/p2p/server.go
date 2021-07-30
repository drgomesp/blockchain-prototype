package p2p

import (
	"context"
	"time"

	p2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	router "github.com/libp2p/go-libp2p-core/routing"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	mplex "github.com/libp2p/go-libp2p-mplex"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	secio "github.com/libp2p/go-libp2p-secio"
	yamux "github.com/libp2p/go-libp2p-yamux"
	"github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const networkStatePeriod = 5 * time.Second

// peerChannels manages channels where peers are sent through.
type peerChannels struct {
	discovered chan *Peer // discovered peers found through Kademlia DHT
	connected  chan *Peer // connected peers in the network
}

const ServerName = "p2p"

// Server manages p2p connections.
type Server struct {
	cfg            Config             // cfg server options.
	logger         *zap.SugaredLogger // logger provided logger.
	peer           *Peer              // Peer is the local p2p peer.
	host           host.Host
	dht            *kaddht.IpfsDHT
	pubSub         *pubsub.PubSub
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
		peer:    nil,
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
		return nil, errors.Wrap(err, "failed to initialize peer")
	}

	if err := srv.setupDiscovery(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to initialize discovery")
	}

	if err := srv.setupPubSub(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to initialize pubsub")
	}

	return srv, nil
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
	s.setupSubscriptions(ctx)

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.running = false

	return nil
}

// setupLocalHost sets up the local p2p peer.
func (s *Server) setupLocalHost(ctx context.Context) error {
	h, err := p2p.New(
		ctx,
		p2p.ChainOptions(
			p2p.Transport(tcp.NewTCPTransport),
			p2p.Transport(ws.New),
		),
		p2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/0",
			"/ip4/0.0.0.0/tcp/0/ws",
		),
		p2p.ChainOptions(
			p2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
			p2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
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
		return errors.Wrap(err, "h peer setup failed") // TODO: change to const
	}

	p, err := NewPeer(*host.InfoFromHost(h))
	if err != nil {
		return errors.Wrap(err, "h peer setup failed") // TODO: change to const
	}

	s.host = h
	s.peer = p

	return nil
}

// ping connected peers regularly.
func (s *Server) ping(ctx context.Context) {
	for {
		for _, p := range s.peersConnected {
			if err := s.host.Connect(ctx, p.Info()); err != nil {
				s.RemovePeer(p)
				s.logger.Debug("peer dropped ", p)

				break
			}
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
				s.peersConnected[p.info.ID] = p
				s.logger.Debug("peer added ", p)
			}
		case <-time.After(networkStatePeriod):
			{
				s.logger.Infow("online", "connected", len(s.peersConnected))
			}
		}
	}
}

// AddPeer adds a peer to the network.
func (s *Server) AddPeer(ctx context.Context, peer *Peer) {
	for {
		var err error

		if s.PeerConnected(peer) {
			return
		}

		if err = s.dht.Host().Connect(ctx, peer.info); err != nil {
			s.logger.Warnw("couldn't connect to peer", "err", err)

			continue
		}

		var p *Peer

		if p, err = NewPeer(peer.info); err != nil {
			s.logger.Error("failed to initialize peer: ", err)

			continue
		}

		s.peerChan.connected <- p

		break
	}
}

// RemovePeer removes a peer from the network.
func (s *Server) RemovePeer(p *Peer) {
	delete(s.peersConnected, p.info.ID)
}

// PeerConnected checks if the peer is connected to the network.
func (s *Server) PeerConnected(p *Peer) bool {
	_, isConnected := s.peersConnected[p.info.ID]

	return isConnected
}
