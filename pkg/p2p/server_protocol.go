package p2p

import (
	"context"
)

type Streaming func(data []byte) ([]byte, Protocol, error)

func (s *Server) registerProtocols(ctx context.Context) {
}
