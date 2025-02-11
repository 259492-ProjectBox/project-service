package dtos

import (
	"mime/multipart"

	"github.com/project-box/models"
)

type ProjectResourceConfig struct {
	ID             int                 `json:"id"`
	Title          string              `json:"title"`
	IconName       *string             `json:"icon_name"`
	IsActive       bool                `json:"is_active"`
	ResourceTypeID *int                `json:"resource_type_id"`
	ResourceType   models.ResourceType `json:"resource_type"`
	ProgramID      *int                `json:"program_id"`
	Program        models.Program      `json:"program"`
	IconURL        string              `json:"icon_url"`
}

type CreateProjectResourceConfigRequest struct {
	ProjectResourceConfig *models.ProjectResourceConfig `form:"projectResourceConfig"`
	Icon                  *multipart.FileHeader         `form:"icon"`
}
