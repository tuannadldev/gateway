// auto generate code, dont edit it
package customer

import (
    "gateway/pkg/monitor"
    "gateway/proto/customer"
    
    "github.com/gin-gonic/gin"
    "go.elastic.co/apm/v2"
)

type createCustomerHandler struct {
}

func NewCreateCustomerHandler() *createCustomerHandler {
    return &createCustomerHandler{}
}

// @Summary permission: 
// @Tags CustomerService
// @Produce json
// @Param name  body  string false " customer name"
// @Param email  body  string false " customer email"
// @Param age  body  int32 false " customer age"
// @Param address  body  string false " customer address"
// @Param body  body  customer.CreateCustomerRequest true "Body example"
// @Success 200 {object} customer.CreateCustomerResponse
// @Router /api/customers/ [post]
func (handler *createCustomerHandler) Handle(ctx *gin.Context) (interface{}, error) {
    monitor.SetApmContext(apm.DetachedContext(ctx.Request.Context()))
    data := customer.CreateCustomerRequest{}
    if err := ctx.BindJSON(&data); err != nil {
        return nil, err
    }

    return &data, nil
}
