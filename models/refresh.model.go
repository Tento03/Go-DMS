package models

import (
	"time"
)

type Refresh struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       uint      `json:"userId" gorm:"not null;index"`
	RefreshToken string    `json:"refreshToken" gorm:"uniqueIndex"`
	ExpiresAt    time.Time `gorm:"not null"`
	RevokedAt    *time.Time
	CreatedAt    time.Time
}
