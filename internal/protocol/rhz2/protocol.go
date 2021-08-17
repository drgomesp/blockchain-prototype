package rhz2

import (
	"context"
	"io/ioutil"
	"log"

	pb "github.com/drgomesp/rhizom/internal/protocol/rhz2/pb"
	"github.com/drgomesp/rhizom/pkg/p2p"
	oldproto "github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func ProtocolHandlerFunc() p2p.StreamHandlerFunc {
	return func(ctx context.Context, rw p2p.MsgReadWriter) (
		p2p.ProtocolType, interface{}, error,
	) {
		msg, err := rw.ReadMsg(ctx)
		if err != nil {
			return p2p.NilProtocol, nil, errors.Wrap(err, "message read failed")
		}

		data, err := ioutil.ReadAll(msg.Payload)
		if err != nil {
			return p2p.NilProtocol, nil, errors.Wrap(err, "message read failed")
		}

		req := &pb.GetBlocks_Request{}
		if err := proto.Unmarshal(data, oldproto.MessageV2(req)); err != nil {
			log.Fatalln("Failed to parse address book:", err)
		}

		s, err := protojson.Marshal(oldproto.MessageV2(req))
		if err != nil {
			return p2p.NilProtocol, nil, errors.Wrap(err, "message read failed")
		}

		log.Println(string(s))

		return p2p.NilProtocol, nil, nil
	}
}
