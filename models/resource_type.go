package models

type ResourceType struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	MimeType string `json:"mime_type" gorm:"unique"`
}
