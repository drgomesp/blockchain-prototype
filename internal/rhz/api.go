package rhz

// Message used for direct p2p communication.
type Message interface {
	// Decode ...
	Decode(v interface{}) error
}

// PeerExchange is a bi-directional protocol for peer message exchange.
type PeerExchange interface {
	// HandleRequest something from the peer.
	// If there's a response, it will come in the form of a message.
	HandleRequest(*Peer, MessagePacket) (MessagePacket, error)

	// HandleResponse handles a message as a reply from the peer.
	HandleResponse(*Peer, MessagePacket) error
}
