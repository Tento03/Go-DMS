package controllers

import (
	"go-dms/config"
	"go-dms/models"
	"go-dms/requests"
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
	var body requests.CreateUserRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to hash password"})
		return
	}

	user := models.User{
		Name:      body.Name,
		Email:     body.Email,
		Username:  body.Username,
		Password:  string(hashed),
		Phone:     body.Phone,
		Gender:    body.Gender,
		Jabatan:   body.Jabatan,
		BirthDate: body.BirthDate,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "user created", "data": user})
}

func Update(c *gin.Context) {
	var id = c.Param("id")
	var body models.User
	if err := config.DB.First(&body, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	var req requests.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Gender != "" {
		updates["gender"] = req.Gender
	}
	if req.Jabatan != "" {
		updates["jabatan"] = req.Jabatan
	}
	if !req.BirthDate.IsZero() {
		updates["birth_date"] = req.BirthDate
	}

	if err := config.DB.Model(&body).Updates(updates).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	config.DB.First(&body, id)
	c.JSON(200, gin.H{"message": "user updated", "data": body})
}

func Delete(c *gin.Context) {
	var id = c.Param("id")
	if err := config.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "user deleted"})
}
