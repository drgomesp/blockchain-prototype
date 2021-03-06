package p2p

import "time"

// Config of the server.
type Config struct {
	NetworkName    string
	ServiceTag     string        // ServiceTag is an identifier for services managed b the server.
	MaxPeers       uint          // MaxPeers allowed to be connected to the network.
	PingTimeout    time.Duration // PingTimeout duration used to regularly ping connected peers.
	BootstrapAddrs []string      // BootstrapAddrs holds a list of bootstrap peer addresses.
	Topics         []string      // Topics which the peer wants to subscribe to.
}
