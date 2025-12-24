package models

import (
	"time"

	"gorm.io/gorm"
)

type Refresh struct {
	gorm.Model
	UserID       uint   `json:"userId"`
	RefreshToken string `json:"refreshToken" gorm:"uniqueIndex"`
	ExpiresAt    time.Time
	RevokedAt    time.Time
}
