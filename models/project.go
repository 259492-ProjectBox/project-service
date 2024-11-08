package models

import "time"

// Project model
type Project struct {
	ID                  int         `json:"id" gorm:"primaryKey;autoIncrement"`
	OldProjectNo        *string     `json:"old_project_no"`
	ProjectNo           string      `json:"project_no"`
	TitleTH             *string     `json:"title_th"`
	TitleEN             *string     `json:"title_en"`
	Abstract            *string     `json:"abstract"`
	ProjectStatus       string      `json:"project_status"`
	RelationDescription string      `json:"relation_description"`
	AdvisorID           int         `json:"advisor_id" gorm:"not null"`
	CourseID            int         `json:"course_id" gorm:"not null"`
	SectionID           *int        `json:"section_id"`
	AcademicYear        int         `json:"academic_year" gorm:"not null"`
	Semester            int         `json:"semester"`
	MajorID             int         `json:"major_id" gorm:"not null"`
	CreatedAt           time.Time   `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	Advisor             Employee    `json:"advisor" gorm:"foreignKey:AdvisorID;constraint:OnDelete:SET NULL"`
	Major               Major       `json:"major" gorm:"foreignKey:MajorID;constraint:OnDelete:SET NULL"`
	Course              Course      `json:"course" gorm:"foreignKey:CourseID;constraint:OnDelete:SET NULL"`
	Section             Section     `json:"section" gorm:"foreignKey:SectionID;constraint:OnDelete:SET NULL"`
	Committees          []*Employee `json:"committees" gorm:"many2many:project_employees;constraint:OnDelete:CASCADE;"`
	Members             []*Student  `json:"members" gorm:"many2many:project_students;constraint:OnDelete:CASCADE;"`
}
