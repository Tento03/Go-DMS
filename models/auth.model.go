package models

import "gorm.io/gorm"

type Auth struct {
	gorm.Model
	UserID   string `gorm:"not null;unique"`
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
}
