package main

import (
	"go-dms/config"
	"go-dms/models"
	"go-dms/routes"
	"go-dms/validators"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
	config.InitRedis()
	config.DB.AutoMigrate(&models.User{}, &models.Refresh{}, &models.Document{})
	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", validators.PasswordValidator)
		v.RegisterValidation("birthdate", validators.BirthDateValidator)
	}
	routes.UserRoutes(r)
	routes.AuthRoutes(r)
	routes.DocumentRoutes(r)
	r.Run(":8080")

}
