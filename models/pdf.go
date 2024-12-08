package models

// PDF represents the entire PDF document.
type PDF struct {
	Pages []PDFPage `json:"pages"` // List of pages in the PDF
}

// PDFPage represents a single page in the PDF.
type PDFPage struct {
	PageNumber int    `json:"page_number"` // The page number
	Content    string `json:"content"`     // The text or content of the page
}
