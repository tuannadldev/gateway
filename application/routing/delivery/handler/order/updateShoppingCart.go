// auto generate code, dont edit it
package order

import (
    "gateway/pkg/monitor"
    "gateway/proto/order"
    
    "github.com/gin-gonic/gin"
    "go.elastic.co/apm/v2"
)

type updateShoppingCartHandler struct {
}

func NewUpdateShoppingCartHandler() *updateShoppingCartHandler {
    return &updateShoppingCartHandler{}
}

// @Summary permission: 
// @Tags OrderService
// @Produce json
// @Param AggregateID  body  string false "Missing comment! you should fill its for long-term"
// @Param ShopItems  body  []order.ShopItem false "Missing comment! you should fill its for long-term"
// @Param body  body  order.UpdateShoppingCartReq true "Body example"
// @Success 200 {object} order.UpdateShoppingCartRes
// @Router /api/orders/cart/:id [put]
func (handler *updateShoppingCartHandler) Handle(ctx *gin.Context) (interface{}, error) {
    monitor.SetApmContext(apm.DetachedContext(ctx.Request.Context()))
    data := order.UpdateShoppingCartReq{}
    if err := ctx.BindJSON(&data); err != nil {
        return nil, err
    }

    return &data, nil
}
