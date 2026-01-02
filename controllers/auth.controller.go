package controllers

import (
	"fmt"
	"go-dms/config"
	"go-dms/models"
	"go-dms/requests"
	"go-dms/utils"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func Login(c *gin.Context) {
	var req requests.Login

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": utils.ValidationError(err)})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password do not match"})
		return
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	})
	accessString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to signed access token"})
		return
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	refreshString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to signed refresh token"})
		return
	}

	rt := models.Refresh{
		UserID:       user.ID,
		RefreshToken: refreshString,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	if err := config.DB.Create(&rt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add refresh token"})
		return
	}

	keys := fmt.Sprintf("rl:login:%s", c.ClientIP())
	config.Client.Del(config.Ctx, keys)

	secured := os.Getenv("APP_ENV") == "production"
	c.SetCookie("accessToken", accessString, 15*60, "/", "", secured, true)
	c.SetCookie("refreshToken", refreshString, 7*24*60*60, "/", "", secured, true)

	c.JSON(200, gin.H{"message": "login success"})
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

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
		return
	}
	userId := uint(claims["id"].(float64))

	var old models.Refresh
	if err := config.DB.Where("user_id = ? AND refresh_token = ? AND revoked_at IS NULL", userId, refreshToken).First(&old).Error; err != nil {
		now := time.Now()
		config.DB.Model(&models.Refresh{}).Where("user_id = ?", userId).UpdateColumn("revoked_at", &now)

		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token invalid or reused"})
		return
	}

	now := time.Now()
	config.DB.Model(&old).Update("revoked_at", &now)

	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userId,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	})
	newAccessString, err := newAccessToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to signed new access token"})
		return
	}

	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userId,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	newRefreshString, err := newRefreshToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to signed new refresh token"})
		return
	}

	rt := models.Refresh{
		UserID:       userId,
		RefreshToken: newRefreshString,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}
	if err := config.DB.Create(&rt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add refresh token"})
		return
	}

	ip := c.ClientIP()
	hashedRT := utils.HashToken(refreshToken)
	keys := fmt.Sprintf("rl:refresh:%s:%s", ip, hashedRT)
	config.Client.Del(config.Ctx, keys)

	secured := os.Getenv("APP_ENV") == "production"
	c.SetCookie("accessToken", newAccessString, 15*60, "/", "", secured, true)
	c.SetCookie("refreshToken", newRefreshString, 15*60, "/", "", secured, true)

	c.JSON(http.StatusOK, gin.H{"message": "token refreshed"})
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

	secure := os.Getenv("APP_ENV") == "production"

	c.SetCookie("accessToken", "", -1, "/", "", secure, true)
	c.SetCookie("refreshToken", "", -1, "/", "", secure, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successfull"})
}
