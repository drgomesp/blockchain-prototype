package p2p

import (
	"context"
	"time"

	p2p "github.com/libp2p/go-libp2p"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	router "github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"
)

const networkStatePeriod = 5 * time.Second

// peerChannels where peers are sent through.
type peerChannels struct {
	discovered chan *Peer // discovered peers found through Kademlia DHT
	connected  chan *Peer // connected peers in the network
}

const ServerName = "p2p"

// Server manages a p2p network.
type Server struct {
	// Dependencies
	cfg       Config                   // cfg server options.
	logger    *zap.SugaredLogger       // logger provided logger.
	protocols []Protocol               // protocols supported by the server.
	host      host.Host                // host is the actual p2p node within the network.
	dht       *kaddht.IpfsDHT          // dht discovery service.
	pubSub    *pubsub.PubSub           // pubSub is a p2p publish/subscribe service.
	topics    map[string]*pubsub.Topic // topics of that the server is subscribed to.
	peer      *Peer                    // Peer is the local p2p peer.

	// Control flags and channels
	running         bool             // running controls the run loop.
	quit            chan bool        // quit channel to receive the stop signal.
	peerChan        peerChannels     // peerChan manages channel-sent peers.
	peersDiscovered map[string]*Peer // peersDiscovered holds the discovered Peer nodes.
	peersConnected  map[string]*Peer // peersConnected holds recently connected Peer nodes
	disc            mdns.Service
}

type ServerOption func(*Server)

func WithLogger(l *zap.SugaredLogger) ServerOption {
	return func(srv *Server) {
		srv.logger = l
	}
}

// NewServer initializes a p2p Server from a given Config capable of managing a network.
func NewServer(config Config, opt ...ServerOption) (*Server, error) {
	srv := &Server{
		cfg:     config,
		dht:     new(kaddht.IpfsDHT),
		topics:  make(map[string]*pubsub.Topic),
		running: false,
		quit:    make(chan bool),
		peerChan: peerChannels{
			discovered: make(chan *Peer),
			connected:  make(chan *Peer),
		},
		peersDiscovered: make(map[string]*Peer),
		peersConnected:  make(map[string]*Peer),
	}

	for _, option := range opt {
		option(srv)
	}

	return srv, nil
}

func (s *Server) Name() string {
	return ServerName
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.setupLocalHost(ctx); err != nil {
		return errors.Wrap(err, "failed to initialize peer")
	}

	if err := s.setupDiscovery(); err != nil {
		return errors.Wrap(err, "failed to initialize discovery")
	}

	if err := s.setupPubSub(ctx); err != nil {
		return errors.Wrap(err, "failed to initialize pubsub")
	}

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
		p2p.Routing(func(h host.Host) (router.PeerRouting, error) {
			var err error
			s.dht, err = kaddht.New(ctx, h)
			if err != nil {
				return nil, errors.Wrap(err, "failed to initialize dht")
			}

			return s.dht, nil
		}),
	)

	for _, addr := range h.Addrs() {
		log.Debug().Msgf(addr.String())
	}

	if err != nil {
		// TODO: change to const
		return errors.Wrap(err, "peer setup failed")
	}

	p, err := NewPeer(host.InfoFromHost(h), s.pubSub)
	if err != nil {
		// TODO: change to const
		return errors.Wrap(err, "peer setup failed")
	}

	s.host = h
	s.peer = p

	return nil
}

// ping connected peers regularly.
func (s *Server) ping(ctx context.Context) {
	for {
		for _, p := range s.peersConnected {
			if err := s.host.Connect(ctx, *p.info); err != nil {
				s.RemovePeer(p)
				log.Debug().Msgf("peer dropped %v", p)

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
				s.peersConnected[p.info.String()] = p
			}
		case <-time.After(networkStatePeriod):
			{
				// s.logger.Infow("online", "connected", len(s.peersConnected))
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

		if err = s.dht.Host().Connect(ctx, *peer.info); err != nil {
			//log.Warn().Msgf("couldn't connect to peer", "err", err)

			continue
		}

		var p *Peer

		if p, err = NewPeer(peer.info, s.pubSub); err != nil {
			s.logger.Error("failed to initialize peer: ", err)

			continue
		}

		s.peerChan.connected <- p

		break
	}
}

// RemovePeer removes a peer from the network.
func (s *Server) RemovePeer(p *Peer) {
	delete(s.peersConnected, p.String())
}

// PeerConnected checks if the peer is connected to the network.
func (s *Server) PeerConnected(p *Peer) bool {
	_, is := s.peersConnected[p.String()]

	return is
}

// RegisterProtocols registers the server protocol set.
func (s *Server) RegisterProtocols(protocols ...Protocol) {
	s.protocols = protocols
	s.registerProtocols(context.Background())
}

func (s *Server) connectPeerByAddr(ctx context.Context, addr string) (*Peer, error) {
	peerAddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize multiaddr")
	}

	peerInfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load addr info from multiaddr")
	}

	return s.setupConnection(ctx, peerInfo), nil
}

// setupProtocolConnection sets up a peer connection from an incoming network.Stream.
func (s *Server) setupProtocolConnection(ctx context.Context, peerInfo *peer.AddrInfo) *Peer {
	p := &Peer{
		info:   peerInfo,
		pubSub: s.pubSub,
		conn:   &connection{},
	}

	go s.AddPeer(ctx, p)

	return p
}

// setupConnection sets up a peer connection, runs the handshakes
// and tries to add the connection as a Peer.
func (s *Server) setupConnection(ctx context.Context, peerInfo *peer.AddrInfo) *Peer {
	p := &Peer{
		info:   peerInfo,
		pubSub: s.pubSub,
	}

	go s.AddPeer(ctx, p)

	return p
}