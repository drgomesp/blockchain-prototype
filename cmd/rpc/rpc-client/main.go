package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/drgomesp/rhizom/proto/gen/entity"
	"github.com/drgomesp/rhizom/proto/gen/message"
	"github.com/drgomesp/rhizom/proto/gen/stream"
	"google.golang.org/grpc"
)

const serverAddr = "localhost:7000"

var blockchain []*entity.Block

func streamData(stream stream.Block_GetBlockClient, ch chan<- uint32) {
	defer close(ch)
	for {
		switch resp, err := stream.Recv(); err {
		case nil:
			if resp == nil || resp.Err != "" {
				log.Fatal(err)
			}
			blockchain = append(blockchain, resp.Block)
			ch <- resp.Block.Index

		case io.EOF:
			return

		default:
			panic(err)
		}
	}
}

func main() {
	fmt.Println("running client...")

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	stream, err := stream.NewBlockClient(conn).GetBlock(context.Background())
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := stream.CloseSend(); err != nil {
			panic(err)
		}
	}()

	ch := make(chan uint32)

	go streamData(stream, ch)
	req := &message.GetBlockRequest{}

	if err := stream.Send(req); err != nil {
		log.Fatalf("failed to send req: %s", err)
	}

	tick := time.Tick(time.Second * 10)

	for {
		fmt.Printf("\nblocks: %d\n", len(blockchain))

		select {
		case i := <-ch:
			i++
			req = &message.GetBlockRequest{Want: i}

			if err := stream.Send(req); err != nil {
				log.Fatalf("failed to send req: %s", err)
			}
			time.Sleep(time.Second)

		case <-tick:
			return
		}
	}
}
