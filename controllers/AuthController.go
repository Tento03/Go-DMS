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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password invalid"})
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
	}

	if err := config.DB.Create(&rt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "login berhasil",
		"access_token":  accessString,
		"refresh_token": refreshString,
	})
}

func RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := jwt.Parse(body.RefreshToken, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalid"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := uint(claims["id"].(float64))

	var refresh models.Refresh
	if err := config.DB.Where("user_id = ? AND refresh_token = ?", userId, body.RefreshToken).First(&refresh).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userId,
		"username": claims["username"],
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	})
	accessString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "refresh success",
		"access_token": accessString,
	})
}
