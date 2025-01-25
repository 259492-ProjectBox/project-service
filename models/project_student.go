package models

type ProjectStudent struct {
	ProjectID int     `json:"project_id"`
	Project   Project `json:"project" gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	StudentID string  `json:"student_id"`
	Student   Student `json:"student" gorm:"foreignKey:StudentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
