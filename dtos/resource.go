package dtos

import (
	"mime/multipart"

	"github.com/project-box/models"
)

type Resource struct {
	ID              int            `json:"id"`
	Title           string         `json:"title"`
	ResourceName    *string        `json:"resource_name"`
	Path            *string        `json:"path"`
	URL             string         `json:"url"`
	PDF             *PDF           `json:"pdf"`
	FileExtensionID *int           `json:"file_extension_id"`
	FileExtension   *FileExtension `json:"file_extension"`
	ResourceTypeID  int            `json:"resource_type_id"`
	ResourceType    ResourceType   `json:"resource_type"`
	CreatedAt       string         `json:"created_at"`
}

type ProjectResource struct {
	ID       int      `json:"id" gorm:"primaryKey;autoIncrement"`
	Resource Resource `json:"resource"`
}

type UploadAssetResource struct {
	File          *multipart.FileHeader `form:"file"`
	AssetResource *models.AssetResource `form:"assetResource"`
	Title         string                `form:"title"`
}
