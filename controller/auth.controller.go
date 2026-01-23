package controller

import (
	"go-dms/requests"
	"go-dms/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req requests.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.Register(req.Username, req.Password)
	if err != nil {
		if err == services.ErrUsernameExist {
			c.JSON(http.StatusConflict, gin.H{"error": "username existed"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "user created",
		"data":    gin.H{"username": user.Username},
	})
}

func Login(c *gin.Context) {
	var req requests.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := services.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	secured := os.Getenv("APP_ENV") == "production"
	c.SetCookie("accessToken", accessToken, 15*60, "/", "", secured, true)
	c.SetCookie("refreshToken", refreshToken, 7*24*60*60, "/", "", secured, true)

	c.JSON(http.StatusOK, gin.H{
		"success": "true",
		"message": "login success",
		"data":    gin.H{"username": req.Username}})
}
