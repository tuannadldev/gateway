// auto generate code, dont edit it
package delivery

import (
    "gateway/application/routing/delivery/handler/customer"
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
}
