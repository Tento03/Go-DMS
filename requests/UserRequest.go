package requests

import "time"

type UpdateUserRequest struct {
	Name      string    `json:"name"`
	Email     string    `json:"email" binding:"omitempty,email"`
	Username  string    `json:"username"`
	Phone     string    `json:"phone"`
	Gender    string    `json:"gender"`
	Jabatan   string    `json:"jabatan"`
	BirthDate time.Time `json:"birth_date"`
}

type CreateUserRequest struct {
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	Username  string    `json:"username" binding:"required"`
	Password  string    `json:"password" binding:"required,min=8"`
	Phone     string    `json:"phone"`
	Gender    string    `json:"gender"`
	Jabatan   string    `json:"jabatan"`
	BirthDate time.Time `json:"birth_date"`
}
