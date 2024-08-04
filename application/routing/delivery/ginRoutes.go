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

    routes = r.Group("/api/orders")
    routes.POST("/", handle)
    routes.POST("/pay/:id", handle)
    routes.PUT("/submit/:id", handle)
    routes.PUT("/cart/:id", handle)
    routes.PUT("/cancel/:id", handle)
    routes.PUT("/complete/:id", handle)
    routes.PUT("/address/:id", handle)
    routes.GET("/:id", handle)
    routes.GET("/search", handle)

}