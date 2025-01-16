package dtos

type FileExtension struct {
	ID            int    `json:"id"`
	ExtensionName string `json:"extension_name"`
	MimeType      string `json:"mime_type"`
}
