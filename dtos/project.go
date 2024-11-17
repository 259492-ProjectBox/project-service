package dtos

import (
	"time"
)

type ProjectData struct {
	ID                  int        `json:"id"`
	OldProjectNo        string     `json:"old_project_no"`
	ProjectNo           string     `json:"project_no"`
	TitleTH             string     `json:"title_th"`
	TitleEN             string     `json:"title_en"`
	Abstract            string     `json:"abstract"`
	ProjectStatus       string     `json:"project_status"`
	RelationDescription string     `json:"relation_description"`
	AcademicYear        int        `json:"academic_year"`
	Semester            int        `json:"semester"`
	CreatedAt           time.Time  `json:"created_at"`
	Advisor             Employee   `json:"advisor"`
	Major               Major      `json:"major"`
	Course              Course     `json:"course"`
	Section             Section    `json:"section"`
	Committees          []Employee `json:"committees"`
	Members             []Student  `json:"members"`
}
