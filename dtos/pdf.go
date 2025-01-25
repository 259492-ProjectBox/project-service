package dtos

type PDF struct {
	ID                int       `json:"id"`
	ProjectResourceID int       `json:"project_resource_id"`
	Pages             []PDFPage `json:"pages"`
}

type PDFPage struct {
	ID         int    `json:"id"`
	PDFID      int    `json:"pdf_id"`
	PageNumber int    `json:"page_number"`
	Content    string `json:"content"`
}
