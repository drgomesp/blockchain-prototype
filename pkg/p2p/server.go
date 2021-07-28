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

	running  bool         // running controls the run loop
	quit     chan bool    // quit channel to receive the stop signal
	peerChan peerChannels // peerChan manages channel-sent peers
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
	}

	if err := srv.setupLocalHost(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to initialize host")
	}

	if err := srv.setupDiscovery(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to initialize discovery")
	}

	return srv, nil
}

func (s *Server) Name() string {
	return ServerName
}

func (s *Server) Start(ctx context.Context) error {
	s.running = true

	go s.listen(ctx)
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

// addPeer adds a peer to the network.
func (s *Server) addPeer(peerInfo peer.AddrInfo) {
	s.peerChan.connected <- peerInfo
}

func (s *Server) listen(ctx context.Context) {
listening:
	for {
		select {
		case <-ctx.Done():
			s.quit <- true

			break listening
		case peerInfo := <-s.peerChan.discovered:
			s.logger.Infow("peer discovered", "peer", peerInfo.ID.Pretty())

			if err := s.dht.Host().Connect(ctx, peerInfo); err != nil {
				s.logger.With(ctx.Err()).Error("couldn't connect to peer")

				continue
			}

			go func() {
				s.logger.Infow("trying to connect to peer", "peer", peerInfo.ID.Pretty())
				s.addPeer(peerInfo)
			}()
		}
	}
}

func (s *Server) run(_ context.Context) {
running:
	for s.running {
		select {
		case <-s.quit:
			break running
		case peerInfo := <-s.peerChan.connected:
			s.logger.Infow("peer connected", "peer", peerInfo.ID.Pretty())
		}
	}
}
