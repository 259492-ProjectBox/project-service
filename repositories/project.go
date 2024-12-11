package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
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
	CreateProjectWithFiles(ctx context.Context, project *models.Project, files []*multipart.FileHeader, titles []string) (*models.Project, error)
	UpdateProject(ctx context.Context, id int, project *models.Project) error
}

type projectRepositoryImpl struct {
	db          *gorm.DB
	minioClient *minio.Client
	*repositoryImpl[models.Project]
}

func NewProjectRepository(db *gorm.DB, minioClient *minio.Client) ProjectRepository {
	return &projectRepositoryImpl{
		db:             db,
		minioClient:    minioClient,
		repositoryImpl: newRepository[models.Project](db),
	}
}

func (r *projectRepositoryImpl) GetProjectByID(ctx context.Context, id int) (*models.Project, error) {
	project := &models.Project{}
	if err := r.db.WithContext(ctx).
		Preload("Major").
		Preload("Course.Major").
		Preload("Employees.Major").
		Preload("Members.Major").
		Preload("ProjectResources.Resource.ResourceType").
		First(project, "projects.id = ?", id).Error; err != nil {
		return nil, err
	}
	return project, nil
}

func (r *projectRepositoryImpl) GetProjectWithPDFByID(ctx context.Context, id int) (*models.Project, error) {
	project := &models.Project{}
	if err := r.db.WithContext(ctx).
		Preload("Major").
		Preload("Course.Major").
		Preload("Employees.Major").
		Preload("Members.Major").
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
		Preload("Major").
		Preload("Course").
		Preload("Employees.Major").
		Preload("Members.Major").
		Preload("ProjectResources.Resource.ResourceType").
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepositoryImpl) CreateProjectWithFiles(ctx context.Context, project *models.Project, files []*multipart.FileHeader, titles []string) (*models.Project, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := r.createProject(ctx, tx, project); err != nil {
		return nil, err
	}

	uploadedFilePaths, err := r.handleFiles(ctx, tx, project.ID, files, titles)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		r.deleteUploadedFiles(ctx, uploadedFilePaths)
		return nil, err
	}

	return r.GetProjectWithPDFByID(ctx, project.ID)
}

func (r *projectRepositoryImpl) createProject(ctx context.Context, tx *gorm.DB, project *models.Project) error {
	if err := tx.WithContext(ctx).Create(project).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func isPDFFile(fileType string) bool { return fileType == "pdf" || fileType == "application/pdf" }

func (r *projectRepositoryImpl) handleFiles(ctx context.Context, tx *gorm.DB, projectID int, files []*multipart.FileHeader, titles []string) ([]string, error) {
	var uploadedFilePaths []string

	for i, file := range files {
		filePath, fileURL, err := r.uploadFileToMinio(ctx, file)
		if err != nil {
			tx.Rollback()
			return uploadedFilePaths, err
		}
		uploadedFilePaths = append(uploadedFilePaths, filePath)

		resourceType, err := r.getResourceType(ctx, tx, file)
		if err != nil {
			r.deleteUploadedFiles(ctx, uploadedFilePaths)
			tx.Rollback()
			return uploadedFilePaths, err
		}

		pdf := &models.PDF{}
		if isPDFFile(resourceType.MimeType) {
			if pdf, err = utils.ReadPdf(file); err != nil {
				tx.Rollback()
				return uploadedFilePaths, err
			}
		}

		if err := r.createProjectResourceAndResource(ctx, tx, projectID, fileURL, titles[i], pdf, resourceType.ID); err != nil {
			r.deleteUploadedFiles(ctx, uploadedFilePaths)
			tx.Rollback()
			return uploadedFilePaths, err
		}
	}

	return uploadedFilePaths, nil
}

func (r *projectRepositoryImpl) uploadFileToMinio(ctx context.Context, file *multipart.FileHeader) (string, string, error) {
	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	filePath := fmt.Sprintf("uploads/%s", fileName)

	src, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	_, err = r.minioClient.PutObject(ctx, "projects", filePath, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", "", err
	}

	fileURL := fmt.Sprintf("http://localhost:9000/projects/%s", filePath)
	return filePath, fileURL, nil
}

func (r *projectRepositoryImpl) getResourceType(ctx context.Context, tx *gorm.DB, file *multipart.FileHeader) (*models.ResourceType, error) {
	fileType := file.Header.Get("Content-Type")
	if fileType == "" {
		fileType = "unknown"
	}

	resourceType := &models.ResourceType{}
	if err := tx.WithContext(ctx).Where("mime_type = ?", fileType).First(resourceType).Error; err != nil {
		return nil, fmt.Errorf("file type not found: %w", err)
	}

	return resourceType, nil
}

func (r *projectRepositoryImpl) createProjectResourceAndResource(ctx context.Context, tx *gorm.DB, projectID int, fileURL string, title string, pdf *models.PDF, resourceTypeID int) error {
	projectResource := &models.ProjectResource{
		ProjectID: projectID,
	}

	if err := tx.WithContext(ctx).Create(projectResource).Error; err != nil {
		return err
	}

	resource := &models.Resource{
		Title:             title,
		URL:               fileURL,
		ProjectResourceID: &projectResource.ID,
		PDF:               pdf,
		ResourceTypeID:    resourceTypeID,
	}

	if err := tx.WithContext(ctx).Create(resource).Error; err != nil {
		return err
	}

	return nil
}

func (r *projectRepositoryImpl) deleteUploadedFiles(ctx context.Context, filePaths []string) {
	for _, filePath := range filePaths {
		err := r.minioClient.RemoveObject(ctx, "projects", filePath, minio.RemoveObjectOptions{})
		if err != nil {
			log.Printf("Failed to delete file from MinIO: %s, error: %v", filePath, err)
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
	if err := tx.Where("project_id = ?", id).Delete(&models.ProjectEmployee{}).Error; err != nil {
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
