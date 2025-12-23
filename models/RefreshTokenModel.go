package models

import (
	"time"

	"gorm.io/gorm"
)

type Refresh struct {
	gorm.Model
	UserID       uint
	RefreshToken string `gorm:"uniqueIndex"`
	ExpiresAt    time.Time
}
