package dtos

import (
	"mime/multipart"

	"github.com/project-box/models"
)

type ProjectResource struct {
	ID             int          `json:"id" gorm:"primaryKey;autoIncrement"`
	Title          *string      `json:"title"`
	ResourceName   *string      `json:"resource_name"`
	Path           *string      `json:"path"`
	URL            *string      `json:"url"`
	PDF            *PDF         `json:"pdf" gorm:"constraint:OnDelete:CASCADE"`
	ResourceTypeID int          `json:"resource_type_id" gorm:"not null"`
	ResourceType   ResourceType `json:"resource_type" gorm:"foreignKey:ResourceTypeID;constraint:OnDelete:CASCADE"`
	ProjectID      int          `json:"project_id"`
	CreatedAt      string       `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}
type UploadAssetResource struct {
	File          *multipart.FileHeader `form:"file"`
	AssetResource *models.AssetResource `form:"asset_resource"`
}
