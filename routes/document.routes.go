package routes

import (
	"go-dms/controller"

	"github.com/gin-gonic/gin"
)

func DocumentRoutes(r *gin.Engine) {
	document := r.Group("/document")
	{
		document.GET("/", controller.GetAllDocuments)
		document.GET("/:id", controller.GetByDocumentId)
		document.POST("/", controller.CreateDocument)
		document.PUT("/:id", controller.UpdateDocument)
		document.DELETE("/:id", controller.DeleteDocument)
	}
}
