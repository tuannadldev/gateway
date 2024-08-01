// auto generate code, dont edit it
package customer

import(
    "gateway/pkg/monitor"
    "gateway/proto/customer"

    "github.com/gin-gonic/gin"
    "go.elastic.co/apm/v2"
)

type getCustomerHandler struct {
}

func NewGetCustomerHandler() *getCustomerHandler {
    return &getCustomerHandler{}
}

// @Summary permission: 
// @Tags CustomerService
// @Produce json
// @Param id  query  string false " input id"
// @Success 200 {object} customer.GetCustomerResponse
// @Router /api/customers/:id [get]
func (handler *getCustomerHandler) Handle(ctx *gin.Context) (interface{}, error) {
    monitor.SetApmContext(apm.DetachedContext(ctx.Request.Context()))
    data := &customer.GetCustomerRequest{}
    data.Id = ctx.Query("id")
    

    return data, nil
}
