package delivery

import(
    "github.com/gin-gonic/gin"
)

type RequestHandler interface {
    Handle(ctx *gin.Context) (interface{}, error)
}