package controllers

import (
	"go-dms/config"
	"go-dms/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetAll(c *gin.Context) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "users found", "data": users})
}

func GetById(c *gin.Context) {
	var id = c.Param("id")
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "user found", "data": user})
}

func Create(c *gin.Context) {
	var body models.User
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(500, gin.H{"error": "hashed failed"})
		return
	}
	body.Password = string(hashed)

	if err := config.DB.Create(&body).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "user created", "user": body})
}

func Update(c *gin.Context) {
	var id = c.Param("id")
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	var UpdatedUser models.User
	c.ShouldBindJSON(&UpdatedUser)
	if err := config.DB.Model(&user).Updates(UpdatedUser).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "user updated", "user": UpdatedUser})
}

func Delete(c *gin.Context) {
	var id = c.Param("id")
	if err := config.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "user deleted"})
}

func ResetPassword(c *gin.Context) {
	var id = c.Param("id")
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	var body struct {
		NewPassword string `json:"newPassword"`
	}
	c.ShouldBindJSON(&body)

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)
	if err != nil {
		c.JSON(500, gin.H{"error": "hashing failed"})
		return
	}
	body.NewPassword = string(hashed)

	if err := config.DB.Model(&user).UpdateColumn("password", &body.NewPassword).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "reset password success"})
}

func ChangePassword(c *gin.Context) {
	var id = c.Param("id")
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "user tidak ditemukan"})
		return
	}

	var body struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	c.ShouldBindJSON(&body)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.OldPassword)); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)
	if err != nil {
		c.JSON(500, gin.H{"error": "hashed failed"})
		return
	}

	body.NewPassword = string(hashed)

	if err := config.DB.Model(&user).UpdateColumn("password", body.NewPassword).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "update password berhasil"})
}
