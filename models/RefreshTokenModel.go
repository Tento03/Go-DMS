package models

import "gorm.io/gorm"

type Refresh struct {
	gorm.Model
	RefreshToken string `json:"refreshToken"`
}
