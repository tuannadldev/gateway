package middleware

import (
	"fmt"
	r9 "github.com/redis/go-redis/v9"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CanMiddleware(permission_key string, rdb *r9.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if permission_key == "" {
			c.Next()
			return
		}

		user_id := c.GetInt64("user_login_id")
		if user_id < 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Permission denied"})
			return
		}

		ok, err := rdb.SIsMember(c, fmt.Sprintf("user::permission::%d", user_id), permission_key).Result()
		if err != nil || !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "permission denied"})
			return
		}

		c.Next()
	}
}
