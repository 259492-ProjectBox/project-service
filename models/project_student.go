package models

// ProjectStudent model
type ProjectStudent struct {
	ProjectID int     `json:"project_id" gorm:"not null"`
	StudentID string  `json:"student_id" gorm:"not null"`
	Project   Project `json:"project" gorm:"foreignKey:ProjectID;constraint:OnDelete:SET NULL"`
	Student   Student `json:"student" gorm:"foreignKey:StudentID;constraint:OnDelete:SET NULL"`
}
