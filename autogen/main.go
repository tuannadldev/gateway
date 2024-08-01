package main

import (
	"fmt"

	"gateway/autogen/generator"
	"gateway/autogen/model"
	"gateway/autogen/parser"
	"gateway/autogen/util"
)

const PROTO_PATH = "proto"

func main() {
	// list proto files
	// for each proto file:
	//   list services
	//   create handlers: for each method
	//     create handler
	//   create remote client
	//   register routes
	//   register handlers
	//   register remote clients

	protoPaths := util.GetAllProtoFiles(PROTO_PATH)

	protos := []*model.Proto{}
	for _, path := range protoPaths {
		fmt.Println("Found proto file at:", path)
		proto := parser.ParseProtoFile(path)
		protos = append(protos, proto)
	}

	handlers, handlerPaths, err := generator.GenerateHandlers(protos)
	if err != nil {
		fmt.Println("failed to generate handlers:", err)
	}

	err = generator.GenerateRoutingRegistry(handlers, handlerPaths)
	if err != nil {
		fmt.Println("failed to generate routing registry:", err)
	}

	err = generator.GenerateGinRoutes(protos)
	if err != nil {
		fmt.Println("failed to generate gin routes:", err)
	}

	err = generator.GenerateGrpcClients(protos)
	if err != nil {
		fmt.Println("failed to generate grpc clients:", err)
	}

	err = generator.GenerateServiceClientRegistry(protos)
	if err != nil {
		fmt.Println("failed to generate service client registry:", err)
	}
}
