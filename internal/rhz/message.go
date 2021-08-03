package rhz

type MessagePacket struct{}

type GetBlocksRequest struct {
	MessagePacket

	IndexHave uint64
	IndexNeed uint64
}

type GetBlocksResponse struct {
	MessagePacket
	Blocks []interface{}
}
