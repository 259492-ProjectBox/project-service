package models

import "time"

type Resource struct {
	ID                int           `json:"id" gorm:"primaryKey;autoIncrement"`
	Title             string        `json:"title"`
	ResourceName      *string       `json:"resource_name"`
	Path              *string       `json:"path"`
	URL               string        `json:"url"`
	PDF               *PDF          `json:"pdf" gorm:"constraint:OnDelete:CASCADE"`
	ProjectResourceID *int          `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	AssetResourceID   *int          `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	FileExtensionID   *int          `json:"file_extension_id"`
	FileExtension     FileExtension `json:"file_extension" gorm:"foreignKey:FileExtensionID;constraint:OnDelete:CASCADE"`
	ResourceTypeID    int           `json:"resource_type_id" gorm:"not null"`
	ResourceType      ResourceType  `json:"resource_type" gorm:"foreignKey:ResourceTypeID;constraint:OnDelete:CASCADE"`
	CreatedAt         time.Time     `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

type ProjectResource struct {
	ID        int      `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectID int      `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Resource  Resource `json:"resource"`
}

type AssetResource struct {
	ID          int      `json:"id" gorm:"primaryKey;autoIncrement"`
	Description string   `json:"description"`
	Program     Program  `json:"program" gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE"`
	ProgramID   int      `json:"program_id"`
	Resource    Resource `json:"resource"`
}
