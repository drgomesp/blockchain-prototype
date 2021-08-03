package main

import (
	"fmt"
	"log"

	"github.com/drgomesp/rhizom/pkg/rpc"
)

func main() {
	fmt.Println("running server...")

	serv := rpc.NewServer("server", rpc.NewStreamService())

	if err := serv.Start(rpc.NewListener(7000)); err != nil {
		log.Fatalf("server failed: %s", err)
	}

	fmt.Println("down")
}
