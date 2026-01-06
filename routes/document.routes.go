package routes

import (
	"go-dms/controllers"
	"go-dms/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func DocumentRoutes(r *gin.Engine) {
	doc := r.Group("/document")
	doc.Use(middleware.RequireAuth, middleware.DocumentLimiter(5, 2*time.Minute))
	{
		doc.GET("/", controllers.GetAll)
		doc.GET("/:id", controllers.GetById)
		doc.POST("/", controllers.CreateDocument)
		doc.PUT("/:id", controllers.UpdateDocument)
		doc.DELETE("/:id", controllers.DeleteDocument)
	}
}
