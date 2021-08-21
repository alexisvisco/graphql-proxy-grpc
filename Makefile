BIN_PROTOC_PLUGIN := protoc-gen-graphql-proxy-grpc

build:
	go build -o $(BIN_PROTOC_PLUGIN) cmd/graphql-proxy-grpc-protoc/main.go \
	&& mv $(BIN_PROTOC_PLUGIN) ~/go/bin/$(BIN_PROTOC_PLUGIN)

