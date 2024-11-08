package dtos

import (
	"time"

	"github.com/project-box/models"
)

type ProjectData struct {
	ID                  int               `json:"id"`
	OldProjectNo        string            `json:"old_project_no"`
	ProjectNo           string            `json:"project_no"`
	TitleTH             string            `json:"title_th"`
	TitleEN             string            `json:"title_en"`
	Abstract            string            `json:"abstract"`
	ProjectStatus       string            `json:"project_status"`
	RelationDescription string            `json:"relation_description"`
	AcademicYear        int               `json:"academic_year"`
	Semester            int               `json:"semester"`
	CreatedAt           time.Time         `json:"created_at"`
	Advisor             models.Employee   `json:"advisor"`
	Major               models.Major      `json:"major"`
	Course              models.Course     `json:"course"`
	Committees          []models.Employee `json:"committees"`
	Members             []models.Student  `json:"members"`
}
