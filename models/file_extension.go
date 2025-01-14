package models

type FileExtension struct {
	ID            int    `json:"id" gorm:"primaryKey;autoIncrement"`
	ExtensionName string `json:"extension_name"`
	MimeType      string `json:"mime_type" gorm:"unique;"`
}
