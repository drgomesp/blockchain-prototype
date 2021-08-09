package rhz

const (
	rhz     = "/rhz/"
	local   = "drgomesp"
	devNet  = "default_55fd187b-b29e-4856-81fc-ba1e7bc18287"
	testNet = "rhz_testnet_e19d2c16-8c39-4f0f-8c88-2427e37c12bb"
	net     = devNet

	TopicBlocks       = rhz + "blk/" + net
	TopicProducers    = rhz + "prc/" + net
	TopicTransactions = rhz + "tx/" + net
	TopicRequestSync  = rhz + "blkchain/req/" + net

	ProtocolRequestBlocks     = rhz + "blocks/req/" + net
	ProtocolResponseBlocks    = rhz + "blocks/resp/" + net
	ProtocolRequestDelegates  = rhz + "delegates/req/" + net
	ProtocolResponseDelegates = rhz + "delegates/resp/" + net
)
