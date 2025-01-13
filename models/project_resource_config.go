package models

type ProjectResourceConfig struct {
	ID             int          `json:"id" gorm:"primaryKey;autoIncrement"`
	Title          string       `json:"title"`
	ResourceTypeID int          `json:"resource_type_id"`
	ResourceType   ResourceType `json:"resource_type" gorm:"foreignKey:ResourceTypeID;constraint:OnDelete:CASCADE"`
	Program        Program      `json:"program" gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE"`
	ProgramID      int          `json:"program_id"`
}
