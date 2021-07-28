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

// Server manages p2p connections.
type Server struct {
	Config

	logger *zap.SugaredLogger
	host   host.Host
	dht    *kaddht.IpfsDHT

	running bool

	// run control channels
	quit      chan struct{}
	peerAdded chan peer.AddrInfo
}

func NewServer(ctx context.Context, logger *zap.SugaredLogger, config Config) (*Server, error) {
	srv := &Server{
		Config:    config,
		logger:    logger,
		host:      nil,
		dht:       new(kaddht.IpfsDHT),
		running:   false,
		quit:      make(chan struct{}),
		peerAdded: make(chan peer.AddrInfo),
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

func (s *Server) Name() string {
	return "p2p"
}

func (s *Server) Start(ctx context.Context) error {
	s.running = true
	go s.run(ctx)

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	return nil
}

func (s *Server) run(ctx context.Context) {
running:
	for {
		select {
		case <-ctx.Done():
			s.logger.With(ctx.Err()).Error("context done")

			break running

		case <-s.quit:
			break running

		case p := <-s.peerAdded:
			if err := s.dht.Host().Connect(ctx, p); err != nil {
				s.logger.With(ctx.Err()).Error("couldn't connect to peer")
			}

			s.logger.With("peer", p).Info("connected to peer")
		}
	}
}

// HandlePeerFound Receive a peer info in an channel.
func (s *Server) HandlePeerFound(pi peer.AddrInfo) {
	s.peerAdded <- pi
}
