package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/pkg/errors"
)

// setupPubSub initializes the pub/sub mechanism.
func (s *Server) setupPubSub(ctx context.Context) error {
	ps, err := pubsub.NewGossipSub(ctx, s.node)
	if err != nil {
		s.logger.Error()
		return errors.Wrap(err, "failed to initialize gossip sub")
	}

	s.pubSub = ps

	return nil
}

// setupSubscriptions sets the pub/sub topic subscriptions.
func (s *Server) setupSubscriptions(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		{
			for _, topic := range s.cfg.Topics {
				go func(topicName string) {
					s.subscribe(ctx, topicName)
				}(topic)
			}
		}
	}
}

// subscribe to a topic.
func (s *Server) subscribe(ctx context.Context, topicName string) {
	topic, err := s.pubSub.Join(topicName)
	if err != nil {
		s.logger.Error("failed to join to topic: ", err)
	}

	sub, err := topic.Subscribe()
	if err != nil {
		s.logger.Error("failed to subscribe to topic: ", err)
	}

	s.logger.Debug("subscribed to topic: ", sub.Topic())
}

// RegisterPeerSubscription sets up the topic subscriptions for a given peer.
func (s *Server) RegisterPeerSubscription(ctx context.Context, peerInfo peer.AddrInfo) {
}
