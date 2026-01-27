package requests

type PaginationRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit" binding:"min=10"`
}
