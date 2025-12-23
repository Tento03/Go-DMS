package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"-"`
	Role      string    `json:"role" gorm:"type:varchar(20);default:role"`
	Status    string    `json:"status" gorm:"type:varchar(20);default:active"`
	BirthDate time.Time `json:"birthDate"`
	Phone     string    `json:"phone"`
	Gender    string    `json:"gender"`
	Jabatan   string    `json:"jabatan"`
}
