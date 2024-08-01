// auto generate code, dont edit it
package order

import(
    "gateway/pkg/monitor"
    "gateway/proto/order"

    "github.com/gin-gonic/gin"
    "go.elastic.co/apm/v2"
)

type getOrderByIDHandler struct {
}

func NewGetOrderByIDHandler() *getOrderByIDHandler {
    return &getOrderByIDHandler{}
}

// @Summary permission: 
// @Tags OrderService
// @Produce json
// @Param AggregateID  query  string false "Missing comment! you should fill its for long-term"
// @Success 200 {object} order.GetOrderByIDRes
// @Router /api/orders/:id [get]
func (handler *getOrderByIDHandler) Handle(ctx *gin.Context) (interface{}, error) {
    monitor.SetApmContext(apm.DetachedContext(ctx.Request.Context()))
    data := &order.GetOrderByIDReq{}
    data.AggregateID = ctx.Query("AggregateID")
    

    return data, nil
}
