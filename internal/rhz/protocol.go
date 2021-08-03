package rhz

type GetBlocksRequest struct {
	IndexHave uint64
	IndexNeed uint64
}

type GetBlocksResponse struct {
	Blocks []interface{}
}
