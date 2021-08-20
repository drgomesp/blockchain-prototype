package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/drgomesp/rhizom/proto/gen/entity"
	"github.com/drgomesp/rhizom/proto/gen/message"
	"github.com/drgomesp/rhizom/proto/gen/stream"
	"google.golang.org/grpc"
)

const serverAddr = "localhost:7000"

var blockchain []*entity.Block

func streamData(stream stream.Block_GetBlockClient) {
	req := &message.RequestStreamGetBlock{
		IndexWant: 1,
	}
	for {
		switch resp, err := stream.Recv(); err {
		case nil:
			if resp == nil || resp.Err != "" {
				log.Fatal(err)
			}
			blockchain = append(blockchain, resp.Block)
			req = &message.RequestStreamGetBlock{
				IndexWant: resp.Block.Index + 1,
			}
			if err := stream.Send(req); err != nil {
				log.Fatalf("failed to send req: %s", err)
			}
			fmt.Println(resp)

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

	streamData(stream)
}
