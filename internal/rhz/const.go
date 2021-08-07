package rhz

const (
	RHZ                       = "/rhz/"
	local                     = "drgomesp"
	devNet                    = "default_2b678c95-27d5-4f09-bf38-a62be2c5339b"
	testNet                   = "rhz_testnet_e19d2c16-8c39-4f0f-8c88-2427e37c12bb"
	net                       = testNet
	TopicBlocks               = RHZ + "blk/" + net
	TopicProducers            = RHZ + "prc/" + net
	TopicTransactions         = RHZ + "tx/" + net
	TopicRequestSync          = RHZ + "blkchain/req/" + net
	ProtocolRequestBlocks     = RHZ + "blocks/req/" + net
	ProtocolResponseBlocks    = RHZ + "blocks/resp/" + net
	ProtocolRequestDelegates  = RHZ + "delegates/req/" + net
	ProtocolResponseDelegates = RHZ + "delegates/resp/" + net

	ProtocolBlocks = RHZ + "blocks/" + net
)
