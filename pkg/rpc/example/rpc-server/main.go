package main

import (
	"fmt"
	"log"
	"net"

	"github.com/drgomesp/rhizom/pkg/rpc"
	"github.com/drgomesp/rhizom/proto/gen/stream"
	"google.golang.org/grpc"
)

const (
	_port       = 7000
	_nameServer = "server"
)

func main() {
	fmt.Println("running server...")

	setup := grpc.NewServer()
	stream.RegisterBlockServer(setup, rpc.NewBlockStream())
	server := rpc.NewServer(_nameServer, setup)

	net, err := net.Listen("tcp", fmt.Sprintf(":%d", _port))
	if err != nil {
		panic(err)
	}

	if err := server.Start(net); err != nil {
		log.Fatalf("server failed: %s", err)
	}

	fmt.Println("down")
}
