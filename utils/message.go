package utils

import (
	"time"

	"github.com/project-box/dtos"
	"github.com/project-box/models"
)

func formatTime(value *time.Time) string {
	if value != nil {
		return value.Format("2006-01-02")
	}
	return ""
}

func SanitizeProjectMessage(project *models.Project) dtos.ProjectData {
	if project == nil {
		return dtos.ProjectData{}
	}
	projectMessage := dtos.ProjectData{
		ID:           project.ID,
		ProjectNo:    project.ProjectNo,
		TitleTH:      project.TitleTH,
		TitleEN:      project.TitleEN,
		AbstractText: project.AbstractText,
		AcademicYear: project.AcademicYear,
		SectionID:    project.SectionID,
		Semester:     project.Semester,
		ProgramID:    project.ProgramID,
		Program: dtos.Program{
			ID:          project.Program.ID,
			ProgramName: project.Program.ProgramName,
		},
		CourseID: project.Course.ID,
		Course: dtos.Course{
			ID:         project.Course.ID,
			CourseNo:   project.Course.CourseNo,
			CourseName: project.Course.CourseName,
			ProgramID:  project.Course.ProgramID,
			Program: dtos.Program{
				ID:          project.Course.Program.ID,
				ProgramName: project.Course.Program.ProgramName,
			},
		},
		CreatedAt: formatTime(project.CreatedAt),
		UpdatedAt: formatTime(project.UpdatedAt),
	}

	for _, staff := range project.Staffs {
		projectMessage.ProjectStaffs = append(projectMessage.ProjectStaffs, dtos.ProjectStaffMessage{
			ID:        staff.ID,
			Prefix:    staff.Prefix,
			FirstName: staff.FirstName,
			LastName:  staff.LastName,
			Email:     staff.Email,
			ProgramID: staff.ProgramID,
			Program: dtos.Program{
				ID:          staff.Program.ID,
				ProgramName: staff.Program.ProgramName,
			},
			ProjectRole: dtos.ProjectRole{},
		})
	}
	// Map Members
	for _, member := range project.Members {
		projectMessage.Members = append(projectMessage.Members, dtos.Student{
			ID:        member.ID,
			Prefix:    member.Prefix,
			FirstName: member.FirstName,
			LastName:  member.LastName,
			Email:     member.Email,
		})
	}
	// Map Project Resources
	for _, projectResource := range project.ProjectResources {
		resourceType := dtos.ResourceType{
			ID:       projectResource.Resource.ResourceTypeID,
			TypeName: projectResource.Resource.ResourceType.TypeName,
		}

		resource := dtos.Resource{
			ID:             projectResource.Resource.ID,
			Title:          projectResource.Resource.Title,
			ResourceTypeID: projectResource.Resource.ResourceTypeID,
			ResourceType:   resourceType,
			URL:            projectResource.Resource.URL,
			CreatedAt:      formatTime(projectResource.Resource.CreatedAt),
		}

		if projectResource.Resource.ResourceType.TypeName != "url" {
			resource.ResourceName = projectResource.Resource.ResourceName
			resource.Path = projectResource.Resource.Path
			resource.FileExtension = &dtos.FileExtension{}
			resource.FileExtensionID = projectResource.Resource.FileExtensionID
			resource.FileExtension.ID = projectResource.Resource.FileExtension.ID
			resource.FileExtension.ExtensionName = projectResource.Resource.FileExtension.ExtensionName
			resource.FileExtension.MimeType = projectResource.Resource.FileExtension.MimeType
		}

		if projectResource.Resource.FileExtension != nil && projectResource.Resource.FileExtension.ExtensionName == "pdf" {
			pages := []dtos.PDFPage{}
			for _, page := range projectResource.Resource.PDF.Pages {
				pages = append(pages, dtos.PDFPage{
					ID:         page.ID,
					PDFID:      page.PDFID,
					PageNumber: page.PageNumber,
					Content:    page.Content,
				})
			}
			resource.PDF = &dtos.PDF{
				ID:         projectResource.Resource.PDF.ID,
				ResourceID: projectResource.Resource.PDF.ResourceID,
				Pages:      pages,
			}
		}
		projectMessage.ProjectResources = append(projectMessage.ProjectResources, dtos.ProjectResource{
			ID:       projectResource.ID,
			Resource: resource,
		})
	}
	return projectMessage
}
