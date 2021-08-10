package p2p

import (
	"context"
	"time"

	p2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	router "github.com/libp2p/go-libp2p-core/routing"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	mplex "github.com/libp2p/go-libp2p-mplex"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	secio "github.com/libp2p/go-libp2p-secio"
	yamux "github.com/libp2p/go-libp2p-yamux"
	"github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
	"github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
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
	running         bool              // running controls the run loop.
	quit            chan bool         // quit channel to receive the stop signal.
	peerChan        peerChannels      // peerChan manages channel-sent peers.
	peersDiscovered map[peer.ID]*Peer // peersDiscovered holds the discovered Peer nodes.
	peersConnected  map[peer.ID]*Peer // peersConnected holds recently connected Peer nodes
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
		topics:  make(map[string]*pubsub.Topic, 0),
		running: false,
		quit:    make(chan bool),
		peerChan: peerChannels{
			discovered: make(chan *Peer),
			connected:  make(chan *Peer),
		},
		peersDiscovered: make(map[peer.ID]*Peer),
		peersConnected:  make(map[peer.ID]*Peer),
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

	if err := s.setupDiscovery(ctx); err != nil {
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
		ctx,
		p2p.ChainOptions(
			p2p.Transport(tcp.NewTCPTransport),
			p2p.Transport(ws.New),
		),
		p2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/0",
			"/ip4/0.0.0.0/tcp/0/ws", // WebSocker address
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

	p, err := NewPeer(host.InfoFromHost(h), s.pubSub)
	if err != nil {
		return errors.Wrap(err, "peer setup failed") // TODO: change to const
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
			s.logger.Warnw("couldn't connect to peer", "err", err)

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
	delete(s.peersConnected, p.info.ID)
}

// PeerConnected checks if the peer is connected to the network.
func (s *Server) PeerConnected(p *Peer) bool {
	_, is := s.peersConnected[p.info.ID]

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

	return s.setupConnection(ctx, peerInfo)
}

// setupProtocolConnection sets up a peer connection from an incoming network.Stream.
func (s *Server) setupProtocolConnection(
	ctx context.Context,
	peerInfo *peer.AddrInfo,
	stream network.Stream,
) (*Peer, error) {
	s.logger.Debugw(
		"setting up protocol connection",
		"protocol", stream.Protocol(), "peer", peerInfo.ID.ShortString(),
	)

	p, err := s.setupConnection(ctx, peerInfo)
	if err != nil {
		return nil, err
	}

	p.conn = &connection{}

	return p, nil
}

// setupProtocolConnection sets up a peer connection, runs the handshakes and tries to
// add the connection as a Peer.
func (s *Server) setupConnection(
	ctx context.Context,
	peerInfo *peer.AddrInfo,
) (*Peer, error) {
	p := &Peer{
		info:   peerInfo,
		pubSub: s.pubSub,
	}

	go s.AddPeer(ctx, p)

	return p, nil
}
