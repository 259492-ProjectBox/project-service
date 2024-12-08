package dtos

import (
	"mime/multipart"

	"github.com/project-box/models"
)

type CreateProjectRequest struct {
	Files   []*multipart.FileHeader `form:"files"`
	Titles  []string                `form:"titles[]"`
	Project models.Project          `form:"project"`
}

type ProjectData struct {
	ID               int               `json:"id"`
	ProjectNo        string            `json:"project_no"`
	TitleTH          string            `json:"title_th"`
	TitleEN          string            `json:"title_en"`
	AbstractText     string            `json:"abstract_text"`
	AcademicYear     int               `json:"academic_year"`
	Semester         int               `json:"semester"`
	IsApproved       bool              `json:"is_approved"`
	SectionID        string            `json:"section_id"`
	CreatedAt        string            `json:"created_at"`
	MajorID          int               `json:"major_id"`
	Major            Major             `json:"major"`
	CourseID         int               `json:"course_id"`
	Course           Course            `json:"course"`
	Employees        []Employee        `json:"employees"`
	Members          []Student         `json:"members"`
	ProjectResources []ProjectResource `json:"project_resources"`
}
