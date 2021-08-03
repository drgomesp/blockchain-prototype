package p2p

type RequestResponseProtocol interface{}

// ProtoRunFunc ...
type ProtoRunFunc func([]byte) error

// Protocol defines a sub-protocol for communication in the network.
type Protocol struct {
	// ID is the unique identifier of the protocol (three-letter word).
	ID string

	// Run ...
	Run ProtoRunFunc
}
