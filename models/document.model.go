package models

import "gorm.io/gorm"

type Document struct {
	gorm.Model
	DocumentID  string `gorm:"not null"`
	Title       string `gorm:"size:100;not null"`
	Description string `gorm:"size:255;not null"`
	Type        string `gorm:"size:20"`
	Path        string
}
