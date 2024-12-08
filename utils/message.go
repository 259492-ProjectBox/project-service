package utils

import (
	"fmt"

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
		IsApproved:   project.IsApproved,
		SectionID:    getStringValue(project.SectionID),
		Semester:     project.Semester,
		CreatedAt:    project.CreatedAt.Format("2006-01-02"),
		MajorID:      project.MajorID,
		Major: dtos.Major{
			ID:        project.Major.ID,
			MajorName: project.Major.MajorName,
		},
		Course: dtos.Course{
			ID:         project.Course.ID,
			CourseNo:   project.Course.CourseNo,
			CourseName: project.Course.CourseName,
			MajorID:    project.Course.MajorID,
			Major: dtos.Major{
				ID:        project.Course.Major.ID,
				MajorName: project.Course.Major.MajorName,
			},
		},
	}

	for _, committee := range project.Employees {
		projectData.Employees = append(projectData.Employees, dtos.Employee{
			ID:        committee.ID,
			Prefix:    committee.Prefix,
			FirstName: committee.FirstName,
			LastName:  committee.LastName,
			Email:     committee.Email,
			MajorID:   committee.MajorID,
			Major: dtos.Major{
				ID:        committee.Major.ID,
				MajorName: committee.Major.MajorName,
			},
		})
	}

	for _, member := range project.Members {
		projectData.Members = append(projectData.Members, dtos.Student{
			ID:        member.ID,
			Prefix:    member.Prefix,
			FirstName: member.FirstName,
			LastName:  member.LastName,
			Email:     member.Email,
			MajorID:   member.MajorID,
			Major: dtos.Major{
				ID:        member.Major.ID,
				MajorName: member.Major.MajorName,
			},
		})
	}

	for _, projectResource := range project.ProjectResources {
		fmt.Println(projectResource.Resource.ResourceType.MimeType)
		resourceType := dtos.ResourceType{
			ID:       projectResource.Resource.ResourceTypeID,
			MimeType: projectResource.Resource.ResourceType.MimeType,
		}

		resource := dtos.Resource{
			ID:             projectResource.Resource.ID,
			Title:          projectResource.Resource.Title,
			URL:            projectResource.Resource.URL,
			CreatedAt:      projectResource.Resource.CreatedAt.Format("2006-01-02"),
			ResourceTypeID: projectResource.Resource.ResourceTypeID,
			ResourceType:   resourceType,
		}

		projectData.ProjectResources = append(projectData.ProjectResources, dtos.ProjectResource{
			ID:       projectResource.ID,
			Resource: resource,
		})

	}
	return projectData
}
