package generator

import (
	"fmt"
	"gateway/autogen/model"
	"gateway/autogen/util"
	"os"
	"runtime"
	"strings"
)

const GRPC_CLIENT_PARENT_PATH = "application/routing/delivery/remote/grpc"
const SERVICE_CLIENT_REGISTRY_PATH = "application/routing/delivery/service/remoteServiceClientRegistry.go"

func GenerateGrpcClients(protos []*model.Proto) error {
	for _, proto := range protos {
		generateGrpcClients(proto)
	}

	// TODO: capture all errors
	return nil
}

func generateGrpcClients(proto *model.Proto) error {
	for _, service := range proto.Services {
		generateGrpcClient(proto, service)
	}

	// TODO: capture all errors
	return nil
}

func generateGrpcClient(proto *model.Proto, service *model.Service) error {
	fmt.Println("Generate grpc client for:", service.Name)
	// segments := strings.Split(proto.FilePath, "/")
	var segments []string
	if runtime.GOOS == "windows" {
		segments = strings.Split(proto.FilePath, "\\")
	} else {
		segments = strings.Split(proto.FilePath, "/")
	}
	if len(segments) < 3 {
		return fmt.Errorf("unsupported proto file path: %s", proto.FilePath)
	}

	template, err := util.ReadFileAsString("autogen/template/grpc_client.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	protoFolder := strings.Join(segments[1:len(segments)-1], "/")
	template = strings.ReplaceAll(template, "<proto_folder>", protoFolder)
	template = strings.ReplaceAll(template, "<proto_package>", proto.Package)
	template = strings.ReplaceAll(template, "<service_name>", service.Name)
	template = strings.ReplaceAll(template, "<service_name_lower_1st>", util.LowerFirstChar(service.Name))

	methodConfigTemplate, err := util.ReadFileAsString("autogen/template/grpc_method_entry.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	methodEntries := []string{}
	for _, api := range service.Apis {
		entry := methodConfigTemplate
		entry = strings.ReplaceAll(entry, "<api_name>", api.Name)
		entry = strings.ReplaceAll(entry, "<api_name_lower_1st>", util.LowerFirstChar(api.Name))
		methodEntries = append(methodEntries, entry)
	}
	template = strings.ReplaceAll(template, "<method_entries>", strings.Join(methodEntries, "\n"))

	apiInvocationTemplate, err := util.ReadFileAsString("autogen/template/grpc_api_invocation.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	apiInvocations := []string{}
	for _, api := range service.Apis {
		apiInvocation := apiInvocationTemplate
		apiInvocation = strings.ReplaceAll(apiInvocation, "<proto_package>", proto.Package)
		apiInvocation = strings.ReplaceAll(apiInvocation, "<service_name_lower_1st>", util.LowerFirstChar(service.Name))
		apiInvocation = strings.ReplaceAll(apiInvocation, "<api_name_lower_1st>", util.LowerFirstChar(api.Name))
		apiInvocation = strings.ReplaceAll(apiInvocation, "<api_name>", api.Name)
		apiInvocation = strings.ReplaceAll(apiInvocation, "<request_type>", api.RequestType)
		apiInvocations = append(apiInvocations, apiInvocation)
		apiInvocations = append(apiInvocations, "\n") // method separator line
	}
	template = strings.ReplaceAll(template, "<api_invocations>", strings.Join(apiInvocations, "\n"))

	err = os.WriteFile(GRPC_CLIENT_PARENT_PATH+"/"+util.LowerFirstChar(service.Name)+"Client.go", []byte(template), 0644)
	if err != nil {
		fmt.Println("failed to write grpc client file:", err)
		return err
	}

	return nil
}

func GenerateServiceClientRegistry(protos []*model.Proto) error {
	fmt.Println("Generating service client registry")
	template, err := util.ReadFileAsString("autogen/template/service_client_registry.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}

	clientEntryTemplate, err := util.ReadFileAsString("autogen/template/grpc_service_client_entry.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	clientEntries := []string{}
	for _, proto := range protos {
		for _, service := range proto.Services {
			clientEntry := clientEntryTemplate
			clientEntry = strings.ReplaceAll(clientEntry, "<service_name>", service.Name)
			clientEntries = append(clientEntries, clientEntry)
		}
	}
	template = strings.ReplaceAll(template, "<grpc_service_client_entries>", strings.Join(clientEntries, "\n"))

	err = os.WriteFile(SERVICE_CLIENT_REGISTRY_PATH, []byte(template), 0644)
	if err != nil {
		fmt.Println("failed to write service registry file:", err)
		return err
	}

	return nil
}
