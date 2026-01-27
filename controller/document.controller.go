package controller

import (
	"fmt"
	"go-dms/requests"
	"go-dms/services"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAllDocuments(c *gin.Context) {
	document, err := services.GetAllDocument()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": services.ErrDocumentNotFound})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "documents fetched",
		"data":    document,
	})
}

func GetByDocumentId(c *gin.Context) {
	documentId := c.Param("id")
	document, err := services.GetByDocumentId(documentId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": services.ErrDocumentNotFound})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "document fetched",
		"data":    document,
	})
}

func CreateDocument(c *gin.Context) {
	var req requests.DocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	if file.Size > 5<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size max 5mb"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedType := map[string][]string{
		"pdf":  {".pdf"},
		"jpg":  {".jpg", ".jpeg"},
		"docx": {".docx"},
	}

	if !slices.Contains(allowedType[ext], req.Type) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed"})
		return
	}

	documentId := uuid.NewString()
	filename := fmt.Sprintf(
		"%s_%d%s",
		documentId,
		time.Now().UnixNano(),
		ext,
	)
	path := filepath.Join("/uploads", filename)

	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
		return
	}

	document, err := services.CreateDocument(documentId, req.Title, req.Description, req.Type, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create document"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "document created",
		"data":    document,
	})
}

func UpdateDocument(c *gin.Context) {
	documentId := c.Param("id")

	var req requests.DocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, _ := c.FormFile("file")
	if file != nil {
		if file.Size > 5<<20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file size max 5mb"})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedType := map[string][]string{
			"pdf":  {".pdf"},
			"jpg":  {".jpg", ".jpeg"},
			"docx": {".docx"},
		}
		if !slices.Contains(allowedType[ext], req.Type) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed"})
			return
		}

		filename := fmt.Sprintf(
			"%s_%d%s",
			documentId,
			time.Now().UnixNano(),
			ext,
		)
		path := filepath.Join("/uploads", filename)

		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
			return
		}

		_, err := services.UpdateDocument(documentId, req.Title, req.Description, req.Type, path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update document"})
			return
		}
	}

	document, err := services.GetByDocumentId(documentId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": services.ErrDocumentNotFound})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "document updated",
		"data":    document,
	})
}

func DeleteDocument(c *gin.Context) {
	documentId := c.Param("id")
	if err := services.DeleteDocument(documentId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete document"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "document deleted",
	})
}
