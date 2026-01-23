package services

import (
	"errors"
	"go-dms/models"
	"go-dms/repository"
	"go-dms/utils"
	"time"

	"github.com/google/uuid"
)

var ErrUsernameExist = errors.New("username is existed")
var ErrUsernameNotFound = errors.New("username not found")
var ErrInvalidCredentials = errors.New("password do not matched")

func Register(username string, password string) (*models.Auth, error) {
	exist, err := repository.IsUsernameExist(username)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, ErrUsernameExist
	}

	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.Auth{
		UserID:   uuid.NewString(),
		Username: username,
		Password: hashed,
	}

	if err := repository.Register(user); err != nil {
		return nil, err
	}

	return user, nil
}

func Login(username string, password string) (string, string, error) {
	user, err := repository.FindByUsername(username)
	if err != nil {
		return "", "", ErrUsernameNotFound
	}

	if !utils.ComparePassword(user.Password, password) {
		return "", "", ErrInvalidCredentials
	}

	accessToken, err := utils.GenerateAccessToken(user.UserID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.UserID)
	if err != nil {
		return "", "", err
	}

	hashRT := utils.HashToken(refreshToken)
	refresh := &models.Refresh{
		UserID:    user.UserID,
		Token:     hashRT,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := repository.SaveRefreshToken(refresh); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func Refresh(refreshToken string) (string, string, error) {
	hashRT := utils.HashToken(refreshToken)

	old, err := repository.FindValidRefreshToken(hashRT)
	if err != nil {
		return "", "", err
	}

	if err := repository.RevokeToken(old); err != nil {
		return "", "", err
	}

	newAccessToken, _ := utils.GenerateAccessToken(old.UserID)
	newRefreshToken, _ := utils.GenerateRefreshToken(old.UserID)

	newHashRT := utils.HashToken(newRefreshToken)
	refresh := &models.Refresh{
		UserID:    old.UserID,
		Token:     newHashRT,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := repository.SaveRefreshToken(refresh); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func Logout(refreshToken string) error {
	hashRT := utils.HashToken(refreshToken)

	old, err := repository.FindValidRefreshToken(hashRT)
	if err != nil {
		return err
	}

	return repository.RevokeToken(old)
}
