package requests

type CreateUserRequest struct {
	Name      string `json:"name" binding:"required,min=3"`
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required,min=3"`
	Password  string `json:"password" binding:"required,password"`
	Role      string `json:"role" binding:"omitempty,oneof=USER ADMIN"`
	Status    int    `json:"status" binding:"omitempty,oneof=1 2"`
	BirthDate string `json:"birthDate" binding:"omitempty,datetime=2006-01-02"`
	Phone     string `json:"phone" binding:"omitempty,numeric,min=10,max=15"`
	Gender    string `json:"gender" binding:"omitempty,oneof=male female"`
	Jabatan   string `json:"jabatan" binding:"omitempty,min=3"`
}
