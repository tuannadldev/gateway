// auto generate code, dont edit it
package order

import(
    "gateway/pkg/monitor"
    "gateway/proto/order"
    "strconv"

    "github.com/gin-gonic/gin"
    "go.elastic.co/apm/v2"
)

type searchHandler struct {
}

func NewSearchHandler() *searchHandler {
    return &searchHandler{}
}

// @Summary permission: 
// @Tags OrderService
// @Produce json
// @Param SearchText  query  string false "Missing comment! you should fill its for long-term"
// @Param Page  query  int64 false "Missing comment! you should fill its for long-term"
// @Param Size  query  int64 false "Missing comment! you should fill its for long-term"
// @Success 200 {object} order.SearchRes
// @Router /api/orders/search [get]
func (handler *searchHandler) Handle(ctx *gin.Context) (interface{}, error) {
    monitor.SetApmContext(apm.DetachedContext(ctx.Request.Context()))
    data := &order.SearchReq{}
    data.SearchText = ctx.Query("SearchText")
    
    Page_str := ctx.Query("Page")
    if Page_str != "" {
        Page_value, err := strconv.ParseInt(Page_str, 10, 64)
        if err != nil {
            return nil, err
        }
        data.Page = Page_value
    }

    Size_str := ctx.Query("Size")
    if Size_str != "" {
        Size_value, err := strconv.ParseInt(Size_str, 10, 64)
        if err != nil {
            return nil, err
        }
        data.Size = Size_value
    }


    return data, nil
}
