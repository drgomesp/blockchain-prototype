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

type PeerChannels struct {
	discovered chan peer.AddrInfo
	connected  chan peer.AddrInfo
}

// Server manages p2p connections.
type Server struct {
	cfg    Config
	logger *zap.SugaredLogger
	host   host.Host
	dht    *kaddht.IpfsDHT

	running bool

	// listen control channels
	quit     chan bool
	peerChan PeerChannels
}

func NewServer(ctx context.Context, logger *zap.SugaredLogger, config Config) (*Server, error) {
	srv := &Server{
		cfg:     config,
		logger:  logger,
		host:    nil,
		dht:     new(kaddht.IpfsDHT),
		running: false,
		quit:    make(chan bool),
		peerChan: PeerChannels{
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

func (s *Server) setupDiscovery(ctx context.Context) error {
	const serviceTag = "rhizom"

	disc, err := discovery.NewMdnsService(ctx, s.host, time.Second, serviceTag)
	if err != nil {
		return errors.Wrap(err, "failed to initialize disc")
	}

	disc.RegisterNotifee(s)

	return nil
}

func (s *Server) setupLocalHost(ctx context.Context) error {
	p2pHost, err := p2p.New(
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

	s.host = p2pHost

	return nil
}

func (s *Server) setupPeerConnection(peerInfo peer.AddrInfo) {
	s.peerChan.connected <- peerInfo
}

func (s *Server) Name() string {
	return "p2p"
}

func (s *Server) Start(ctx context.Context) error {
	s.running = true

	go s.listen(ctx)
	go s.run(ctx)

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	return nil
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
				s.setupPeerConnection(peerInfo)
			}()
		}
	}
}

func (s *Server) run(_ context.Context) {
running:
	for {
		select {
		case <-s.quit:
			break running
		case peerInfo := <-s.peerChan.connected:
			s.logger.Infow("peer connected", "peer", peerInfo.ID.Pretty())
		}
	}
}

// HandlePeerFound Receive a peer info in an channel.
func (s *Server) HandlePeerFound(peerInfo peer.AddrInfo) {
	s.peerChan.discovered <- peerInfo
}
