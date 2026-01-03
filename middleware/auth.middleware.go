package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func RequireAuth(c *gin.Context) {
	accessToken, err := c.Cookie("accessToken")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access token not found"})
		return
	}

	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if !token.Valid || err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token invalid"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
		return
	}
	userId := uint(claims["id"].(float64))
	username := claims["username"]

	c.Set("userId", userId)
	c.Set("username", username)
	c.Next()
}
