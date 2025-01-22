package dtos

import (
	"mime/multipart"

	"github.com/project-box/models"
)

type CreateProjectRequest struct {
	Files   []*multipart.FileHeader `form:"files"`
	Titles  string                  `form:"titles"`
	Urls    string                  `form:"urls"`
	Project *models.ProjectRequest  `form:"project"`
}

type UpdateProjectRequest struct {
	Files   []*multipart.FileHeader `form:"files"`
	Titles  string                  `form:"titles"`
	Urls    string                  `form:"urls"`
	Project *models.ProjectRequest  `form:"project"`
}

type ProjectData struct {
	ID               int                   `json:"id"`
	ProjectNo        string                `json:"project_no"`
	TitleTH          *string               `json:"title_th"`
	TitleEN          *string               `json:"title_en"`
	AbstractText     *string               `json:"abstract_text"`
	AcademicYear     int                   `json:"academic_year"`
	Semester         int                   `json:"semester"`
	SectionID        *string               `json:"section_id"`
	ProgramID        int                   `json:"program_id"`
	Program          Program               `json:"program"`
	CourseID         int                   `json:"course_id"`
	Course           Course                `json:"course"`
	ProjectStaffs    []ProjectStaffMessage `json:"staffs"`
	Members          []Student             `json:"members"`
	ProjectResources []ProjectResource     `json:"project_resources"`
	CreatedAt        string                `json:"created_at"`
	UpdatedAt        string                `json:"updated_at"`
}
