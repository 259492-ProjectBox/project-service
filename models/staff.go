package models

type Staff struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	PrefixTH    string    `json:"prefix_th"`
	PrefixEN    string    `json:"prefix_en"`
	FirstNameTH string    `json:"first_name_th"`
	LastNameTH  string    `json:"last_name_th"`
	FirstNameEN string    `json:"first_name_en"`
	LastNameEN  string    `json:"last_name_en"`
	Email       string    `json:"email" gorm:"uniqueIndex:idx_email_program"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Program     Program   `json:"program" gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE" swaggerignore:"true"`
	ProgramID   int       `json:"program_id" gorm:"uniqueIndex:idx_email_program"`
	Projects    []Project `json:"projects" gorm:"many2many:project_staffs;constraint:OnDelete:CASCADE;" swaggerignore:"true"`
}
