package repositories

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/project-box/dtos"
	"github.com/project-box/models"
	"github.com/project-box/utils"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	repository[models.Project]
	GetProjectMessageByID(ctx context.Context, id int) (*dtos.ProjectData, error)
	GetProjectByID(ctx context.Context, id int) (*dtos.ProjectData, error)
	GetProjectWithPDFByID(ctx context.Context, id int) (*dtos.ProjectData, error)
	GetProjectsByStudentId(ctx context.Context, studentId string) ([]models.Project, error)
	CheckDuplicateProjectByTitleAndSemester(ctx context.Context, titleTH, titleEN string, academicYear, semester int) (bool, error)
	CreateProjects(ctx context.Context, projectReq []models.ProjectRequest) ([]*dtos.ProjectData, error)
	CreateProjectWithFiles(ctx context.Context, tx *gorm.DB, project *models.ProjectRequest, projectResources []*models.ProjectResource, files []*multipart.FileHeader) (*dtos.ProjectData, error)
	UpdateProjectWithFiles(ctx context.Context, tx *gorm.DB, project *models.ProjectRequest, projectResources []*models.ProjectResource, files []*multipart.FileHeader) (*dtos.ProjectData, error)
	CreateProjectNumber(ctx context.Context, tx *gorm.DB, project *models.ProjectRequest) (*models.ProjectRequest, error)
}

type projectRepositoryImpl struct {
	db                       *gorm.DB
	projectBucketName        string
	fileExtensionRepo        FileExtensionRepository
	projectStaffRepo         ProjectStaffRepository
	resourceRepo             ResourceRepository
	resourceTypeRepo         ResourceTypeRepository
	uploadRepo               UploadRepository
	projectNumberCounterRepo ProjectNumberCounterRepository

	*repositoryImpl[models.Project]
}

func NewProjectRepository(db *gorm.DB, fileExtensionRepo FileExtensionRepository, projectStaffRepo ProjectStaffRepository, projectNumberCounterRepo ProjectNumberCounterRepository, resourceRepo ResourceRepository, resourceTypeRepo ResourceTypeRepository, uploadRepo UploadRepository) ProjectRepository {
	return &projectRepositoryImpl{
		db:                       db,
		projectBucketName:        os.Getenv("MINIO_PROJECT_BUCKET"),
		fileExtensionRepo:        fileExtensionRepo,
		projectStaffRepo:         projectStaffRepo,
		resourceRepo:             resourceRepo,
		resourceTypeRepo:         resourceTypeRepo,
		uploadRepo:               uploadRepo,
		projectNumberCounterRepo: projectNumberCounterRepo,
		repositoryImpl:           newRepository[models.Project](db),
	}
}
func isPDFFile(fileType string) bool { return fileType == "pdf" || fileType == "application/pdf" }

func (r *projectRepositoryImpl) CreateProjectNumber(ctx context.Context, tx *gorm.DB, project *models.ProjectRequest) (*models.ProjectRequest, error) {
	nextProjectNumber, err := r.projectNumberCounterRepo.GetNextProjectNumber(ctx, tx, project.AcademicYear, project.Semester)
	if err != nil {
		return nil, err
	}
	projectNumber := utils.FormatProjectID(project.Semester, project.AcademicYear, nextProjectNumber)
	project.ProjectNo = projectNumber
	return project, nil
}

func (r *projectRepositoryImpl) GetProjectMessageByID(ctx context.Context, id int) (*dtos.ProjectData, error) {
	projectData, err := r.GetProjectWithPDFByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return projectData, nil
}

func (r *projectRepositoryImpl) GetProjectByID(ctx context.Context, id int) (*dtos.ProjectData, error) {
	project := &models.Project{}
	if err := r.db.WithContext(ctx).
		Preload("Program").
		Preload("Staffs.Program").
		Preload("Members.Program").
		Preload("ProjectResources.ResourceType").
		Preload("ProjectResources.FileExtension").
		First(project, "projects.id = ?", id).Error; err != nil {
		return nil, err
	}

	return r.buildProjectData(ctx, project)
}

func (r *projectRepositoryImpl) GetProjectWithPDFByID(ctx context.Context, id int) (*dtos.ProjectData, error) {
	project := &models.Project{}
	if err := r.db.WithContext(ctx).
		Preload("Program").
		Preload("Staffs.Program").
		Preload("Members.Program").
		Preload("ProjectResources.ResourceType").
		Preload("ProjectResources.FileExtension").
		Preload("ProjectResources.PDF.Pages").
		First(project, "projects.id = ?", id).Error; err != nil {
		return nil, err
	}

	return r.buildProjectData(ctx, project)
}

func (r *projectRepositoryImpl) buildProjectData(ctx context.Context, project *models.Project) (*dtos.ProjectData, error) {
	projectData := utils.SanitizeProjectMessage(project)

	for i, projectStaff := range project.Staffs {
		projectStaff, err := r.projectStaffRepo.GetProjectStaffByProjectIdAndStaffId(ctx, project.ID, projectStaff.ID)
		if err != nil {
			return nil, err
		}
		projectData.ProjectStaffs[i].ProjectRole = dtos.ProjectRole{
			ID:         projectStaff.ProjectRole.ID,
			RoleNameTH: projectStaff.ProjectRole.RoleNameTH,
			RoleNameEN: projectStaff.ProjectRole.RoleNameEN,
		}
		projectData.ProjectStaffs[i].ProjectRole.Program = dtos.Program{
			ID:            projectStaff.ProjectRole.ProgramID,
			Abbreviation:  projectStaff.ProjectRole.Program.Abbreviation,
			ProgramNameTH: projectStaff.ProjectRole.Program.ProgramNameTH,
			ProgramNameEN: projectStaff.ProjectRole.Program.ProgramNameEN,
		}
		projectData.ProjectStaffs[i].ProjectRole.ProgramID = projectStaff.ProjectRole.ProgramID
	}

	return projectData, nil
}

func (r *projectRepositoryImpl) GetProjectsByStudentId(ctx context.Context, studentId string) ([]models.Project, error) {
	var projectIds []int
	if err := r.db.WithContext(ctx).
		Model(&models.ProjectStudent{}).
		Where("project_students.student_id = ?", studentId).
		Pluck("project_students.project_id", &projectIds).Error; err != nil {
		return nil, err
	}

	if len(projectIds) == 0 {
		return []models.Project{}, nil
	}

	var projects []models.Project
	if err := r.db.WithContext(ctx).
		Where("id IN ?", projectIds).
		Preload("Program").
		Preload("Staffs.Program").
		Preload("Members.Program").
		Preload("ProjectResources.ResourceType").
		Preload("ProjectResources.FileExtension").
		Find(&projects).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *projectRepositoryImpl) CreateProjectWithFiles(ctx context.Context, tx *gorm.DB, projectReq *models.ProjectRequest, projectResources []*models.ProjectResource, files []*multipart.FileHeader) (*dtos.ProjectData, error) {
	if tx == nil {
		tx = r.db.Begin()
		if tx.Error != nil {
			return nil, tx.Error
		}
	}

	project, err := r.createProject(ctx, tx, projectReq)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	uploadedFilePaths, err := r.handleCreateProjectResources(ctx, tx, project, projectResources, files)
	if err != nil {
		err := r.uploadRepo.DeleteUploadedFiles(ctx, r.projectBucketName, uploadedFilePaths, minio.RemoveObjectOptions{})
		if err != nil {
			return nil, err
		}
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		err := r.uploadRepo.DeleteUploadedFiles(ctx, r.projectBucketName, uploadedFilePaths, minio.RemoveObjectOptions{})
		if err != nil {
			return nil, err
		}
		tx.Rollback()
		return nil, err
	}

	projectData, err := r.GetProjectByID(ctx, project.ID)
	if err != nil {
		err := r.uploadRepo.DeleteUploadedFiles(ctx, r.projectBucketName, uploadedFilePaths, minio.RemoveObjectOptions{})
		if err != nil {
			return nil, err
		}
		tx.Rollback()
		return nil, err
	}

	return projectData, nil
}

func (r *projectRepositoryImpl) UpdateProjectWithFiles(ctx context.Context, tx *gorm.DB, projectReq *models.ProjectRequest, projectResources []*models.ProjectResource, files []*multipart.FileHeader) (*dtos.ProjectData, error) {
	if tx == nil {
		tx = r.db.Begin()
		if tx.Error != nil {
			return nil, tx.Error
		}
	}

	if err := r.deleteProjectAssociations(ctx, tx, projectReq.ID); err != nil {
		tx.Rollback()
		return nil, err
	}

	project, err := r.updateProject(ctx, tx, projectReq)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	uploadedFilePaths, err := r.handleCreateProjectResources(ctx, tx, project, projectResources, files)
	if err != nil {
		err := r.uploadRepo.DeleteUploadedFiles(ctx, r.projectBucketName, uploadedFilePaths, minio.RemoveObjectOptions{})
		if err != nil {
			return nil, err
		}
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		err := r.uploadRepo.DeleteUploadedFiles(ctx, r.projectBucketName, uploadedFilePaths, minio.RemoveObjectOptions{})
		if err != nil {
			return nil, err
		}
		tx.Rollback()
		return nil, err
	}

	projectData, err := r.GetProjectByID(ctx, project.ID)
	if err != nil {
		err := r.uploadRepo.DeleteUploadedFiles(ctx, r.projectBucketName, uploadedFilePaths, minio.RemoveObjectOptions{})
		if err != nil {
			return nil, err
		}
		tx.Rollback()
		return nil, err
	}

	return projectData, nil
}

func (r *projectRepositoryImpl) deleteProjectAssociations(ctx context.Context, tx *gorm.DB, projectID int) error {
	if err := r.deleteProjectStudents(ctx, tx, projectID); err != nil {
		return err
	}

	if err := r.deleteProjectStaffs(ctx, tx, projectID); err != nil {
		return err
	}

	if err := r.deleteProjectResources(ctx, tx, projectID); err != nil {
		return err
	}
	return nil
}

func (r *projectRepositoryImpl) deleteProjectResources(ctx context.Context, tx *gorm.DB, projectID int) error {
	var projectResources []models.ProjectResource
	if err := tx.WithContext(ctx).Where("project_id = ?", projectID).Find(&projectResources).Error; err != nil {
		return err
	}

	for _, projectResource := range projectResources {
		if projectResource.Path != nil {
			// Trim the bucket name from the path
			objectPath := strings.TrimPrefix(*projectResource.Path, r.projectBucketName+"/")
			if err := r.uploadRepo.DeleteUploadedFile(ctx, r.projectBucketName, objectPath, minio.RemoveObjectOptions{}); err != nil {
				return err
			}
		}

		if err := tx.WithContext(ctx).Delete(&projectResource).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *projectRepositoryImpl) deleteProjectStudents(ctx context.Context, tx *gorm.DB, projectID int) error {
	if err := tx.WithContext(ctx).Where("project_id = ?", projectID).Delete(&models.ProjectStudent{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *projectRepositoryImpl) deleteProjectStaffs(ctx context.Context, tx *gorm.DB, projectID int) error {
	if err := tx.WithContext(ctx).Where("project_id = ?", projectID).Delete(&models.ProjectStaff{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *projectRepositoryImpl) CreateProjects(ctx context.Context, projectReqs []models.ProjectRequest) ([]*dtos.ProjectData, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	projects, err := r.createProjectsInTransaction(ctx, tx, projectReqs)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	projectMessages, err := r.getProjectMessages(ctx, projects)
	if err != nil {
		return nil, err
	}

	return projectMessages, nil
}

func (r *projectRepositoryImpl) createProjectsInTransaction(ctx context.Context, tx *gorm.DB, projectReqs []models.ProjectRequest) ([]*models.Project, error) {
	var projects []*models.Project
	for _, projectReq := range projectReqs {
		project, err := r.createProject(ctx, tx, &projectReq)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *projectRepositoryImpl) getProjectMessages(ctx context.Context, projects []*models.Project) ([]*dtos.ProjectData, error) {
	var projectMessages []*dtos.ProjectData
	for _, project := range projects {
		projectMessage, err := r.GetProjectByID(ctx, project.ID)
		if err != nil {
			return nil, err
		}
		projectMessages = append(projectMessages, projectMessage)
	}
	return projectMessages, nil
}

func (r *projectRepositoryImpl) createProject(ctx context.Context, tx *gorm.DB, projectReq *models.ProjectRequest) (*models.Project, error) {
	projectReq, err := r.CreateProjectNumber(ctx, tx, projectReq)
	if err != nil {
		return nil, err
	}

	project := &models.Project{
		ProjectNo:    projectReq.ProjectNo,
		TitleTH:      projectReq.TitleTH,
		TitleEN:      projectReq.TitleEN,
		AbstractText: projectReq.AbstractText,
		AcademicYear: projectReq.AcademicYear,
		SectionID:    projectReq.SectionID,
		Semester:     projectReq.Semester,
		Members:      projectReq.Members,
	}

	if err := tx.WithContext(ctx).Create(project).Error; err != nil {
		return nil, err
	}

	if err := tx.WithContext(ctx).Preload("Program").First(project).Error; err != nil {
		return nil, err
	}

	for _, projectStaff := range projectReq.ProjectStaffs {
		projectStaff.ProjectID = project.ID
		if err := r.projectStaffRepo.CreateProjectStaff(ctx, tx, &projectStaff); err != nil {
			return nil, err
		}
	}

	return project, nil
}

func (r *projectRepositoryImpl) updateProject(ctx context.Context, tx *gorm.DB, projectReq *models.ProjectRequest) (*models.Project, error) {
	project := &models.Project{
		ID:           projectReq.ID,
		ProjectNo:    projectReq.ProjectNo,
		TitleTH:      projectReq.TitleTH,
		TitleEN:      projectReq.TitleEN,
		AbstractText: projectReq.AbstractText,
		AcademicYear: projectReq.AcademicYear,
		SectionID:    projectReq.SectionID,
		Semester:     projectReq.Semester,
		ProgramID:    projectReq.ProgramID,
		Members:      projectReq.Members,
	}

	if err := tx.WithContext(ctx).Save(project).Error; err != nil {
		return nil, err
	}

	if err := tx.WithContext(ctx).Preload("Program").First(project).Error; err != nil {
		return nil, err
	}

	for _, projectStaff := range projectReq.ProjectStaffs {
		projectStaff.ProjectID = project.ID
		if err := r.projectStaffRepo.CreateProjectStaff(ctx, tx, &projectStaff); err != nil {
			return nil, err
		}
	}
	return project, nil
}

func generateUniqueFileName(fileName string) string {
	return fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileName)
}

func buildFilePath(bucketName, programName, projectNo string, title, uniqueFileName string) string {
	projectNo = strings.ReplaceAll(projectNo, "/", "_")
	return fmt.Sprintf("%s/%s/%s/%s/%s",
		bucketName, programName, projectNo, title, uniqueFileName)
}

func buildObjectName(programName, projectNo string, title, uniqueFileName string) string {
	projectNo = strings.ReplaceAll(projectNo, "/", "_")
	return fmt.Sprintf("%s/%s/%s/%s",
		programName, projectNo, title, uniqueFileName)
}

func (r *projectRepositoryImpl) handleCreateProjectResources(ctx context.Context, tx *gorm.DB, project *models.Project, projectResources []*models.ProjectResource, files []*multipart.FileHeader) ([]string, error) {
	var uploadedObjectNames []string

	for _, projectResource := range projectResources {
		if projectResource.URL != nil {
			if err := r.processURL(ctx, tx, project, projectResource); err != nil {
				return nil, err
			}
		} else {
			if len(files) == 0 {
				return nil, fmt.Errorf("not enough files to uploaded")
			}
			if err := r.processFile(ctx, tx, project, projectResource, files[0], &uploadedObjectNames); err != nil {
				return nil, err
			}
			files = files[1:]
		}
	}

	return uploadedObjectNames, nil
}

func (r *projectRepositoryImpl) processFile(ctx context.Context, tx *gorm.DB, project *models.Project, projectResource *models.ProjectResource, file *multipart.FileHeader, uploadedObjectNames *[]string) error {
	if projectResource.Title == nil {
		return fmt.Errorf("title is required")
	}
	title := projectResource.Title
	uniqueFileName := generateUniqueFileName(file.Filename)
	objectName := buildObjectName(project.Program.ProgramNameTH, project.ProjectNo, *title, uniqueFileName)
	filePath := buildFilePath(r.projectBucketName, project.Program.ProgramNameTH, project.ProjectNo, *title, uniqueFileName)

	if err := r.uploadRepo.UploadFile(ctx, r.projectBucketName, objectName, file, minio.PutObjectOptions{}); err != nil {
		return err
	}
	*uploadedObjectNames = append(*uploadedObjectNames, objectName)

	resourceType, err := r.resourceTypeRepo.GetResourceTypeByName(ctx, tx, "file")
	if err != nil {
		return err
	}

	fileExtension, err := r.fileExtensionRepo.GetFileExtension(ctx, tx, file)
	if err != nil {
		return err
	}
	var pdf *models.PDF
	if isPDFFile(fileExtension.MimeType) {
		pdf, err = utils.ReadPdf(file)
		if err != nil {
			return err
		}
	}
	projectResource.ResourceName = &file.Filename
	projectResource.Path = &filePath
	projectResource.PDF = pdf
	projectResource.ResourceTypeID = resourceType.ID
	projectResource.FileExtensionID = &fileExtension.ID
	projectResource.ProjectID = project.ID

	if err := r.resourceRepo.CreateProjectResource(ctx, tx, projectResource); err != nil {
		return err
	}

	return nil
}

func (r *projectRepositoryImpl) processURL(ctx context.Context, tx *gorm.DB, project *models.Project, projectResource *models.ProjectResource) error {
	resourceType, err := r.resourceTypeRepo.GetResourceTypeByName(ctx, tx, "url")
	if err != nil {
		return err
	}

	projectResource.ResourceTypeID = resourceType.ID
	projectResource.ProjectID = project.ID

	if err := r.resourceRepo.CreateProjectResource(ctx, tx, projectResource); err != nil {
		return err
	}

	return nil
}

func (r *projectRepositoryImpl) CheckDuplicateProjectByTitleAndSemester(ctx context.Context, titleTH, titleEN string, academicYear, semester int) (bool, error) {

	var project models.Project
	err := r.db.WithContext(ctx).
		Where("(title_th = ? OR title_en = ?) AND academic_year = ? AND semester = ?", titleTH, titleEN, academicYear, semester).
		First(&project).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
