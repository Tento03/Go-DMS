package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string `json:"name"`
	Email     string `json:"email" gorm:"unique"`
	Username  string `json:"username" gorm:"unique"`
	Password  string `json:"password"`
	Role      string `json:"role" gorm:"type:varchar(20);default:USER"`
	Status    int    `json:"status" gorm:"type:int(11);default:1"`
	BirthDate string `json:"birthDate"`
	Phone     string `json:"phone"`
	Gender    string `json:"gender"`
	Jabatan   string `json:"jabatan"`
}
