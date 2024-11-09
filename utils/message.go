package utils

import (
	dto "github.com/project-box/dtos"
	"github.com/project-box/models"
)

func SanitizeProjectMessage(project *models.Project) dto.ProjectData {
	// Create the base ProjectData structure
	projectData := dto.ProjectData{
		ID:                  project.ID,
		OldProjectNo:        *project.OldProjectNo,
		ProjectNo:           project.ProjectNo,
		TitleTH:             *project.TitleTH,
		TitleEN:             *project.TitleEN,
		Abstract:            *project.Abstract,
		ProjectStatus:       project.ProjectStatus,
		RelationDescription: project.RelationDescription,
		AcademicYear:        project.AcademicYear,
		Semester:            project.Semester,
		CreatedAt:           project.CreatedAt,
		Advisor: models.Employee{
			ID:           project.Advisor.ID,
			Prefix:       project.Advisor.Prefix,
			FirstName:    project.Advisor.FirstName,
			LastName:     project.Advisor.LastName,
			Email:        project.Advisor.Email,
			RoleID:       project.Advisor.RoleID,
			Role: models.Role{
				ID:       project.Advisor.Role.ID,
				RoleName: project.Advisor.Role.RoleName,
			},
		},
		Major: models.Major{
			ID:        project.Major.ID,
			MajorName: project.Major.MajorName,
		},
		Course: models.Course{
			ID:         project.Course.ID,
			CourseName: project.Course.CourseName,
		},
	}

	// Loop through Committees and add each to the DTO
	for _, committee := range project.Committees {
		projectData.Committees = append(projectData.Committees, models.Employee{
			ID:        committee.ID,
			Prefix:    committee.Prefix,
			FirstName: committee.FirstName,
			LastName:  committee.LastName,
			Email:     committee.Email,
			RoleID:    committee.RoleID,
			Role: models.Role{
				ID:       committee.Role.ID,
				RoleName: committee.Role.RoleName,
			},
		})
	}

	// Loop through Members and add each to the DTO
	for _, member := range project.Members {
		projectData.Members = append(projectData.Members, models.Student{
			ID:        member.ID,
			Prefix:    member.Prefix,
			FirstName: member.FirstName,
			LastName:  member.LastName,
			Email:     member.Email,
			Major: models.Major{
				ID:        member.Major.ID,
				MajorName: member.Major.MajorName,
			},
		})
	}

	return projectData
}
