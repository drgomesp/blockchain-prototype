
all: clean
	@protoc \
		--go_out=./proto \
		--go_opt=paths=import \
    --go-grpc_out=./proto \
		--go-grpc_opt=paths=import \
		./proto/*.proto
	mv ./proto/github.com/drgomesp/rhizom/proto/gen ./proto
	rm -rf ./proto/github.com

clean:
	rm -rf ./proto/gen
