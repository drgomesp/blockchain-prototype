package rpc

import (
	"fmt"
	"net"
)

type Listener net.Listener

// NewListener instantiates a new new.NetListener with
// the provided port to be exposed for TCP connections.
func NewListener(port int) Listener {
	tcp, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	return Listener(tcp)
}
