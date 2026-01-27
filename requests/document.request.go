package requests

type DocumentRequest struct {
	Title       string `form:"title" binding:"required,min=5"`
	Description string `form:"description" binding:"required"`
	Type        string `form:"type" binding:"required,oneof=pdf jpg docx"`
}
