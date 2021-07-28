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
	"github.com/libp2p/go-libp2p/p2p/discovery"
	"github.com/libp2p/go-tcp-transport"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// peerChannels manages channels where peers are sent through.
type peerChannels struct {
	discovered chan peer.AddrInfo // discovered peers found through Kademlia DHT
	connected  chan peer.AddrInfo // connected peers in the network
}

const ServerName = "p2p.server"

// Server manages p2p connections.
type Server struct {
	cfg    Config
	logger *zap.SugaredLogger
	host   host.Host
	dht    *kaddht.IpfsDHT

	running        bool                      // running controls the run loop
	quit           chan bool                 // quit channel to receive the stop signal
	peerChan       peerChannels              // peerChan manages channel-sent peers
	peersConnected map[peer.ID]peer.AddrInfo // peersConnected holds recently connected peers
}

// NewServer initializes a p2p Server from a given Config capable of managing a network.
func NewServer(ctx context.Context, logger *zap.SugaredLogger, config Config) (*Server, error) {
	srv := &Server{
		cfg:     config,
		logger:  logger,
		host:    nil,
		dht:     new(kaddht.IpfsDHT),
		running: false,
		quit:    make(chan bool),
		peerChan: peerChannels{
			discovered: make(chan peer.AddrInfo),
			connected:  make(chan peer.AddrInfo),
		},
		peersConnected: make(map[peer.ID]peer.AddrInfo),
	}

	if err := srv.setupLocalHost(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to initialize host")
	}

	if err := srv.setupDiscovery(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to initialize discovery")
	}

	return srv, nil
}

// HandlePeerFound receives a discovered peer.
func (s *Server) HandlePeerFound(peerInfo peer.AddrInfo) {
	s.peerChan.discovered <- peerInfo
}

func (s *Server) Name() string {
	return ServerName
}

func (s *Server) Start(ctx context.Context) error {
	s.running = true

	go s.discover(ctx)
	go s.ping(ctx)
	go s.run(ctx)

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.running = false

	return nil
}

// setupLocalHost sets up the local p2p host.
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
		return errors.Wrap(err, "local host setup failed")
	}

	s.host = h

	return nil
}

// setupDiscovery sets up the peer discovery mechanism.
func (s *Server) setupDiscovery(ctx context.Context) error {
	const serviceTag = "rhizom"

	disc, err := discovery.NewMdnsService(ctx, s.host, time.Second, serviceTag)
	if err != nil {
		return errors.Wrap(err, "failed to initialize disc")
	}

	disc.RegisterNotifee(s)

	return nil
}

// discover for incoming discovered peers.
func (s *Server) discover(ctx context.Context) {
listening:
	for {
		select {
		case <-ctx.Done():
			{
				s.quit <- true

				break listening
			}
		case peerInfo := <-s.peerChan.discovered:
			{
				s.logger.Info("peer discovered ", peerInfo.ID.ShortString())
				s.AddPeer(ctx, peerInfo)
			}
		}
	}
}

// ping connected peers regularly.
func (s *Server) ping(ctx context.Context) {
	for {
		for _, p := range s.peersConnected {
			if err := s.host.Connect(ctx, p); err != nil {
				s.RemovePeer(p)

				break
			}

			s.logger.Debug("peer check ", p.ID.ShortString())
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
		case peerInfo := <-s.peerChan.connected:
			{
				s.peersConnected[peerInfo.ID] = peerInfo
				s.logger.Info("peer added ", peerInfo.ID.ShortString())

				break
			}
		default:
			{
				s.logger.Info("connected peers: ", len(s.peersConnected))

				time.Sleep(time.Second)
			}
		}
	}
}

// AddPeer adds a peer to the network.
func (s *Server) AddPeer(ctx context.Context, peerInfo peer.AddrInfo) {
	_, isConnected := s.peersConnected[peerInfo.ID]
	if isConnected {
		return
	}

	if err := s.dht.Host().Connect(ctx, peerInfo); err != nil {
		s.logger.Warnw("couldn't connect to peer", "err", err)
	}

	s.peerChan.connected <- peerInfo
}

// RemovePeer removes a peer from the network.
func (s *Server) RemovePeer(peerInfo peer.AddrInfo) {
	delete(s.peersConnected, peerInfo.ID)
	s.logger.Info("peer removed ", peerInfo.ID.ShortString())
}
