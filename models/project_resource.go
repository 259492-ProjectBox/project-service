package models

import "time"

type ProjectResource struct {
	ID             int          `json:"id" gorm:"primaryKey;autoIncrement"`
	Title          *string      `json:"title"`
	ResourceName   *string      `json:"resource_name"`
	Path           *string      `json:"path"`
	URL            *string      `json:"url"`
	PDF            *PDF         `json:"pdf" gorm:"constraint:OnDelete:CASCADE" swaggerignore:"true"`
	ResourceTypeID int          `json:"resource_type_id" gorm:"not null"`
	ResourceType   ResourceType `json:"resource_type" gorm:"foreignKey:ResourceTypeID;constraint:OnDelete:CASCADE"`
	ProjectID      int          `json:"project_id"`
	Project        Project      `json:"project" gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" swaggerignore:"true"`
	CreatedAt      *time.Time   `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}
