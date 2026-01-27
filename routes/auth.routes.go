package routes

import (
	"go-dms/controller"
	"go-dms/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/login", middleware.LoginRateLimiter(5, 2*time.Minute), controller.Login)
		auth.POST("/refresh", middleware.RequireAuth, middleware.RefreshTokenRateLimiter(5, 2*time.Minute), controller.Refresh)
		auth.POST("/logout", middleware.RequireAuth, controller.Logout)
	}
}
