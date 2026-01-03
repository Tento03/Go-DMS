package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"type:varchar(20);default:USER;not null"`
	Status    int    `gorm:"default:1;not null"`
	BirthDate string `gorm:"not null"`
	Phone     string `gorm:"not null"`
	Gender    string `gorm:"not null"`
	Jabatan   string `gorm:"not null"`
}
