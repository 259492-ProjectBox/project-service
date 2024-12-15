package models

import "time"

type Resource struct {
	ID                int          `json:"id" gorm:"primaryKey;autoIncrement"`
	Title             string       `json:"title"`
	CreatedAt         time.Time    `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	PDF               *PDF         `json:"pdf" gorm:"constraint:OnDelete:CASCADE"`
	ProjectResourceID *int         `json:"-"`
	AssetResourceID   *int         `json:"-"`
	ResourceTypeID    int          `json:"resource_type_id" gorm:"not null"`
	ResourceType      ResourceType `json:"resource_type" gorm:"foreignKey:ResourceTypeID;constraint:OnDelete:CASCADE"`
}

type ProjectResource struct {
	ID        int      `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectID int      `json:"-"`
	Resource  Resource `json:"resource"`
}

type AssetResource struct {
	ID          int      `json:"id" gorm:"primaryKey;autoIncrement"`
	Description string   `json:"description"`
	MajorID     int      `json:"major_id"`
	Major       Major    `json:"major" gorm:"foreignKey:MajorID;constraint:OnDelete:SET NULL"`
	Resource    Resource `json:"resource"`
}

type FileResponse struct {
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
	URL          string    `json:"url"`
}

type UploadResourceResponse struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	ProjectID      int    `json:"project_id"`
	ResourceTypeID int    `json:"resource_type_id"`
	URL            string `json:"url"`
}
