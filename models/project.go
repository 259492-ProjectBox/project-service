package models

import "time"

type Project struct {
	ID               int               `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectNo        string            `json:"project_no"`
	TitleTH          *string           `json:"title_th"`
	TitleEN          *string           `json:"title_en"`
	AbstractText     *string           `json:"abstract_text"`
	AcademicYear     int               `json:"academic_year"`
	Semester         int               `json:"semester"`
	SectionID        *string           `json:"section_id"`
	Program          Program           `json:"program" gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE"`
	ProgramID        int               `json:"program_id"`
	Staffs           []Staff           `json:"staffs" gorm:"many2many:project_staffs;constraint:OnDelete:CASCADE;"`
	Members          []Student         `json:"members" gorm:"many2many:project_students;constraint:OnDelete:CASCADE;"`
	ProjectResources []ProjectResource `json:"project_resources"`
	Keywords         []Keyword         `json:"keywords" gorm:"many2many:project_keywords;constraint:OnDelete:CASCADE;"`
	IsPublic         bool              `json:"is_public"`
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
	IsPublic         bool              `json:"is_public"`
	ProjectStaffs    []ProjectStaff    `json:"staffs"`
	Keywords         []Keyword         `json:"keywords"`
	Members          []Student         `json:"members"`
	ProjectResources []ProjectResource `json:"project_resources"`
}
