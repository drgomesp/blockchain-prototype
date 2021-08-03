package rhz

import (
	"github.com/drgomesp/rhizom/pkg/p2p"
)

// request/response messages.
const (
	// MsgGetBlocksRequest represents a request for blocks.
	MsgGetBlocksRequest = p2p.MsgType("/rhz/blocks/req/default_2b678c95-27d5-4f09-bf38-a62be2c5339b")
	// MsgGetBlocksResponse represents a response to a request for blocks.
	MsgGetBlocksResponse = p2p.MsgType("/rhz/blocks/resp/default_2b678c95-27d5-4f09-bf38-a62be2c5339b")
)

// async messages.
const (
	// MsgNewBlock represents a new block.
	MsgNewBlock = p2p.MsgType("NewBlock")
)
