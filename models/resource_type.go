package models

type ResourceType struct {
	ID           int    `json:"id" gorm:"primaryKey;autoIncrement"`
	ResourceType string `json:"resource_type" gorm:"unique"`
}
