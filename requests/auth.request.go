package requests

type AuthRequest struct {
	Username string `json:"username" binding:"required,min=3,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
}
