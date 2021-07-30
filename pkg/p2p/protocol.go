package p2p

import "io"

type Protocol struct {
	Name string
	Run  func(p **Peer, writer io.ReadWriter) error
}
