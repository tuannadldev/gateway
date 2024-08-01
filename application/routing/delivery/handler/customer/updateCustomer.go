// auto generate code, dont edit it
package customer

import (
    "gateway/pkg/monitor"
    "gateway/proto/customer"
    
    "github.com/gin-gonic/gin"
    "go.elastic.co/apm/v2"
)

type updateCustomerHandler struct {
}

func NewUpdateCustomerHandler() *updateCustomerHandler {
    return &updateCustomerHandler{}
}

// @Summary permission: 
// @Tags CustomerService
// @Produce json
// @Param id  body  string false " name is comment line"
// @Param name  body  string false " customer name"
// @Param email  body  string false " customer email"
// @Param age  body  int32 false " customer age"
// @Param address  body  string false " customer address"
// @Param body  body  customer.UpdateCustomerRequest true "Body example"
// @Success 200 {object} customer.UpdateCustomerResponse
// @Router /api/customers [put]
func (handler *updateCustomerHandler) Handle(ctx *gin.Context) (interface{}, error) {
    monitor.SetApmContext(apm.DetachedContext(ctx.Request.Context()))
    data := customer.UpdateCustomerRequest{}
    if err := ctx.BindJSON(&data); err != nil {
        return nil, err
    }

    return &data, nil
}
