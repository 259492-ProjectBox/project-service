package utils

import (
	"io"
	"mime/multipart"
)

// ConvertMultipartFileToBytes converts a multipart.File to a byte slice ([]byte).
func ConvertMultipartFileToBytes(file multipart.File) ([]byte, error) {
	defer file.Close() // Ensure the file is closed after reading

	// Read the entire content of the file into a byte slice
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}
