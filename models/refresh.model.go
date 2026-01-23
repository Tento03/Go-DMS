package models

import (
	"time"

	"gorm.io/gorm"
)

type Refresh struct {
	gorm.Model
	UserID    string    `gorm:"not null"`
	Token     string    `gorm:"type:varchar(64);not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	RevokedAt *time.Time
}
