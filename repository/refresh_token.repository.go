package repository

import (
	"go-dms/config"
	"go-dms/models"
	"time"
)

func SaveRefreshToken(refresh *models.Refresh) error {
	return config.DB.Create(refresh).Error
}

func FindValidRefreshToken(token string) (*models.Refresh, error) {
	var refresh models.Refresh
	err := config.DB.Model(&models.Refresh{}).Where("token = ? AND revoked_at IS NULL", token).First(&refresh).Error
	return &refresh, err
}

func RevokeAllUser(userId string) error {
	now := time.Now()
	return config.DB.Model(&models.Refresh{}).Where("user_id = ?", userId).UpdateColumn("revoked_at", &now).Error
}

func RevokeToken(refresh *models.Refresh) error {
	now := time.Now()
	return config.DB.Model(refresh).UpdateColumn("revoked_at", &now).Error
}
