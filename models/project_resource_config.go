package models

type ProjectResourceConfig struct {
	ID             int          `json:"id" gorm:"primaryKey;autoIncrement"`
	Title          string       `json:"title"`
	IconName       *string      `json:"icon_name"`
	IsActive       bool         `json:"is_active"`
	ResourceTypeID *int         `json:"resource_type_id"`
	ResourceType   ResourceType `json:"resource_type" gorm:"foreignKey:ResourceTypeID;constraint:OnDelete:CASCADE"`
	ProgramID      *int         `json:"program_id"`
	Program        Program      `json:"program" gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE"`
}
