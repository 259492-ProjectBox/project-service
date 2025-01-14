package models

type PDF struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ResourceID int       `json:"resource_id"`
	Pages      []PDFPage `json:"pages" gorm:"constraint:OnDelete:CASCADE"`
}

type PDFPage struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	PDFID      int    `json:"pdf_id"`
	PageNumber int    `json:"page_number"`
	Content    string `json:"content"`
}
