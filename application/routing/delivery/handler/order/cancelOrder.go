// auto generate code, dont edit it
package order

import (
    "gateway/pkg/monitor"
    "gateway/proto/order"
    
    "github.com/gin-gonic/gin"
    "go.elastic.co/apm/v2"
)

type cancelOrderHandler struct {
}

func NewCancelOrderHandler() *cancelOrderHandler {
    return &cancelOrderHandler{}
}

// @Summary permission: 
// @Tags OrderService
// @Produce json
// @Param AggregateID  body  string false "Missing comment! you should fill its for long-term"
// @Param CancelReason  body  string false "Missing comment! you should fill its for long-term"
// @Param body  body  order.CancelOrderReq true "Body example"
// @Success 200 {object} order.CancelOrderRes
// @Router /api/orders/cancel/:id [put]
func (handler *cancelOrderHandler) Handle(ctx *gin.Context) (interface{}, error) {
    monitor.SetApmContext(apm.DetachedContext(ctx.Request.Context()))
    data := order.CancelOrderReq{}
    if err := ctx.BindJSON(&data); err != nil {
        return nil, err
    }

    return &data, nil
}
