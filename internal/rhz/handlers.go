package rhz

// HandleGetBlocks handles an incoming message from a peer that is requesting us something.
func HandleGetBlocks(backend API, msg Message, peer *Peer) (Message, error) {
	req := new(MsgGetBlocks)
	if err := msg.Decode(&req); err != nil {
		return nil, err
	}

	//res, err := backend.GetBlocks(peer, req.(Message))
	//if err != nil {
	//	return nil, err
	//}
	//
	//return res.(*MessagePacket), nil
	return nil, nil
}

// HandleBlocks handles an incoming message from a peer that responded to our request.
func HandleBlocks(backend API, msg Message, peer *Peer) error {
	res := new(MsgBlocks)
	if err := msg.Decode(&res); err != nil {
		return err
	}

	return backend.Blocks(peer, res)
}
