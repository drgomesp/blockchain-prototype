package p2p

type StreamingFunc func(data []byte)

// Protocol defines a sub-protocol for communication in the network.
type Protocol struct {
	// Name of the protocol (three-letter word).
	Name string

	// Run ...
	// Run func(p *Peer, msg *Message) error
}
