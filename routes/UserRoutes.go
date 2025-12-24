package routes

import (
	"go-dms/controllers"
	"go-dms/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.GET("/users", controllers.GetAll)
	r.GET("/users/:id", controllers.GetById)
	r.POST("/users", controllers.Create)
	r.PUT("/users/:id", controllers.Update)
	r.DELETE("/users/:id", controllers.Delete)
	r.PATCH("users/reset-password/:id", controllers.ResetPassword)
	r.PATCH("/users/change-password/:id", controllers.ChangePassword)

	r.POST("/auth/login", controllers.Login)
	r.POST("/auth/refresh", controllers.RefreshToken)
	r.POST("/auth/logout", controllers.Logout)

	auth := r.Group("/auth")
	auth.Use(middleware.RequireAuth)
	auth.GET("/me", func(ctx *gin.Context) {
		userId := ctx.GetUint("userId")
		username := ctx.GetString("username")
		ctx.JSON(200, gin.H{"message": "authenticated", "user Id": userId, "username": username})
	})
}
