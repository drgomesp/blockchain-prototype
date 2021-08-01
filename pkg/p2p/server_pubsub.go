package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/pkg/errors"
)

// setupPubSub initializes the pub/sub mechanism.
func (s *Server) setupPubSub(ctx context.Context) error {
	ps, err := pubsub.NewGossipSub(ctx, s.host)
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
				go func(ctx context.Context, topicName string) {
					sub, err := s.subscribe(ctx, topicName)
					if err != nil {
						s.logger.Error("failed to setup subscriptions: ", err)

						return
					}

					s.handleSubscription(ctx, sub)
				}(ctx, topic)
			}
		}
	}
}

// subscribe to a topic.
func (s *Server) subscribe(_ context.Context, topicName string) (*pubsub.Subscription, error) {
	topic, err := s.pubSub.Join(topicName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to join topic %s", topicName)
	}

	sub, err := topic.Subscribe()
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe to topic: ")
	}

	s.logger.Debug("subscribed to topic: ", sub.Topic())

	return sub, nil
}

// handleSubscription handles messages from a given subscription.
func (s *Server) handleSubscription(ctx context.Context, sub *pubsub.Subscription) {
	for {
		select {
		case <-ctx.Done():
			{
				s.logger.Error(ctx.Err())
				sub.Cancel()

				return
			}
		default:
			{
				msg, err := sub.Next(ctx)
				if err != nil {
					s.logger.Error("failed get next topic message: ", err)

					continue
				}

				var pm Message
				if err = pm.Decode(msg); err != nil {
					s.logger.Error("unmarshal block failed: ", err)

					continue
				}

				s.logger.Debugw("message received", "msg", pm)
			}
		}
	}
}

// RegisterPeerSubscription sets up the topic subscriptions for a given peer.
func (s *Server) RegisterPeerSubscription(ctx context.Context, peerInfo peer.AddrInfo) {
}
