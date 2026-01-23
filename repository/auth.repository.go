package repository

import (
	"go-dms/config"
	"go-dms/models"
)

func Register(user *models.Auth) error {
	return config.DB.Create(user).Error
}

func IsUsernameExist(username string) (bool, error) {
	var count int64
	err := config.DB.Model(&models.Auth{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
