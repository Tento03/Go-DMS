package controllers

import (
	"go-dms/config"
	"go-dms/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func Login(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", body.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(400, gin.H{"error": "password invalid"})
		return
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	})
	accessString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign access token"})
		return
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	refreshString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign refresh token"})
		return
	}

	rt := models.Refresh{
		UserID:       user.ID,
		RefreshToken: refreshString,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		RevokedAt:    nil,
	}

	if err := config.DB.Create(&rt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save rt"})
		return
	}

	c.SetCookie("accessToken", accessString, 15*60, "/", "", true, true)
	c.SetCookie("refreshToken", refreshString, 7*24*60*60, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "login successfull"})
}

func RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}

	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if !token.Valid || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalid"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := uint(claims["id"].(float64))

	var refresh models.Refresh
	if err := config.DB.Where("user_id = ? AND refresh_token = ? AND revoked_at IS NULL", userId, refreshToken).First(&refresh).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}

	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userId,
		"username": claims["username"],
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	})
	newAccessString, _ := newAccessToken.SignedString(jwtSecret)

	c.SetCookie("accessToken", newAccessString, 15*60, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "access token refreshed"})
}

func Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token not found"})
		return
	}

	now := time.Now()
	if err := config.DB.Model(&models.Refresh{}).Where("refresh_token = ?", refreshToken).Update("revoked_at", &now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoked token"})
		return
	}

	c.SetCookie("accessToken", "", -1, "", "", true, true)
	c.SetCookie("refreshToken", "", -1, "", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successfull"})
}
