package models

import "time"

type Project struct {
	ID               int               `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectNo        string            `json:"project_no" gorm:"unique"`
	TitleTH          *string           `json:"title_th"`
	TitleEN          *string           `json:"title_en"`
	AbstractText     *string           `json:"abstract_text"`
	AcademicYear     int               `json:"academic_year"`
	Semester         int               `json:"semester"`
	SectionID        *string           `json:"section_id"`
	Program          Program           `json:"program" gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE"`
	ProgramID        int               `json:"program_id"`
	CourseID         int               `json:"course_id"`
	Course           Course            `json:"course" gorm:"foreignKey:CourseID;constraint:OnDelete:SET NULL"`
	Staffs           []Staff           `json:"staffs" gorm:"many2many:project_staffs;constraint:OnDelete:CASCADE;"`
	Members          []Student         `json:"members" gorm:"many2many:project_students;constraint:OnDelete:CASCADE;"`
	ProjectResources []ProjectResource `json:"project_resources"`
	CreatedAt        *time.Time        `json:"created_at" gorm:"default:CURRENT_DATE"`
	UpdatedAt        *time.Time        `json:"updated_at"`
}

type ProjectRequest struct {
	ID               int               `json:"id"`
	ProjectNo        string            `json:"project_no"`
	TitleTH          *string           `json:"title_th"`
	TitleEN          *string           `json:"title_en"`
	AbstractText     *string           `json:"abstract_text"`
	AcademicYear     int               `json:"academic_year"`
	Semester         int               `json:"semester"`
	SectionID        *string           `json:"section_id"`
	ProgramID        int               `json:"program_id"`
	CourseID         int               `json:"course_id"`
	ProjectStaffs    []ProjectStaff    `json:"staffs"`
	Members          []Student         `json:"members"`
	ProjectResources []ProjectResource `json:"project_resources"`
}
