package services

import (
	"errors"
	"fmt"
	"go-dms/cache"
	"go-dms/config"
	"go-dms/models"
	"go-dms/repository"
	"go-dms/requests"
	"log"
	"time"
)

var ErrDocumentNotFound = errors.New("Document not found")

func GetAllDocument() (*[]models.Document, error) {
	var req requests.PaginationRequest

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}

	cacheKey := fmt.Sprintf(
		"documents:page=%d:limit=%d",
		req.Page,
		req.Limit,
	)

	documentCache, err := cache.GetCacheListDocument(config.Ctx, cacheKey)
	if err == nil && len(*documentCache) > 0 {
		log.Println("GET ALL DOCUMENTS -> CACHE HIT:", documentCache)
		return documentCache, nil
	}

	log.Println("GET ALL DOCUMENTS -> CACHE MISS:", documentCache)

	document, err := repository.GetAllDocuments()
	if err != nil {
		return nil, ErrDocumentNotFound
	}

	_ = cache.SetCacheDocumentList(config.Ctx, cacheKey, document, 5*time.Minute)

	return document, nil
}

func GetByDocumentId(documentId string) (*models.Document, error) {
	cacheKey := fmt.Sprintf(
		"document:documentId=%s", documentId,
	)

	documentCache, err := cache.GetCacheDocument(config.Ctx, cacheKey)
	if err == nil && documentCache != nil {
		log.Println("GET DOCUMENT BY ID -> CACHE HIT:", documentId)
		return documentCache, nil
	}

	log.Println("GET DOCUMENT BY ID -> CACHE MISS:", documentId)

	document, err := repository.GetByDocumentId(documentId)
	if err != nil {
		return nil, ErrDocumentNotFound
	}

	_ = cache.SetCacheDocument(config.Ctx, cacheKey, document, 5*time.Minute)
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

	_ = cache.DeleteCacheDocument(config.Ctx, documentId)

	return document, nil
}

func DeleteDocument(documentId string) error {
	_, err := repository.GetByDocumentId(documentId)
	if err != nil {
		return ErrDocumentNotFound
	}

	_ = cache.DeleteCacheDocument(config.Ctx, documentId)
	return repository.DeleteDocument(documentId)
}
