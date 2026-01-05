package controllers

import (
	"fmt"
	"go-dms/config"
	"go-dms/models"
	"go-dms/requests"
	"net/http"
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
	var doc models.Document
	docId := c.Param("id")
	userId := c.GetUint("userId")
	if err := config.DB.Where("user_id = ? AND id = ?", userId, docId).First(&doc).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to fetched document"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "document found", "data": doc})
}

func CreateDocument(c *gin.Context) {
	//1. Bind JSON
	var req requests.CreateDocumentRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	//2. Ambil dari form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	//3. Validasi ukuran file
	if file.Size > 5<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size max 5mb"})
		return
	}

	//4. Validasi extensi file
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

	//5. Ambil UserId
	userId := c.GetUint("userId")

	//6. Generate nama file aman
	filename := fmt.Sprintf(
		"%d_%d%s",
		userId,
		time.Now().UnixNano(),
		ext,
	)

	path := filepath.Join("uploads", filename)

	//7. Simpan file
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	//8. Simpan ke DB
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

	//9. Hasilkan response
	c.JSON(http.StatusCreated, gin.H{"message": "document uploaded", "data": doc})
}
