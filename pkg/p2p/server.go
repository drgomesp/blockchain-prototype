package p2p

import (
	"context"
	"time"

	p2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
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

	logger  *zap.SugaredLogger
	host    host.Host
	dht     *kaddht.IpfsDHT
	notifee *notifee
}

func NewServer(ctx context.Context, logger *zap.SugaredLogger, config Config) (*Server, error) {
	var dht *kaddht.IpfsDHT

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
			dht, err = kaddht.New(ctx, h)
			if err != nil {
				return nil, errors.Wrap(err, "failed to initialize dht")
			}

			return dht, nil
		}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize host")
	}

	const serviceTag = "rhizom"

	disc, err := discovery.NewMdnsService(ctx, p2pHost, time.Second, serviceTag)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize disc")
	}

	n := newNotifee()
	disc.RegisterNotifee(n)

	srv := &Server{
		Config:  config,
		logger:  logger,
		host:    p2pHost,
		dht:     dht,
		notifee: n,
	}

	return srv, nil
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				s.logger.With(ctx.Err()).Error("context done")
			case p := <-s.notifee.PeerChan:
				if err := s.dht.Host().Connect(ctx, p); err != nil {
					s.logger.With(ctx.Err()).Error("couldn't connect to peer")
				}

				s.logger.With("peer", p).Info("connected to peer")
			}
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return nil
}
