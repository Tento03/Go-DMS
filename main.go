package main

import (
	"go-dms/config"
	"go-dms/models"
	"go-dms/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{}, &models.Refresh{})
	r := gin.Default()
	routes.UserRoutes(r)
	r.Run(":8080")

}
