// auto generate code, dont edit it
package order

import (
    "gateway/pkg/monitor"
    "gateway/proto/order"
    
    "github.com/gin-gonic/gin"
    "go.elastic.co/apm/v2"
)

type payOrderHandler struct {
}

func NewPayOrderHandler() *payOrderHandler {
    return &payOrderHandler{}
}

// @Summary permission: 
// @Tags OrderService
// @Produce json
// @Param AggregateID  body  string false "Missing comment! you should fill its for long-term"
// @Param Payment  body  order.Payment false "Missing comment! you should fill its for long-term"
// @Param body  body  order.PayOrderReq true "Body example"
// @Success 200 {object} order.PayOrderRes
// @Router /api/orders/pay/:id [post]
func (handler *payOrderHandler) Handle(ctx *gin.Context) (interface{}, error) {
    monitor.SetApmContext(apm.DetachedContext(ctx.Request.Context()))
    data := order.PayOrderReq{}
    if err := ctx.BindJSON(&data); err != nil {
        return nil, err
    }

    return &data, nil
}
