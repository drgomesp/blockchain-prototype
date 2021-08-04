package rhz

import "github.com/drgomesp/rhizom/pkg/p2p"

type MessagePacket interface {
	Type() p2p.MsgType
}

type GetBlocksRequest struct {
	IndexHave uint64
	IndexNeed uint64
}

type GetBlocksResponse struct {
	IsUpdated bool
	Chain     []struct {
		Header struct {
			Index uint64
		}
	}
}

func (g *GetBlocksRequest) Type() p2p.MsgType {
	return MsgGetBlocksRequest
}

func (g *GetBlocksResponse) Type() p2p.MsgType {
	return MsgGetBlocksResponse
}
