// auto generate code, dont edit it
package order

import(
    "gateway/pkg/monitor"
    "gateway/proto/order"
    "strconv"
    "strings"

    "github.com/gin-gonic/gin"
    "go.elastic.co/apm/v2"
)

type createRequestHandler struct {
}

func NewCreateRequestHandler() *createRequestHandler {
    return &createRequestHandler{}
}

// @Summary permission: 
// @Tags OrderService
// @Produce json
// @Param question_ids  query  string false "value1,value2,value3,..."
// @Param page  query  int64 false "Missing comment! you should fill its for long-term"
// @Param size  query  int64 false "Missing comment! you should fill its for long-term"
// @Success 200 {object} order.CreateOrderRes
// @Router /api/orders/order12/ [get]
func (handler *createRequestHandler) Handle(ctx *gin.Context) (interface{}, error) {
    monitor.SetApmContext(apm.DetachedContext(ctx.Request.Context()))
    data := &order.GetListRequest{}
    question_ids_split := strings.Split(ctx.Query("question_ids"), ",")
    question_ids_value := make([]int64, 0)
    for _, v := range question_ids_split {
        if v == "" {
            continue
        }
        v_int64, err := strconv.ParseInt(v, 10, 64)
        if err != nil {
            return nil, err
        }
        question_ids_value = append(question_ids_value, v_int64)
    }
    data.QuestionIds = question_ids_value

    page_str := ctx.Query("page")
    if page_str != "" {
        page_value, err := strconv.ParseInt(page_str, 10, 64)
        if err != nil {
            return nil, err
        }
        data.Page = page_value
    }

    size_str := ctx.Query("size")
    if size_str != "" {
        size_value, err := strconv.ParseInt(size_str, 10, 64)
        if err != nil {
            return nil, err
        }
        data.Size = size_value
    }


    return data, nil
}
