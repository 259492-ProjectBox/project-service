package models

// ProjectStudent model
type ProjectStudent struct {
	ProjectID int     `json:"project_id" gorm:"not null"`
	StudentID string  `json:"student_id" gorm:"not null"`
	Project   Project `json:"project" gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Student   Student `json:"student" gorm:"foreignKey:StudentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
