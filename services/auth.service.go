package services

import (
	"errors"
	"go-dms/models"
	"go-dms/repository"
	"go-dms/utils"

	"github.com/google/uuid"
)

var ErrUsernameExist = errors.New("username is existed")

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
