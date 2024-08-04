package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	r9 "github.com/redis/go-redis/v9"
)

func BucketRateLimiter(rdb *r9.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		limiter := redis_rate.NewLimiter(rdb)
		ip := c.Request.Header.Get("X-Forwarded-For")
		whitelist, err := rdb.SIsMember(ctx, "white:list:ip", ip).Result()
		if err == nil && whitelist {
			c.Next()
		}
		res, err := limiter.Allow(ctx, "bucket:gateway:api:ip"+ip, redis_rate.PerMinute(120))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.Writer.Header().Set("RateLimit-Remaining", strconv.Itoa(res.Remaining))
		if res.Allowed == 0 {
			// We are rate limited.
			seconds := int(res.RetryAfter / time.Second)
			c.Writer.Header().Set("RateLimit-RetryAfter", strconv.Itoa(seconds))
			// Stop processing and return the error.
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "You are rate limited"})
			return
		}
		// Continue processing as normal.
		c.Next()
	}
}
