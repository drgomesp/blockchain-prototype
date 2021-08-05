package p2p

import "context"

// Protocol defines a sub-protocol for communication in the network.
type Protocol struct {
	// ID is the unique identifier of the protocol (three-letter word).
	ID string

	// Run ...
	Run func(context.Context, MsgReadWriter) error
}
