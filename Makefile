PROTO_FILES=$(shell find proto -type f \( -iname "*.proto" ! -iname "validate.proto" \))

protos:
	protoc -I. \
		--go-grpc_out=require_unimplemented_servers=false:.. \
		--go_out=. --validate_out="lang=go:." \
		--go-grpc_out=. proto/**/*.proto	\
		$(PROTO_FILES)

refreshing:
	protoc --go-grpc_out=require_unimplemented_servers=false:. --go_out=. proto/**/*.proto
	go run autogen/main.go
	swag init --parseDependency -g cmd/main.go

server:
	go run cmd/main.go
