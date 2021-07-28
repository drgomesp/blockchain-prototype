package p2p

import "time"

type Config struct {
	MaxPeers       uint          // MaxPeers allowed to be connected to the network.
	PingTimeout    time.Duration // PingTimeout duration used to regularly ping connected peers.
	BootstrapAddrs []string
}
