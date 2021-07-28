package p2p

import "time"

type Config struct {
	MaxPeers    uint
	PingTimeout time.Duration
}
