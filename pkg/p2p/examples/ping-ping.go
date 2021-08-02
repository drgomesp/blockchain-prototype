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

func (p *pingPongService) Start(_ context.Context) error {
	p.logger.Infow("ping-pong service starting...")

	return nil
}

func (p *pingPongService) Stop(_ context.Context) error {
	p.logger.Infow("ping-pong service stopping...")

	return nil
}

func (p *pingPongService) Protocols() []p2p.Protocol {
	return []p2p.Protocol{
		{
			Name: "ping-pong",
		},
	}
}

const (
	MsgPing = uint64(iota)
	MsgPong
)
