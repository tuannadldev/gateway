package middleware

import (
	"gateway/config"
	"gateway/pkg/monitor"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	r9 "github.com/redis/go-redis/v9"
	"go.elastic.co/apm/v2"
)

func RateLimitByIP(cfg *config.Config, rdb *r9.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		monitor.SetApmContext(apm.DetachedContext(ctx))
		span, ctx := apm.StartSpan(monitor.GetApmContext(), "RateLimitByIP", "middleware")
		defer span.End()

		ip := c.Request.Header.Get("X-Forwarded-For")

		whitelist, err := rdb.SIsMember(ctx, "GATEWAY_WHITELIST", ip).Result()
		if err == nil && whitelist {
			c.Next()
		}

		key := "ratelimit_ip:" + ip
		err = rdb.Incr(ctx, key).Err()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		rateStr, err := rdb.Get(ctx, key).Result()
		if err != nil && rateStr == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		rate, err := strconv.Atoi(rateStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}

		if rate >= cfg.RateLimit.Limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "Too many requests"})
			return
		}

		err = rdb.Expire(ctx, key, time.Duration(cfg.RateLimit.Expire)*time.Second).Err()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}

		c.Next()
	}
}
