package rhz

// HandleGetBlocks handles an incoming message from a peer that is requesting us something.
func HandleGetBlocks(backend PeerExchange, msg Message, peer *Peer) (MessagePacket, error) {
	req := new(MsgGetBlocks)
	if err := msg.Decode(&req); err != nil {
		return nil, err
	}

	return backend.ReceiveRequest(peer, req)
}

// HandleBlocks handles an incoming message from a peer that responded to our request.
func HandleBlocks(backend PeerExchange, msg Message, peer *Peer) error {
	res := new(MsgBlocks)
	if err := msg.Decode(&res); err != nil {
		return err
	}

	return backend.ReceiveResponse(peer, res)
}
