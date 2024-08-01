// auto generate code, dont edit it
package order

import (
    "gateway/pkg/monitor"
    "gateway/proto/order"
    
    "github.com/gin-gonic/gin"
    "go.elastic.co/apm/v2"
)

type submitOrderHandler struct {
}

func NewSubmitOrderHandler() *submitOrderHandler {
    return &submitOrderHandler{}
}

// @Summary permission: 
// @Tags OrderService
// @Produce json
// @Param AggregateID  body  string false "Missing comment! you should fill its for long-term"
// @Param body  body  order.SubmitOrderReq true "Body example"
// @Success 200 {object} order.SubmitOrderRes
// @Router /api/orders/submit/:id [put]
func (handler *submitOrderHandler) Handle(ctx *gin.Context) (interface{}, error) {
    monitor.SetApmContext(apm.DetachedContext(ctx.Request.Context()))
    data := order.SubmitOrderReq{}
    if err := ctx.BindJSON(&data); err != nil {
        return nil, err
    }

    return &data, nil
}
