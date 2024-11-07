package models

// ResourceType model
type ResourceType struct {
	ID           int    `json:"id" gorm:"primaryKey;autoIncrement"`
	ResourceType string `json:"resource_type" gorm:"unique"`
}
