package repositories

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/project-box/models"
	"gorm.io/gorm"
)

type ResourceRepository interface {
	CreateAssetResource(ctx context.Context, file *multipart.FileHeader, assetResource *models.AssetResource) (*models.AssetResource, error)
	DeleteAssetResourceByID(ctx context.Context, id string) error
	FindAssetResourcesByProgramID(ctx context.Context, id string) ([]models.AssetResource, error)

	CreateProjectResource(ctx context.Context, tx *gorm.DB, projectResource *models.ProjectResource) error
	FindDetailedResourceByID(ctx context.Context, id string) (*models.DetailedResource, error)
	DeleteProjectResourceByID(ctx context.Context, id string, filePath *string) error
}

type resourceRepository struct {
	db                *gorm.DB
	resourceTypeRepo  ResourceTypeRepository
	fileExtensionRepo FileExtensionRepository
	uploadRepo        UploadRepository
}

func NewResourceRepository(db *gorm.DB, resourceTypeRepo ResourceTypeRepository, fileExtensionRepo FileExtensionRepository, uploadRepo UploadRepository) ResourceRepository {
	return &resourceRepository{
		db:                db,
		resourceTypeRepo:  resourceTypeRepo,
		uploadRepo:        uploadRepo,
		fileExtensionRepo: fileExtensionRepo,
	}
}

func (r *resourceRepository) generateFilePath(fileName, programName, bucket string) (string, string) {
	uniqueFileName := r.generateUniqueFileName(fileName)
	objectName := r.buildObjectName(programName, uniqueFileName)
	filePath := r.buildFilePath(bucket, objectName)
	return objectName, filePath
}

func (r *resourceRepository) generateUniqueFileName(fileName string) string {
	return fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileName)
}

func (r *resourceRepository) buildObjectName(programName, uniqueFileName string) string {
	return fmt.Sprintf("%s/%s", programName, uniqueFileName)
}

func (r *resourceRepository) buildFilePath(bucket, objectName string) string {
	return fmt.Sprintf("%s/%s", bucket, objectName)
}

func (r *resourceRepository) uploadFileToMinio(ctx context.Context, objectName string, file *multipart.FileHeader) error {
	err := r.uploadRepo.UploadFile(ctx, os.Getenv("MINIO_ASSET_BUCKET"), objectName, file, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to MinIO: %w", err)
	}
	return nil
}

func (r *resourceRepository) CreateAssetResource(ctx context.Context, file *multipart.FileHeader, assetResource *models.AssetResource) (*models.AssetResource, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := tx.Create(assetResource).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create asset resource: %w", err)
	}

	if err := tx.Preload("Program").First(assetResource).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to preload program: %w", err)
	}

	programNameTH := assetResource.Program.ProgramNameTH
	objectName, _ := r.generateFilePath(file.Filename, programNameTH, os.Getenv("MINIO_ASSET_BUCKET"))

	if err := r.uploadFileToMinio(ctx, objectName, file); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("file upload failed: %w", err)
	}

	if err := tx.Create(assetResource).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	return assetResource, nil
}

func (r *resourceRepository) FindAssetResourcesByProgramID(ctx context.Context, program_id string) ([]models.AssetResource, error) {
	var assetResources []models.AssetResource
	if err := r.db.WithContext(ctx).Where("program_id = ?", program_id).Find(&assetResources).Error; err != nil {
		return nil, err
	}
	return assetResources, nil
}

func (r *resourceRepository) FindDetailedResourceByID(ctx context.Context, id string) (*models.DetailedResource, error) {
	var detailedResource models.DetailedResource

	query := r.db.WithContext(ctx).
		Table("project_resources").
		Select(`
        projects.id AS project_id,
        project_resources.id AS project_resource_id,
        projects.*,
        project_resources.*,
    `).
		Joins("LEFT JOIN projects ON project_resources.project_id IS NOT NULL AND projects.id = project_resources.project_id")

	if err := query.Where("project_resources.id = ?", id).
		Scan(&detailedResource).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("resource not found")
		}
		return nil, err
	}

	return &detailedResource, nil
}

func (r *resourceRepository) DeleteProjectResourceByID(ctx context.Context, id string, filePath *string) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if result := tx.Delete(&models.ProjectResource{}, id); result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if filePath != nil {
		if err := r.uploadRepo.DeleteUploadedFile(ctx, os.Getenv("MINIO_PROJECT_BUCKET"), *filePath, minio.RemoveObjectOptions{}); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *resourceRepository) DeleteAssetResourceByID(ctx context.Context, id string) error {
	if result := r.db.WithContext(ctx).Delete(&models.AssetResource{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *resourceRepository) CreateProjectResource(ctx context.Context, tx *gorm.DB, projectResource *models.ProjectResource) error {
	if err := tx.WithContext(ctx).Create(projectResource).Error; err != nil {
		return err
	}
	return nil
}
