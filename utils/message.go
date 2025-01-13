package utils

import (
	"github.com/project-box/dtos"
	"github.com/project-box/models"
)

func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
func SanitizeProjectMessage(project *models.Project) dtos.ProjectData {

	projectData := dtos.ProjectData{
		ID:           project.ID,
		ProjectNo:    project.ProjectNo,
		TitleTH:      getStringValue(project.TitleTH),
		TitleEN:      getStringValue(project.TitleEN),
		AbstractText: getStringValue(project.AbstractText),
		AcademicYear: project.AcademicYear,
		SectionID:    getStringValue(project.SectionID),
		Semester:     project.Semester,
		CreatedAt:    project.CreatedAt.Format("2006-01-02"),
		ProgramID:    project.ProgramID,
		Program: dtos.Program{
			ID:          project.Program.ID,
			ProgramName: project.Program.ProgramName,
		},
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
	}

	for _, staff := range project.Staffs {
		projectData.Staffs = append(projectData.Staffs, dtos.ProjectStaff{
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
			ProjectRole: dtos.ProjectRole{
				ID:       1,
				RoleName: "Advisor",
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
			ProgramID: member.ProgramID,
			Program: dtos.Program{
				ID:          member.Program.ID,
				ProgramName: member.Program.ProgramName,
			},
		})
	}

	for _, projectResource := range project.ProjectResources {
		resourceType := dtos.ResourceType{
			ID:       projectResource.Resource.ResourceTypeID,
			MimeType: projectResource.Resource.ResourceType.MimeType,
		}

		resource := dtos.Resource{
			ID:             projectResource.Resource.ID,
			Title:          projectResource.Resource.Title,
			ResourceName:   projectResource.Resource.ResourceName,
			Path:           projectResource.Resource.Path,
			CreatedAt:      projectResource.Resource.CreatedAt.Format("2006-01-02"),
			ResourceTypeID: projectResource.Resource.ResourceTypeID,
			ResourceType:   resourceType,
		}

		pages := []dtos.PDFPage{}
		if projectResource.Resource.PDF != nil {
			for _, page := range projectResource.Resource.PDF.Pages {
				pageObj := dtos.PDFPage{
					ID:         page.ID,
					PDFID:      page.PDFID,
					PageNumber: page.PageNumber,
					Content:    page.Content,
				}
				pages = append(pages, pageObj)
			}

			resource.PDF = dtos.PDF{
				ID:         projectResource.Resource.PDF.ID,
				ResourceID: projectResource.Resource.PDF.ResourceID,
				Pages:      pages,
			}
		}

		projectData.ProjectResources = append(projectData.ProjectResources, dtos.ProjectResource{
			ID:       projectResource.ID,
			Resource: resource,
		})

	}
	return projectData
}
