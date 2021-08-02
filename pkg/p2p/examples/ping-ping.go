package main

import (
	"context"

	"github.com/drgomesp/rhizom/pkg/p2p"
	"go.uber.org/zap"
)

func main() {
}

// pingPongService runs a simple ping-pong protocol between nodes.
type pingPongService struct {
	logger *zap.SugaredLogger
}

func NewPingPongService(logger *zap.SugaredLogger) *pingPongService {
	return &pingPongService{logger}
}

func (p *pingPongService) Name() string {
	return "ping-pong"
}

func (p *pingPongService) Start(ctx context.Context) error {
	p.logger.Infow("ping-pong service starting...")

	return nil
}

func (p *pingPongService) Stop(ctx context.Context) error {
	p.logger.Infow("ping-pong service stopping...")

	return nil
}

func (p *pingPongService) Protocols() []p2p.Protocol {
	return []p2p.Protocol{
		{
			Name: "ping-pong",
			Run:  p.Run,
		},
	}
}

const (
	MsgPing = uint64(iota)
	MsgPong
)

func (p *pingPongService) Run(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
	//errChan := make(chan error)
	//
	//go func() {
	//	for range time.Tick(3 * time.Second) {
	//		p.logger.Info("sending ping...")
	//
	//		if err := peer.SendMsg(rw, MsgPing, "data"); err != nil {
	//			errChan <- err
	//		}
	//	}
	//}()

	return nil
}
