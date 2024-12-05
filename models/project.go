package models

import "time"

// Project model
type Project struct {
	ID               int               `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectNo        string            `json:"project_no" gorm:"unique"`
	TitleTH          *string           `json:"title_th"`
	TitleEN          *string           `json:"title_en"`
	Abstract         *string           `json:"abstract"`
	AcademicYear     int               `json:"academic_year" gorm:"not null"`
	Semester         int               `json:"semester"`
	IsApproved       bool              `json:"is_approved"`
	CreatedAt        time.Time         `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	SectionID        string            `json:"section_id"`
	MajorID          int               `json:"major_id" gorm:"not null"`
	Major            Major             `json:"major" gorm:"foreignKey:MajorID;constraint:OnDelete:SET NULL"`
	CourseID         int               `json:"course_id" gorm:"not null"`
	Course           Course            `json:"course" gorm:"foreignKey:CourseID;constraint:OnDelete:SET NULL"`
	Employees        []Employee        `json:"employees" gorm:"many2many:project_employees;constraint:OnDelete:CASCADE;"`
	Members          []Student         `json:"members" gorm:"many2many:project_students;constraint:OnDelete:CASCADE;"`
	ProjectResources []ProjectResource `json:"project_resources"`
}
