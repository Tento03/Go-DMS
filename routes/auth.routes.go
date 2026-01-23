package routes

import (
	"go-dms/controller"
	"go-dms/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/login", controller.Login)
		auth.POST("/refresh", middleware.RequireAuth, controller.Refresh)
		auth.POST("/logout", middleware.RequireAuth, controller.Logout)
	}
}
