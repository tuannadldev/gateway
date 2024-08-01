// auto generate code, dont edit it
package service

import(
    "gateway/config"
    "gateway/application/routing/delivery/remote"
    "gateway/application/routing/delivery/remote/grpc"
)

func initRemoteServiceClientRegistry(config *config.Config) map[string]RemoteServiceClient {
    return map[string]RemoteServiceClient {
        "CustomerService": remote.NewGrpcServiceClient(grpc.NewCustomerServiceClient(config)),
        "OrderService": remote.NewGrpcServiceClient(grpc.NewOrderServiceClient(config)),
    }
}
