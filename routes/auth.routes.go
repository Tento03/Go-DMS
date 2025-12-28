package routes

import (
	"go-dms/controllers"
	"go-dms/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")

	auth.POST("/login", middleware.LoginRateLimiter(5, 2*time.Minute), controllers.Login)
	auth.POST("/refresh", controllers.RefreshToken)
	auth.POST("/logout", controllers.Logout)
	auth.GET("/me", middleware.RequireAuth, func(ctx *gin.Context) {
		userId := ctx.GetUint("userId")
		username := ctx.GetString("username")
		ctx.JSON(200, gin.H{"message": "authenticated", "userId": userId, "username": username})
	})
}
