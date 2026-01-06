package middleware

import (
	"fmt"
	"go-dms/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateDocumentLimiter(maxAttempt int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		userId := c.GetUint("userId")
		key := fmt.Sprintf("rl:create:%s:%d", ip, userId)

		count, err := config.Client.Incr(config.Ctx, key).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "redis error",
			})
		}

		if count == 1 {
			config.Client.Expire(config.Ctx, key, window)
		}

		if count > int64(maxAttempt) {
			ttl, _ := config.Client.TTL(config.Ctx, key).Result()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":        "too many create document requests attempts",
				"retry_after":  int(ttl.Seconds()),
				"max_attempts": maxAttempt,
			})
			return
		}
		c.Next()
	}
}
