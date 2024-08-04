// auto generate code, dont edit it
package delivery

import (
    "gateway/application/routing/delivery/handler/customer"
    "gateway/application/routing/delivery/handler/order"
)

type routingConfig struct {
    requestHandler          RequestHandler
    remoteServiceName       string
    remoteServiceMethod     string
    remoteServicePermission string
}

var routingRegistry = map[string]routingConfig{
    "GET:/api/customers/:id": {
        customer.NewGetCustomerHandler(),
        "CustomerService",
        "GetCustomer",
        "",
    },
    "POST:/api/customers/": {
        customer.NewCreateCustomerHandler(),
        "CustomerService",
        "CreateCustomer",
        "",
    },
    "PUT:/api/customers": {
        customer.NewUpdateCustomerHandler(),
        "CustomerService",
        "UpdateCustomer",
        "",
    },
    "POST:/api/orders/": {
        order.NewCreateOrderHandler(),
        "OrderService",
        "CreateOrder",
        "",
    },
    "POST:/api/orders/pay/:id": {
        order.NewPayOrderHandler(),
        "OrderService",
        "PayOrder",
        "",
    },
    "PUT:/api/orders/submit/:id": {
        order.NewSubmitOrderHandler(),
        "OrderService",
        "SubmitOrder",
        "",
    },
    "PUT:/api/orders/cart/:id": {
        order.NewUpdateShoppingCartHandler(),
        "OrderService",
        "UpdateShoppingCart",
        "",
    },
    "PUT:/api/orders/cancel/:id": {
        order.NewCancelOrderHandler(),
        "OrderService",
        "CancelOrder",
        "",
    },
    "PUT:/api/orders/complete/:id": {
        order.NewCompleteOrderHandler(),
        "OrderService",
        "CompleteOrder",
        "",
    },
    "PUT:/api/orders/address/:id": {
        order.NewChangeDeliveryAddressHandler(),
        "OrderService",
        "ChangeDeliveryAddress",
        "",
    },
    "GET:/api/orders/:id": {
        order.NewGetOrderByIDHandler(),
        "OrderService",
        "GetOrderByID",
        "",
    },
    "GET:/api/orders/search": {
        order.NewSearchHandler(),
        "OrderService",
        "Search",
        "",
    },
}
