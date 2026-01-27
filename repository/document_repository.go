package repository

import (
	"go-dms/config"
	"go-dms/models"
)

func GetAllDocuments() (*[]models.Document, error) {
	var document []models.Document
	err := config.DB.Model(&models.Document{}).Find(&document).Error
	return &document, err
}

func GetByDocumentId(documentId string) (*models.Document, error) {
	var document models.Document
	err := config.DB.Model(&models.Document{}).Where("document_id = ?", documentId).First(&document).Error
	return &document, err
}

func CreateDocument(document *models.Document) error {
	return config.DB.Create(document).Error
}

func UpdateDocument(document *models.Document) error {
	return config.DB.Save(document).Error
}

func DeleteDocument(documentId string) error {
	return config.DB.Model(&models.Document{}).Where("document_id = ?", documentId).Delete(&models.Document{}).Error
}
