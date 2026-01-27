package services

import (
	"errors"
	"go-dms/models"
	"go-dms/repository"
)

var ErrDocumentNotFound = errors.New("Document not found")

func GetAllDocument() (*[]models.Document, error) {
	document, err := repository.GetAllDocuments()
	if err != nil {
		return nil, ErrDocumentNotFound
	}
	return document, nil
}

func GetByDocumentId(documentId string) (*models.Document, error) {
	document, err := repository.GetByDocumentId(documentId)
	if err != nil {
		return nil, ErrDocumentNotFound
	}
	return document, nil
}

func CreateDocument(documentId string, title string, description string, types string, path string) (*models.Document, error) {
	document := &models.Document{
		DocumentID:  documentId,
		Title:       title,
		Description: description,
		Type:        types,
		Path:        path,
	}

	if err := repository.CreateDocument(document); err != nil {
		return nil, err
	}

	return document, nil
}

func UpdateDocument(documentId string, title string, description string, types string, path string) (*models.Document, error) {
	document, err := repository.GetByDocumentId(documentId)
	if err != nil {
		return nil, ErrDocumentNotFound
	}

	document.Title = title
	document.Description = description
	document.Type = types
	document.Path = path

	if err := repository.UpdateDocument(document); err != nil {
		return nil, err
	}

	return document, nil
}

func DeleteDocument(documentId string) error {
	_, err := repository.GetByDocumentId(documentId)
	if err != nil {
		return ErrDocumentNotFound
	}

	return repository.DeleteDocument(documentId)
}
