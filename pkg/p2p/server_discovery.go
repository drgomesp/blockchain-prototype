package p2p

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p/p2p/discovery"
	"github.com/pkg/errors"
)

// setupDiscovery sets up the peer discovery mechanism.
func (s *Server) setupDiscovery(ctx context.Context) error {
	const serviceTag = "rhizom"

	disc, err := discovery.NewMdnsService(ctx, s.node, time.Second, serviceTag)
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
		case p := <-s.peerChan.discovered:
			{
				s.logger.Info("peer discovered ", p)
				go s.AddPeer(ctx, p)
			}
		}
	}
}
