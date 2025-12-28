package middleware

import (
	"fmt"
	"go-dms/config"
	"go-dms/requests"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func LoginRateLimiter(maxAttempt int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req requests.Login
		_ = c.ShouldBindJSON(&req)

		ip := c.ClientIP()
		username := req.Username
		key := fmt.Sprintf("rl:login:%s:%s", ip, username)

		count, err := config.Client.Incr(config.Ctx, key).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "redis error",
			})
			return
		}

		if count == 1 {
			config.Client.Expire(config.Ctx, key, window)
		}

		if count > int64(maxAttempt) {
			ttl, _ := config.Client.TTL(config.Ctx, key).Result()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "too many login attemps",
				"retry_after": int(ttl.Seconds()),
				"max_attemps": maxAttempt,
			})
			return
		}
		c.Next()
	}
}
