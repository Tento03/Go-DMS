package middleware

import (
	"go-dms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAuth(c *gin.Context) {
	accessToken, err := c.Cookie("accessToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access token not found"})
		return
	}

	claims, err := utils.ParseToken(accessToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
		return
	}

	userId := claims["user_id"]
	username := claims["username"]

	c.Set("userId", userId)
	c.Set("username", username)
	c.Next()
}
