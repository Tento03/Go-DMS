package requests

type UpdateDocumentRequest struct {
	Title       string `form:"title" binding:"required,min=3,max=100"`
	Description string `form:"description" binding:"required,max=255"`
	Type        string `form:"type" binding:"required,oneof=pdf image docx"`
}
