// auto generate code, dont edit it
package order

import (
    "gateway/pkg/monitor"
    "gateway/proto/order"
    
    "github.com/gin-gonic/gin"
    "go.elastic.co/apm/v2"
)

type changeDeliveryAddressHandler struct {
}

func NewChangeDeliveryAddressHandler() *changeDeliveryAddressHandler {
    return &changeDeliveryAddressHandler{}
}

// @Summary permission: 
// @Tags OrderService
// @Produce json
// @Param AggregateID  body  string false "Missing comment! you should fill its for long-term"
// @Param DeliveryAddress  body  string false "Missing comment! you should fill its for long-term"
// @Param body  body  order.ChangeDeliveryAddressReq true "Body example"
// @Success 200 {object} order.ChangeDeliveryAddressRes
// @Router /api/orders/address/:id [put]
func (handler *changeDeliveryAddressHandler) Handle(ctx *gin.Context) (interface{}, error) {
    monitor.SetApmContext(apm.DetachedContext(ctx.Request.Context()))
    data := order.ChangeDeliveryAddressReq{}
    if err := ctx.BindJSON(&data); err != nil {
        return nil, err
    }

    return &data, nil
}
