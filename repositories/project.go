package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/project-box/models"
	"github.com/project-box/utils"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	repository[models.Project]
	GetProjectByID(ctx context.Context, id int) (*models.Project, error)
	GetProjectWithPDFByID(ctx context.Context, id int) (*models.Project, error)
	GetProjectsByStudentId(ctx context.Context, studentId string) ([]models.Project, error)
	CreateProjectWithFiles(ctx context.Context, project *models.Project, files []*multipart.FileHeader, titles []string, urls []string) (*models.Project, error)
	UpdateProjectWithFiles(ctx context.Context, project *models.Project, files []*multipart.FileHeader, titles []string, urls []string) (*models.Project, error)
	GetProjectsByAdvisorId(ctx context.Context, advisorId int) ([]models.Project, error)
}

type projectRepositoryImpl struct {
	db               *gorm.DB
	minioClient      *minio.Client
	resourceRepo     ResourceRepository
	resourceTypeRepo ResourceTypeRepository
	*repositoryImpl[models.Project]
}

func NewProjectRepository(db *gorm.DB, minioClient *minio.Client, resourceRepo ResourceRepository, resourceTypeRepo ResourceTypeRepository) ProjectRepository {
	return &projectRepositoryImpl{
		db:               db,
		minioClient:      minioClient,
		resourceRepo:     resourceRepo,
		resourceTypeRepo: resourceTypeRepo,
		repositoryImpl:   newRepository[models.Project](db),
	}
}

func (r *projectRepositoryImpl) GetProjectByID(ctx context.Context, id int) (*models.Project, error) {
	project := &models.Project{}
	if err := r.db.WithContext(ctx).
		Preload("Program").
		Preload("Course.Program").
		Preload("Staffs.Program").
		Preload("Members").
		Preload("ProjectResources.Resource.ResourceType").
		First(project, "projects.id = ?", id).Error; err != nil {
		return nil, err
	}
	return project, nil
}

func (r *projectRepositoryImpl) GetProjectWithPDFByID(ctx context.Context, id int) (*models.Project, error) {
	project := &models.Project{}
	if err := r.db.WithContext(ctx).
		Preload("Program").
		Preload("Course.Program").
		Preload("Staffs.Program").
		Preload("Members").
		Preload("ProjectResources.Resource.ResourceType").
		Preload("ProjectResources.Resource.PDF.Pages").
		First(project, "projects.id = ?", id).Error; err != nil {
		return nil, err
	}
	return project, nil
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
		Preload("Course").
		Preload("Staffs.Program").
		Preload("Members").
		Preload("ProjectResources.Resource.ResourceType").
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepositoryImpl) CreateProjectWithFiles(ctx context.Context, project *models.Project, files []*multipart.FileHeader, titles []string, urls []string) (*models.Project, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := r.createProject(ctx, tx, project); err != nil {
		tx.Rollback()
		return nil, err
	}

	uploadedFilePaths, err := r.handleCreateProjectResources(ctx, tx, project, files, titles, urls)
	if err != nil {
		r.deleteUploadedFiles(ctx, uploadedFilePaths)
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		r.deleteUploadedFiles(ctx, uploadedFilePaths)
		tx.Rollback()
		return nil, err
	}

	return r.GetProjectWithPDFByID(ctx, project.ID)
}

func (r *projectRepositoryImpl) UpdateProjectWithFiles(ctx context.Context, project *models.Project, files []*multipart.FileHeader, titles []string, urls []string) (*models.Project, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	if err := r.deleteProjectAssociations(ctx, tx, project.ID); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := r.updateProject(ctx, tx, project); err != nil {
		tx.Rollback()
		return nil, err
	}
	uploadedFilePaths, err := r.handleCreateProjectResources(ctx, tx, project, files, titles, urls)
	if err != nil {

		r.deleteUploadedFiles(ctx, uploadedFilePaths)
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		r.deleteUploadedFiles(ctx, uploadedFilePaths)
		tx.Rollback()
		return nil, err
	}

	return r.GetProjectWithPDFByID(ctx, project.ID)
}

func (r *projectRepositoryImpl) deleteProjectAssociations(ctx context.Context, tx *gorm.DB, projectID int) error {
	if err := r.deleteProjectStudents(ctx, tx, projectID); err != nil {
		return fmt.Errorf("failed to delete project students: %w", err)
	}

	if err := r.deleteProjectStaffs(ctx, tx, projectID); err != nil {
		return fmt.Errorf("failed to delete project staffs: %w", err)
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

func (r *projectRepositoryImpl) createProject(ctx context.Context, tx *gorm.DB, project *models.Project) error {
	if err := tx.WithContext(ctx).Create(project).Error; err != nil {
		return err
	}

	if err := tx.WithContext(ctx).Preload("Program").First(project).Error; err != nil {
		return err
	}

	return nil
}

func (r *projectRepositoryImpl) updateProject(ctx context.Context, tx *gorm.DB, project *models.Project) error {
	if err := tx.WithContext(ctx).Save(project).Error; err != nil {
		return err
	}
	return nil
}

func isPDFFile(fileType string) bool { return fileType == "pdf" || fileType == "application/pdf" }

func (r *projectRepositoryImpl) handleCreateProjectResources(ctx context.Context, tx *gorm.DB, project *models.Project, files []*multipart.FileHeader, titles []string, urls []string) ([]string, error) {
	var uploadedObjectNames []string
	if len(titles) != len(files)+len(urls) {
		return uploadedObjectNames, fmt.Errorf("not enough titles provided for the files")
	}
	for _, file := range files {
		// Pop the first title
		title := titles[0]
		titles = titles[1:]

		uniqueFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		filePath := fmt.Sprintf("%s/%s/%s/%s", os.Getenv("MINIO_PROJECT_BUCKET"), project.Program.ProgramName, title, uniqueFileName)
		objectName := fmt.Sprintf("%s/%s/%s", project.Program.ProgramName, title, uniqueFileName)
		err := r.uploadFileToMinio(ctx, objectName, file)
		if err != nil {
			return uploadedObjectNames, err
		}
		uploadedObjectNames = append(uploadedObjectNames, objectName)

		resourceType, err := r.resourceTypeRepo.GetResourceTypeByName(ctx, tx, "file")
		if err != nil {
			return uploadedObjectNames, err
		}

		fileExtension, err := r.getFileExtension(ctx, tx, file)
		if err != nil {
			return uploadedObjectNames, err
		}

		var pdf *models.PDF
		if isPDFFile(fileExtension.MimeType) {
			if pdf, err = utils.ReadPdf(file); err != nil {
				return uploadedObjectNames, err
			}
		}

		projectResource := &models.ProjectResource{
			ProjectID: project.ID,
		}
		resource := &models.Resource{
			Title:             title,
			ProjectResourceID: nil,
			ResourceName:      &file.Filename,
			Path:              &filePath,
			PDF:               pdf,
			ResourceTypeID:    resourceType.ID,
			FileExtensionID:   &fileExtension.ID,
		}
		if err := r.resourceRepo.CreateProjectResourceAndResource(ctx, tx, projectResource, resource); err != nil {
			return uploadedObjectNames, err
		}
	}

	for _, url := range urls {
		// Pop the first title
		title := titles[0]
		titles = titles[1:]

		resourceType, err := r.resourceTypeRepo.GetResourceTypeByName(ctx, tx, "url")
		if err != nil {
			return uploadedObjectNames, err
		}

		projectResource := &models.ProjectResource{
			ProjectID: project.ID,
		}
		resource := &models.Resource{
			Title:             title,
			ProjectResourceID: nil,
			URL:               url,
			ResourceTypeID:    resourceType.ID,
		}
		if err := r.resourceRepo.CreateProjectResourceAndResource(ctx, tx, projectResource, resource); err != nil {
			return uploadedObjectNames, err
		}
	}

	return uploadedObjectNames, nil
}

func (r *projectRepositoryImpl) uploadFileToMinio(ctx context.Context, objectName string, file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = r.minioClient.PutObject(ctx, os.Getenv("MINIO_PROJECT_BUCKET"), objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *projectRepositoryImpl) getFileExtension(ctx context.Context, tx *gorm.DB, file *multipart.FileHeader) (*models.FileExtension, error) {
	fileType := file.Header.Get("Content-Type")
	if fileType == "" {
		fileType = "unknown"
	}
	fileExtension := &models.FileExtension{}
	if err := tx.WithContext(ctx).Where("mime_type = ?", fileType).First(fileExtension).Error; err != nil {
		return nil, fmt.Errorf("file type not found: %w", err)
	}

	return fileExtension, nil
}

func (r *projectRepositoryImpl) deleteUploadedFiles(ctx context.Context, objectNames []string) {
	for _, objectName := range objectNames {
		err := r.minioClient.RemoveObject(ctx, os.Getenv("MINIO_PROJECT_BUCKET"), objectName, minio.RemoveObjectOptions{})
		if err != nil {
			log.Printf("Failed to delete file from MinIO: %s, error: %v", objectName, err)
		}
	}
}

func (r *projectRepositoryImpl) UpdateProject(ctx context.Context, id int, project *models.Project) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	existingProject := &models.Project{}
	if err := tx.First(existingProject, id).Error; err != nil {
		tx.Rollback()
		return errors.New("project not found")
	}
	if err := tx.Where("project_id = ?", id).Delete(&models.ProjectStaff{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("project_id = ?", id).Delete(&models.ProjectStudent{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	project.ID = id
	if err := tx.Updates(project).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *projectRepositoryImpl) GetProjectsByAdvisorId(ctx context.Context, advisorId int) ([]models.Project, error) {
	var projects []models.Project

	if err := r.db.WithContext(ctx).
		Table("projects as p").
		Select("p.*").
		Joins("JOIN project_staffs as pe ON pe.project_id = p.id").
		Where("pe.staff_id = ?", advisorId).
		Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("failed to get projects for advisor_id %d: %w", advisorId, err)
	}

	return projects, nil
}
