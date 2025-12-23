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
		Username string `json:"username" gorm:"unique"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", body.Username).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(400, gin.H{"error": "password invalid"})
		return
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(15 * time.Minute),
	})
	accessString, _ := accessToken.SignedString(jwtSecret)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(7 * 24 * time.Hour),
	})
	refreshString, _ := refreshToken.SignedString(jwtSecret)

	var refresh models.Refresh
	refresh.RefreshToken = refreshString
	config.DB.Save(&refresh)

	c.JSON(200, gin.H{"message": "Login berhasil", "access token": accessString, "refresh token": refreshString})
}
