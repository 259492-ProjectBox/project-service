package utils

import (
	"github.com/project-box/dtos"
	dto "github.com/project-box/dtos"
	"github.com/project-box/models"
)

func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
func SanitizeProjectMessage(project *models.Project) dto.ProjectData {

	projectData := dto.ProjectData{
		ID:           project.ID,
		ProjectNo:    project.ProjectNo,
		TitleTH:      getStringValue(project.TitleTH),
		TitleEN:      getStringValue(project.TitleEN),
		AbstractText: getStringValue(project.AbstractText),
		AcademicYear: project.AcademicYear,
		Semester:     project.Semester,
		CreatedAt:    project.CreatedAt,
		Major: dtos.Major{
			ID:        project.Major.ID,
			MajorName: project.Major.MajorName,
		},
		Course: dtos.Course{
			ID:         project.Course.ID,
			CourseNo:   project.Course.CourseNo,
			CourseName: project.Course.CourseName,
		},
	}

	for _, committee := range project.Employees {
		projectData.Employees = append(projectData.Employees, dtos.Employee{
			ID:        committee.ID,
			Prefix:    committee.Prefix,
			FirstName: committee.FirstName,
			LastName:  committee.LastName,
			Email:     committee.Email,
		})
	}

	for _, member := range project.Members {
		projectData.Members = append(projectData.Members, dtos.Student{
			ID:        member.ID,
			Prefix:    member.Prefix,
			FirstName: member.FirstName,
			LastName:  member.LastName,
			Email:     member.Email,
		})
	}

	return projectData
}
