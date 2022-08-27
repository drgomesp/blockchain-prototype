PB = $(wildcard *.proto)
GO = $(PB:.proto=.pb.go)

all: $(GO)

%.pb.go: %.proto
		protoc --gofast_out=. message.proto

clean:
		rm -f *.pb.go
		rm -f *.go
