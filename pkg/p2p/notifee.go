package p2p

import "github.com/libp2p/go-libp2p-core/peer"

// notifee Struct to notify when a new peer is found.
type notifee struct {
	PeerChan chan peer.AddrInfo
}

// newNotifee returns a new notifier of mdns.
func newNotifee() *notifee {
	return &notifee{
		PeerChan: make(chan peer.AddrInfo),
	}
}

// HandlePeerFound Receive a peer info in an channel.
func (n *notifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}
