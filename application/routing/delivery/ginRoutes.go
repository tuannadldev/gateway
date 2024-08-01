// auto generate code, dont edit it
package delivery

import (
    "github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine, handle func(ctx *gin.Context)) {
    var routes *gin.RouterGroup

    routes = r.Group("/api/customers")
    routes.GET("/:id", handle)
    routes.POST("/", handle)
    routes.PUT("", handle)

}