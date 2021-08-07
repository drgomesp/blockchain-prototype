package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p-core/network"
)

type ProtocolType string

var NilProtocol = ProtocolType("")

type StreamHandlerFunc func(context.Context, network.Stream) (ProtocolType, interface{}, error)

// Protocol defines a sub-protocol for communication in the network.
type Protocol struct {
	// ID is the unique identifier of the protocol (three-letter word).
	ID string

	// Run ...
	Run StreamHandlerFunc
}

type protoRW struct{}
