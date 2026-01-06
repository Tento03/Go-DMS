package controllers

import (
	"fmt"
	"go-dms/config"
	"go-dms/models"
	"go-dms/requests"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetDocument(c *gin.Context) {
	var doc []models.Document
	userId := c.GetUint("userId")
	if err := config.DB.Where("user_id = ?", userId).First(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetched document"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "documents found", "data": doc})
}

func GetDocumentById(c *gin.Context) {
	docId := c.Param("id")
	userId := c.GetUint("userId")

	var doc models.Document
	if err := config.DB.Where("user_id = ? AND id = ?", userId, docId).First(&doc).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to fetched document"})
		return
	}

	docURL := fmt.Sprintf("%s/uploads/%s", c.Request.Host, filepath.Base(doc.Path))

	c.JSON(http.StatusOK, gin.H{
		"message": "document found",
		"data": gin.H{
			"id":          doc.ID,
			"title":       doc.Title,
			"description": doc.Description,
			"type":        doc.Type,
			"url":         "http://" + docURL,
		},
	})
}

func CreateDocument(c *gin.Context) {
	var req requests.CreateDocumentRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
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
	allowed := map[string][]string{
		"pdf":   {".pdf"},
		"image": {".jpg", ".jpeg", ".png"},
		"docx":  {".docx"},
	}

	if !slices.Contains(allowed[req.Type], ext) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed"})
		return
	}

	userId := c.GetUint("userId")

	filename := fmt.Sprintf(
		"%d_%d%s",
		userId,
		time.Now().UnixNano(),
		ext,
	)

	path := filepath.Join("uploads", filename)

	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	doc := models.Document{
		UserID:      userId,
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Path:        path,
	}
	if err := config.DB.Create(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save document"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "document uploaded", "data": doc})
}

func UpdateDocument(c *gin.Context) {
	docId := c.Param("id")
	userId := c.GetUint("userId")

	var doc models.Document
	if err := config.DB.Where("user_id = ? AND id = ?", userId, docId).First(&doc).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to fetch document"})
		return
	}

	var req requests.UpdateDocumentRequest
	if err := c.ShouldBind(&req); err != nil {
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
		allowed := map[string][]string{
			"pdf":   {".pdf"},
			"image": {".jpg", ".jpeg", ".png"},
			"docx":  {".docx"},
		}
		if !slices.Contains(allowed[req.Type], ext) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed"})
			return
		}

		if err := os.Remove(doc.Path); err != nil {
			fmt.Println("failed to remove old file:", err)
		}

		filename := fmt.Sprintf(
			"%d_%d%s",
			userId,
			time.Now().UnixNano(),
			ext,
		)
		path := filepath.Join("uploads", filename)

		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
			return
		}

		doc.Path = path
		doc.Type = req.Type
	}

	doc.Title = req.Title
	doc.Description = req.Description

	if err := config.DB.Save(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "document updated", "data": doc})
}

func DeleteDocument(c *gin.Context) {
	docId := c.Param("id")
	userId := c.GetUint("userId")

	var doc models.Document
	if err := config.DB.Where("user_id = ? AND id = ?", userId, docId).First(&doc).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	if err := config.DB.Delete(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete document"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "document deleted"})
}
