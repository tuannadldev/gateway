// auto generate code, dont edit it
package grpc

import(
	"fmt"

	"google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"gateway/config"
	"go.elastic.co/apm/module/apmgrpc/v2"
    "gateway/pkg/monitor"
    "gateway/proto/customer"
)

func (client *customerServiceClient) initMethodRegistry() {
    client.methodRegistry = map[string]func(interface{}, map[string]string) (interface{}, error) {
        "GetCustomer": client.getCustomer,
        "CreateCustomer": client.createCustomer,
        "UpdateCustomer": client.updateCustomer,
    }
}

type customerServiceClient struct {
    grpcClient customer.CustomerServiceClient
    methodRegistry map[string]func(interface{}, map[string]string) (interface{}, error)
}

func NewCustomerServiceClient(config *config.Config) *customerServiceClient {
    // using WithInsecure() because no SSL running
    cc, err := grpc.Dial(
        config.Service.CustomerServiceUrl,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithUnaryInterceptor(apmgrpc.NewUnaryClientInterceptor()),
        grpc.WithStreamInterceptor(apmgrpc.NewStreamClientInterceptor()),
    )

    if err != nil {
        fmt.Println("Could not connect:", err)
    }

    client := customerServiceClient{
        grpcClient: customer.NewCustomerServiceClient(cc),
    }
    client.initMethodRegistry()

    return &client
}

func (client *customerServiceClient) GetMethodRegistry() map[string]func(interface{}, map[string]string) (interface{}, error) {
    return client.methodRegistry
}

func (client *customerServiceClient) getCustomer(data interface{}, md map[string]string) (interface{}, error) {
    ctx := monitor.GetApmContext()
    //ctx := context.Background()
	if len(md) > 0{
        ctx = metadata.NewOutgoingContext(
            ctx,
            metadata.New(md),
        )
    }
    return client.grpcClient.GetCustomer(ctx, data.(*customer.GetCustomerRequest))
}


func (client *customerServiceClient) createCustomer(data interface{}, md map[string]string) (interface{}, error) {
    ctx := monitor.GetApmContext()
    //ctx := context.Background()
	if len(md) > 0{
        ctx = metadata.NewOutgoingContext(
            ctx,
            metadata.New(md),
        )
    }
    return client.grpcClient.CreateCustomer(ctx, data.(*customer.CreateCustomerRequest))
}


func (client *customerServiceClient) updateCustomer(data interface{}, md map[string]string) (interface{}, error) {
    ctx := monitor.GetApmContext()
    //ctx := context.Background()
	if len(md) > 0{
        ctx = metadata.NewOutgoingContext(
            ctx,
            metadata.New(md),
        )
    }
    return client.grpcClient.UpdateCustomer(ctx, data.(*customer.UpdateCustomerRequest))
}


