package main

import (
	"go-dms/config"
	"go-dms/models"
	"go-dms/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.InitRedis()
	config.ConnectDB()
	config.DB.AutoMigrate(&models.Auth{}, &models.Refresh{}, &models.Document{})
	r := gin.Default()
	routes.AuthRoutes(r)
	routes.DocumentRoutes(r)
	r.Run(":8080")
}
