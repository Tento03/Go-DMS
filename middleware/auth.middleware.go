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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access token not found"})
		c.Abort()
		return
	}

	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if !token.Valid || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "acces token invalid"})
		c.Abort()
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	c.Set("userId", uint(claims["id"].(float64)))
	c.Set("username", claims["username"])
	c.Next()
}
