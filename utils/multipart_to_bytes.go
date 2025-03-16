package utils

import (
	"io"
	"log"
	"mime/multipart"
)

// ConvertMultipartFileToBytes converts a multipart.File to a byte slice ([]byte).
func ConvertMultipartFileToBytes(file multipart.File) ([]byte, error) {
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			// Log the error, but don't return it
			// Returning the error would mask the original error
			// which is more important
			log.Printf("error closing file: %v", err)
		}
	}(file) // Ensure the file is closed after reading

	// Read the entire content of the file into a byte slice
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}
