package remote

import (
	"errors"
)

type defaultGrpcServiceClient struct {
	grpcClient GrpcServiceClient
}

func NewGrpcServiceClient(grpcClient GrpcServiceClient) *defaultGrpcServiceClient {
	return &defaultGrpcServiceClient{
		grpcClient: grpcClient,
	}
}

func (defaultClient defaultGrpcServiceClient) Invoke(method string, data interface{}, md map[string]string) (interface{}, error) {
	invoke, found := defaultClient.grpcClient.GetMethodRegistry()[method]
	if !found {
		return nil, errors.New("Method not found: " + method)
	}

	return invoke(data, md)
}
