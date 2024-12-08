package utils

import (
	"fmt"
	"mime/multipart"
	"os"

	"github.com/heussd/pdftotext-go"
	"github.com/project-box/models"
)

func readPDFFile(pdfFile string) ([]byte, error) {
	content, err := os.ReadFile(pdfFile)
	if err != nil {
		return nil, fmt.Errorf("error reading PDF file: %v", err)
	}
	return content, nil
}

func extractTextFromPDF(pdfContent []byte) ([]pdftotext.PdfPage, error) {
	pages, err := pdftotext.Extract(pdfContent)
	if err != nil {
		return nil, fmt.Errorf("error extracting text from PDF: %v", err)
	}
	return pages, nil
}

func createPdfObject(pages []pdftotext.PdfPage) *models.PDF {
	pdf := models.PDF{}

	for i, page := range pages {
		pageObj := models.PDFPage{
			PageNumber: i + 1,
			Content:    page.Content,
		}
		pdf.Pages = append(pdf.Pages, pageObj)
	}

	return &pdf
}

func ReadPdf(file *multipart.FileHeader) (*models.PDF, error) {
	srcForPDF, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer srcForPDF.Close()

	pdf, err := readPdfContent(srcForPDF)
	if err != nil {
		return nil, fmt.Errorf("error reading PDF: %w", err)
	}

	return pdf, nil
}

func readPdfContent(file multipart.File) (*models.PDF, error) {
	contentBytes, err := ConvertMultipartFileToBytes(file)
	if err != nil {
		return nil, err
	}

	pdfPages, err := extractTextFromPDF(contentBytes)
	if err != nil {
		return nil, err
	}

	return createPdfObject(pdfPages), nil
}
