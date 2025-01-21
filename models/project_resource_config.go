package models

type ProjectResourceConfig struct {
	ID              int           `json:"id" gorm:"primaryKey;autoIncrement"`
	Title           string        `json:"title"`
	IconName        *string       `json:"icon_name"`
	MaxFileSize     *int          `json:"max_file_size"`
	IsActive        bool          `json:"is_active"`
	FileExtensionID *int          `json:"file_extension_id"`
	FileExtension   FileExtension `json:"file_extension" gorm:"foreignKey:FileExtensionID;constraint:OnDelete:CASCADE"`
	ResourceTypeID  *int          `json:"resource_type_id"`
	ResourceType    ResourceType  `json:"resource_type" gorm:"foreignKey:ResourceTypeID;constraint:OnDelete:CASCADE"`
	ProgramID       *int          `json:"program_id"`
	Program         Program       `json:"program" gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE"`
}
