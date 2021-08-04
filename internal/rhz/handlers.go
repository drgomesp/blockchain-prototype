package rhz

func HandleGetBlocksRequest(backend PeerExchange, msg Message, peer *Peer) (MessagePacket, error) {
	req := new(GetBlocksRequest)
	if err := msg.Decode(&req); err != nil {
		return nil, err
	}

	return backend.HandleRequest(peer, req)
}

func HandleGetBlocksResponse(backend PeerExchange, msg Message, peer *Peer) error {
	res := new(GetBlocksResponse)
	if err := msg.Decode(&res); err != nil {
		return err
	}

	return backend.HandleResponse(peer, res)
}
