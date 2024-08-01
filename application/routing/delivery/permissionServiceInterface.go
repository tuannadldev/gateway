package delivery

import(
    "github.com/gin-gonic/gin"
)

type PermissionService interface {
    Authorize(ctx *gin.Context) bool
}