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

func SanitizeProjectMessage(project *models.Project) *dtos.ProjectData {
	if project == nil {
		return nil
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
		IsPublic:     project.IsPublic,
		ProgramID:    project.ProgramID,
		Program: dtos.Program{
			ID:            project.Program.ID,
			Abbreviation:  project.Program.Abbreviation,
			ProgramNameTH: project.Program.ProgramNameTH,
			ProgramNameEN: project.Program.ProgramNameEN,
		},
		CreatedAt: formatTime(project.CreatedAt),
		UpdatedAt: formatTime(project.UpdatedAt),
	}

	for _, staff := range project.Staffs {
		projectMessage.ProjectStaffs = append(projectMessage.ProjectStaffs, dtos.ProjectStaffMessage{
			ID:          staff.ID,
			PrefixTH:    staff.PrefixTH,
			PrefixEN:    staff.PrefixEN,
			FirstNameTH: staff.FirstNameTH,
			LastNameTH:  staff.LastNameTH,
			FirstNameEN: staff.FirstNameEN,
			LastNameEN:  staff.LastNameEN,
			Email:       staff.Email,
			IsActive:    staff.IsActive,
			ProgramID:   staff.ProgramID,
			Program: dtos.Program{
				ID:            staff.Program.ID,
				Abbreviation:  staff.Program.Abbreviation,
				ProgramNameTH: staff.Program.ProgramNameTH,
				ProgramNameEN: staff.Program.ProgramNameEN,
			},
			ProjectRole: dtos.ProjectRole{},
		})
	}

	for _, member := range project.Members {
		projectMessage.Members = append(projectMessage.Members, dtos.Student{
			ID:           member.ID,
			StudentID:    member.StudentID,
			SecLab:       member.SecLab,
			FirstName:    member.FirstName,
			LastName:     member.LastName,
			Email:        member.Email,
			Semester:     member.Semester,
			AcademicYear: member.AcademicYear,
			ProgramID:    member.ProgramID,
			Program: dtos.Program{
				ID:            member.Program.ID,
				Abbreviation:  member.Program.Abbreviation,
				ProgramNameTH: member.Program.ProgramNameTH,
				ProgramNameEN: member.Program.ProgramNameEN,
			},
		})
	}

	for _, keyword := range project.Keywords {
		projectMessage.Keywords = append(projectMessage.Keywords, dtos.Keyword{
			ID:      keyword.ID,
			Keyword: keyword.Keyword,
			Program: dtos.Program{
				ID:            keyword.Program.ID,
				Abbreviation:  keyword.Program.Abbreviation,
				ProgramNameTH: keyword.Program.ProgramNameTH,
				ProgramNameEN: keyword.Program.ProgramNameEN,
			},
			ProgramID: keyword.ProgramID,
		})
	}
	for _, resource := range project.ProjectResources {
		resourceType := dtos.ResourceType{
			ID:       resource.ResourceTypeID,
			TypeName: resource.ResourceType.TypeName,
		}

		projectResource := dtos.ProjectResource{
			ID:              resource.ID,
			ResourceName:    resource.ResourceName,
			Path:            resource.Path,
			Title:           resource.Title,
			ResourceTypeID:  resource.ResourceTypeID,
			ResourceType:    resourceType,
			FileExtensionID: resource.FileExtensionID,
			FileExtension: dtos.FileExtension{
				ID:            resource.FileExtension.ID,
				ExtensionName: resource.FileExtension.ExtensionName,
				MimeType:      resource.FileExtension.MimeType,
			},
			ProjectID: resource.ProjectID,
			CreatedAt: formatTime(resource.CreatedAt),
		}

		if resource.URL != nil {
			projectResource.URL = resource.URL
		}

		if resource.PDF != nil {
			pages := []dtos.PDFPage{}
			for _, page := range resource.PDF.Pages {
				pages = append(pages, dtos.PDFPage{
					ID:         page.ID,
					PDFID:      page.PDFID,
					PageNumber: page.PageNumber,
					Content:    page.Content,
				})
			}
			projectResource.PDF = &dtos.PDF{
				ID:                resource.PDF.ID,
				ProjectResourceID: resource.PDF.ProjectResourceID,
				Pages:             pages,
			}
		}

		projectMessage.ProjectResources = append(projectMessage.ProjectResources, projectResource)
	}

	return &projectMessage
}
