package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string
	Email     string `gorm:"unique"`
	Username  string `gorm:"unique"`
	Password  string
	Role      string `gorm:"type:varchar(20);default:USER"`
	Status    int    `gorm:"default:1"`
	BirthDate string
	Phone     string
	Gender    string
	Jabatan   string
}
