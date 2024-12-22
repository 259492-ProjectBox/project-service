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
	CreateProjectWithFiles(ctx context.Context, project *models.Project, major *models.Major, files []*multipart.FileHeader, titles []string) (*models.Project, error)
	UpdateProjectWithFiles(ctx context.Context, project *models.Project, major *models.Major, files []*multipart.FileHeader, titles []string) (*models.Project, error)
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

func (r *projectRepositoryImpl) CreateProjectWithFiles(ctx context.Context, project *models.Project, major *models.Major, files []*multipart.FileHeader, titles []string) (*models.Project, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := r.createProject(ctx, tx, project); err != nil {
		return nil, err
	}

	uploadedFilePaths, err := r.handleFilesUpload(ctx, tx, project, major, files, titles)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		r.deleteUploadedFiles(ctx, uploadedFilePaths)
		return nil, err
	}

	return r.GetProjectWithPDFByID(ctx, project.ID)
}

func (r *projectRepositoryImpl) UpdateProjectWithFiles(ctx context.Context, project *models.Project, major *models.Major, files []*multipart.FileHeader, titles []string) (*models.Project, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	if err := r.deleteProjectAssociations(ctx, tx, project.ID); err != nil {
		return nil, err
	}
	if err := r.updateProject(ctx, tx, project); err != nil {
		return nil, err
	}
	uploadedFilePaths, err := r.handleFilesUpload(ctx, tx, project, major, files, titles)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		r.deleteUploadedFiles(ctx, uploadedFilePaths)
		return nil, err
	}

	return r.GetProjectWithPDFByID(ctx, project.ID)
}

func (r *projectRepositoryImpl) deleteProjectAssociations(ctx context.Context, tx *gorm.DB, projectID int) error {
	if err := r.deleteProjectStudents(ctx, tx, projectID); err != nil {
		return fmt.Errorf("failed to delete project students: %w", err)
	}

	if err := r.deleteProjectEmployees(ctx, tx, projectID); err != nil {
		return fmt.Errorf("failed to delete project employees: %w", err)
	}

	return nil
}

func (r *projectRepositoryImpl) deleteProjectStudents(ctx context.Context, tx *gorm.DB, projectID int) error {
	if err := tx.WithContext(ctx).Where("project_id = ?", projectID).Delete(&models.ProjectStudent{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *projectRepositoryImpl) deleteProjectEmployees(ctx context.Context, tx *gorm.DB, projectID int) error {
	if err := tx.WithContext(ctx).Where("project_id = ?", projectID).Delete(&models.ProjectEmployee{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *projectRepositoryImpl) createProject(ctx context.Context, tx *gorm.DB, project *models.Project) error {
	if err := tx.WithContext(ctx).Create(project).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *projectRepositoryImpl) updateProject(ctx context.Context, tx *gorm.DB, project *models.Project) error {
	if err := tx.WithContext(ctx).Save(project).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func isPDFFile(fileType string) bool { return fileType == "pdf" || fileType == "application/pdf" }

func (r *projectRepositoryImpl) handleFilesUpload(ctx context.Context, tx *gorm.DB, project *models.Project, major *models.Major, files []*multipart.FileHeader, titles []string) ([]string, error) {
	var uploadedObjectNames []string
	for i, file := range files {
		title := titles[i]
		uniqueFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		filePath := fmt.Sprintf("%s/%s/%s/%s", os.Getenv("MINIO_PROJECT_BUCKET"), major.MajorName, title, uniqueFileName)
		objectName := fmt.Sprintf("%s/%s/%s", major.MajorName, title, uniqueFileName)
		err := r.uploadFileToMinio(ctx, objectName, file)
		if err != nil {
			tx.Rollback()
			return uploadedObjectNames, err
		}
		uploadedObjectNames = append(uploadedObjectNames, objectName)

		resourceType, err := r.getResourceType(ctx, tx, file)
		if err != nil {
			r.deleteUploadedFiles(ctx, uploadedObjectNames)
			tx.Rollback()
			return uploadedObjectNames, err
		}

		var pdf *models.PDF
		if isPDFFile(resourceType.MimeType) {
			if pdf, err = utils.ReadPdf(file); err != nil {
				tx.Rollback()
				return uploadedObjectNames, err
			}
		}

		if err := r.createProjectResourceAndResource(ctx, tx, project.ID, title, filePath, file.Filename, pdf, resourceType.ID); err != nil {
			r.deleteUploadedFiles(ctx, uploadedObjectNames)
			tx.Rollback()
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

func (r *projectRepositoryImpl) createProjectResourceAndResource(ctx context.Context, tx *gorm.DB, projectID int, title string, filePath string, fileName string, pdf *models.PDF, resourceTypeID int) error {
	projectResource := &models.ProjectResource{
		ProjectID: projectID,
	}

	if err := tx.WithContext(ctx).Create(projectResource).Error; err != nil {
		return err
	}

	resource := &models.Resource{
		Title:             title,
		ProjectResourceID: &projectResource.ID,
		ResourceName:      fileName,
		Path:              filePath,
		PDF:               pdf,
		ResourceTypeID:    resourceTypeID,
	}

	if err := tx.WithContext(ctx).Create(resource).Error; err != nil {
		return err
	}

	return nil
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
