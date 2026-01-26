package middleware

import (
	"fmt"
	"go-dms/config"
	"go-dms/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func LoginRateLimiter(maxAttempt int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("login:rl:%s", ip)

		count, err := config.Client.Incr(config.Ctx, key).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "redis error",
			})
			return
		}

		if count == 1 {
			config.Client.Expire(config.Ctx, key, duration)
		}
		if count > int64(maxAttempt) {
			ttl, _ := config.Client.TTL(config.Ctx, key).Result()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":        "too many login requests",
				"retry_at":     int(ttl.Seconds()),
				"max_attempts": maxAttempt,
			})
			return
		}
		c.Next()
	}
}

func RefreshTokenRateLimiter(maxAttempt int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		refreshToken, err := c.Cookie("refreshToken")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "refresh token not found",
			})
			return
		}
		hashRT := utils.HashToken(refreshToken)
		key := fmt.Sprintf("refresh:rl:%s:%s", ip, hashRT)

		count, err := config.Client.Incr(config.Ctx, key).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "redis error",
			})
			return
		}

		if count == 1 {
			config.Client.Expire(config.Ctx, key, duration)
		}
		if count > int64(maxAttempt) {
			ttl, _ := config.Client.TTL(config.Ctx, key).Result()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":        "too many refresh token attempts",
				"retry_at":     int(ttl.Seconds()),
				"max_attempts": maxAttempt,
			})
			return
		}
		c.Next()
	}
}
