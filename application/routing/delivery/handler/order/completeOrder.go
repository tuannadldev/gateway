// auto generate code, dont edit it
package order

import (
    "gateway/pkg/monitor"
    "gateway/proto/order"
    
    "github.com/gin-gonic/gin"
    "go.elastic.co/apm/v2"
)

type completeOrderHandler struct {
}

func NewCompleteOrderHandler() *completeOrderHandler {
    return &completeOrderHandler{}
}

// @Summary permission: 
// @Tags OrderService
// @Produce json
// @Param AggregateID  body  string false "Missing comment! you should fill its for long-term"
// @Param DeliveryTimestamp  body  string false "Missing comment! you should fill its for long-term"
// @Param body  body  order.CompleteOrderReq true "Body example"
// @Success 200 {object} order.CompleteOrderRes
// @Router /api/orders/complete/:id [put]
func (handler *completeOrderHandler) Handle(ctx *gin.Context) (interface{}, error) {
    monitor.SetApmContext(apm.DetachedContext(ctx.Request.Context()))
    data := order.CompleteOrderReq{}
    if err := ctx.BindJSON(&data); err != nil {
        return nil, err
    }

    return &data, nil
}
