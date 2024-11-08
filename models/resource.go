package models

import "time"

// Resource model
type Resource struct {
	ID             int          `json:"id" gorm:"primaryKey;autoIncrement"`
	Title          *string      `json:"title"`
	ProjectID      int          `json:"project_id" gorm:"not null"`
	ResourceTypeID int          `json:"resource_type_id" gorm:"not null"`
	URL            string       `json:"url"`
	CreatedAt      time.Time    `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	Project        Project      `json:"project" gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	ResourceType   ResourceType `json:"resource_type" gorm:"foreignKey:ResourceTypeID;constraint:OnDelete:CASCADE"`
}

type FileResponse struct {
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
	URL          string    `json:"url"`
}

type UploadResourceResponse struct {
	Title          *string `json:"title"`
	ProjectID      int     `json:"project_id"`
	ResourceTypeID int     `json:"resource_type_id"`
	URL            string  `json:"url"`
}
