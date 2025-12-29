package routes

import (
	"go-dms/controllers"
	"go-dms/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	user := r.Group("/users")
	user.GET("/", controllers.GetAll)
	user.GET("/:id", controllers.GetById)
	user.POST("/", middleware.CreateUserLimiter(5, 2*time.Minute), controllers.Create)
	user.PUT("/:id", middleware.UpdateUserLimiter(5, 2*time.Minute), controllers.Update)
	user.DELETE("/:id", middleware.DeleteUserLimiter(5, 2*time.Minute), controllers.Delete)
	user.PATCH("users/reset-password/:id", controllers.ResetPassword)
	user.PATCH("/users/change-password/:id", controllers.ChangePassword)
}
