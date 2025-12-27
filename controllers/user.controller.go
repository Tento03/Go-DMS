package controllers

import (
	"go-dms/config"
	"go-dms/models"
	"go-dms/requests"
	"go-dms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetAll(c *gin.Context) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no users found"})
		return
	}
	c.JSON(200, gin.H{"message": "users found", "user": users})
}

func GetById(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no user found"})
		return
	}
	c.JSON(200, gin.H{"message": "user found", "user": user})
}

func Create(c *gin.Context) {
	var req requests.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"errors": utils.ValidationError(err)})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)

	user := models.User{
		Name:      req.Name,
		Email:     req.Email,
		Username:  req.Username,
		Password:  string(hashed),
		Role:      "USER",
		Status:    1,
		BirthDate: req.BirthDate,
		Phone:     req.Phone,
		Gender:    req.Gender,
		Jabatan:   req.Jabatan,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "user created"})
}

func Update(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	var req requests.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"errors": utils.ValidationError(err)})
		return
	}

	if err := config.DB.Model(&user).Updates(req).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user updated"})
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
	id := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	var body struct {
		NewPassword string `json:"newPassword" binding:"required,password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"errors": utils.ValidationError(err)})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 12)

	config.DB.Model(&user).
		UpdateColumn("password", string(hashed))

	c.JSON(200, gin.H{"message": "password reset success"})
}

func ChangePassword(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	var body struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"errors": utils.ValidationError(err)})
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(body.OldPassword),
	); err != nil {
		c.JSON(400, gin.H{"error": "old password wrong"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 12)

	config.DB.Model(&user).
		UpdateColumn("password", string(hashed))

	c.JSON(200, gin.H{"message": "password changed"})
}
