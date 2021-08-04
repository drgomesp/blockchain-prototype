package main

import (
	"fmt"
	"log"

	"github.com/drgomesp/rhizom/pkg/rpc"
	"github.com/drgomesp/rhizom/pkg/service"
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
	stream.RegisterBlockServer(setup, service.NewBlockStream())
	server := rpc.NewServer(_nameServer, setup)

	if err := server.Start(rpc.NewListener(_port)); err != nil {
		log.Fatalf("server failed: %s", err)
	}

	fmt.Println("down")
}
