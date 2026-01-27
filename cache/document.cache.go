package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"go-dms/config"
	"go-dms/models"
	"time"
)

func GetCacheListDocument(ctx context.Context, key string) (*[]models.Document, error) {
	val, err := config.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var document []models.Document
	if err := json.Unmarshal([]byte(val), &document); err != nil {
		return nil, err
	}

	return &document, nil
}

func SetCacheDocumentList(ctx context.Context, key string, document *[]models.Document, ttl time.Duration) error {
	bytes, _ := json.Marshal(document)
	return config.Client.Set(ctx, key, bytes, ttl).Err()
}

func GetCacheDocument(ctx context.Context, key string) (*models.Document, error) {
	val, err := config.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var document models.Document
	if err := json.Unmarshal([]byte(val), &document); err != nil {
		return nil, err
	}

	return &document, nil
}

func SetCacheDocument(ctx context.Context, key string, document *models.Document, ttl time.Duration) error {
	bytes, _ := json.Marshal(document)
	return config.Client.Set(ctx, key, bytes, ttl).Err()
}

func DeleteCacheDocument(ctx context.Context, documentId string) error {
	key := fmt.Sprintf("document:documentId=%s", documentId)
	return config.Client.Del(ctx, key).Err()
}
