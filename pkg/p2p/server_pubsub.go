package p2p

import (
	"bytes"
	"context"
	"sync"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/pkg/errors"
)

var mutex sync.Mutex

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
	var wg sync.WaitGroup
	wg.Add(len(s.cfg.Topics))

	select {
	case <-ctx.Done():
		return
	default:
		{
			for _, topicName := range s.cfg.Topics {
				go func(ctx context.Context, topicName string, wg *sync.WaitGroup) {
					sub, _, err := s.subscribe(ctx, topicName)
					if err != nil {
						s.logger.Error("failed to setup subscriptions: ", err)

						return
					}

					go s.handleSubscription(ctx, sub)

					s.logger.Debug("subscribed to topic: ", sub.Topic())
					wg.Done()
				}(ctx, topicName, &wg)
			}
		}
	}

	wg.Wait()
}

// subscribe to a topic.
func (s *Server) subscribe(_ context.Context, topicName string) (*pubsub.Subscription, *pubsub.Topic, error) {
	if t, ok := s.topics[topicName]; !ok {
		topic, err := s.pubSub.Join(topicName)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to join topic %s", topicName)
		}

		sub, err := topic.Subscribe()
		if err != nil {
			return nil, topic, errors.Wrap(err, "failed to subscribe to topic: ")
		}

		mutex.Lock()
		s.topics[topicName] = topic
		mutex.Unlock()

		return sub, topic, nil
	} else {
		sub, err := t.Subscribe()
		if err != nil {
			return nil, nil, err
		}
		return sub, t, nil
	}
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

				m := Message{
					Type:    MsgType(*msg.Topic),
					Payload: bytes.NewReader(msg.Data),
				}
				var msgNewBlock struct {
					Header struct {
						Index uint64
					}
				}
				if err := m.Decode(&msgNewBlock); err != nil {
					s.logger.Error(err)

					continue
				}

				s.logger.Debugw("message received from topic", "topic", msg.Topic, "msg", msgNewBlock)
			}
		}
	}
}
