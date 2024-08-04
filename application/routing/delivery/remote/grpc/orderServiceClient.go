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
    "gateway/proto/order"
)

func (client *orderServiceClient) initMethodRegistry() {
    client.methodRegistry = map[string]func(interface{}, map[string]string) (interface{}, error) {
        "CreateOrder": client.createOrder,
        "PayOrder": client.payOrder,
        "SubmitOrder": client.submitOrder,
        "UpdateShoppingCart": client.updateShoppingCart,
        "CancelOrder": client.cancelOrder,
        "CompleteOrder": client.completeOrder,
        "ChangeDeliveryAddress": client.changeDeliveryAddress,
        "GetOrderByID": client.getOrderByID,
        "Search": client.search,
    }
}

type orderServiceClient struct {
    grpcClient order.OrderServiceClient
    methodRegistry map[string]func(interface{}, map[string]string) (interface{}, error)
}

func NewOrderServiceClient(config *config.Config) *orderServiceClient {
    // using WithInsecure() because no SSL running
    cc, err := grpc.Dial(
        config.Service.OrderServiceUrl,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithUnaryInterceptor(apmgrpc.NewUnaryClientInterceptor()),
        grpc.WithStreamInterceptor(apmgrpc.NewStreamClientInterceptor()),
    )

    if err != nil {
        fmt.Println("Could not connect:", err)
    }

    client := orderServiceClient{
        grpcClient: order.NewOrderServiceClient(cc),
    }
    client.initMethodRegistry()

    return &client
}

func (client *orderServiceClient) GetMethodRegistry() map[string]func(interface{}, map[string]string) (interface{}, error) {
    return client.methodRegistry
}

func (client *orderServiceClient) createOrder(data interface{}, md map[string]string) (interface{}, error) {
    ctx := monitor.GetApmContext()
    //ctx := context.Background()
	if len(md) > 0{
        ctx = metadata.NewOutgoingContext(
            ctx,
            metadata.New(md),
        )
    }
    return client.grpcClient.CreateOrder(ctx, data.(*order.CreateOrderReq))
}


func (client *orderServiceClient) payOrder(data interface{}, md map[string]string) (interface{}, error) {
    ctx := monitor.GetApmContext()
    //ctx := context.Background()
	if len(md) > 0{
        ctx = metadata.NewOutgoingContext(
            ctx,
            metadata.New(md),
        )
    }
    return client.grpcClient.PayOrder(ctx, data.(*order.PayOrderReq))
}


func (client *orderServiceClient) submitOrder(data interface{}, md map[string]string) (interface{}, error) {
    ctx := monitor.GetApmContext()
    //ctx := context.Background()
	if len(md) > 0{
        ctx = metadata.NewOutgoingContext(
            ctx,
            metadata.New(md),
        )
    }
    return client.grpcClient.SubmitOrder(ctx, data.(*order.SubmitOrderReq))
}


func (client *orderServiceClient) updateShoppingCart(data interface{}, md map[string]string) (interface{}, error) {
    ctx := monitor.GetApmContext()
    //ctx := context.Background()
	if len(md) > 0{
        ctx = metadata.NewOutgoingContext(
            ctx,
            metadata.New(md),
        )
    }
    return client.grpcClient.UpdateShoppingCart(ctx, data.(*order.UpdateShoppingCartReq))
}


func (client *orderServiceClient) cancelOrder(data interface{}, md map[string]string) (interface{}, error) {
    ctx := monitor.GetApmContext()
    //ctx := context.Background()
	if len(md) > 0{
        ctx = metadata.NewOutgoingContext(
            ctx,
            metadata.New(md),
        )
    }
    return client.grpcClient.CancelOrder(ctx, data.(*order.CancelOrderReq))
}


func (client *orderServiceClient) completeOrder(data interface{}, md map[string]string) (interface{}, error) {
    ctx := monitor.GetApmContext()
    //ctx := context.Background()
	if len(md) > 0{
        ctx = metadata.NewOutgoingContext(
            ctx,
            metadata.New(md),
        )
    }
    return client.grpcClient.CompleteOrder(ctx, data.(*order.CompleteOrderReq))
}


func (client *orderServiceClient) changeDeliveryAddress(data interface{}, md map[string]string) (interface{}, error) {
    ctx := monitor.GetApmContext()
    //ctx := context.Background()
	if len(md) > 0{
        ctx = metadata.NewOutgoingContext(
            ctx,
            metadata.New(md),
        )
    }
    return client.grpcClient.ChangeDeliveryAddress(ctx, data.(*order.ChangeDeliveryAddressReq))
}


func (client *orderServiceClient) getOrderByID(data interface{}, md map[string]string) (interface{}, error) {
    ctx := monitor.GetApmContext()
    //ctx := context.Background()
	if len(md) > 0{
        ctx = metadata.NewOutgoingContext(
            ctx,
            metadata.New(md),
        )
    }
    return client.grpcClient.GetOrderByID(ctx, data.(*order.GetOrderByIDReq))
}


func (client *orderServiceClient) search(data interface{}, md map[string]string) (interface{}, error) {
    ctx := monitor.GetApmContext()
    //ctx := context.Background()
	if len(md) > 0{
        ctx = metadata.NewOutgoingContext(
            ctx,
            metadata.New(md),
        )
    }
    return client.grpcClient.Search(ctx, data.(*order.SearchReq))
}


