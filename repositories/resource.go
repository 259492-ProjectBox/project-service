package repositories

import (
	"context"
	"errors"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/project-box/models"
	"gorm.io/gorm"
)

type ResourceRepository interface {
	CreateResource(ctx context.Context, resource *models.Resource) (*models.Resource, error)
	FindDetailedResourceByID(ctx context.Context, id string) (*models.DetailedResource, error)
	DeleteProjectResourceByID(ctx context.Context, id string, filePath string) error
	FindByProjectID(ctx context.Context, projectID string) ([]models.Resource, error)
}

type resourceRepository struct {
	db          *gorm.DB
	minioClient *minio.Client
}

func NewResourceRepository(db *gorm.DB, minioClient *minio.Client) ResourceRepository {
	return &resourceRepository{
		db:          db,
		minioClient: minioClient,
	}
}

func (r *resourceRepository) CreateResource(ctx context.Context, resource *models.Resource) (*models.Resource, error) {
	if err := r.db.WithContext(ctx).Create(resource).Error; err != nil {
		return nil, err
	}
	return resource, nil
}
func (r *resourceRepository) FindDetailedResourceByID(ctx context.Context, id string) (*models.DetailedResource, error) {
	var detailedResource models.DetailedResource

	// Start building the query
	query := r.db.WithContext(ctx).
		Table("resources").
		Select(`
			projects.id AS project_id,
			project_resources.id AS project_resource_id,
			asset_resources.id AS asset_resource_id,
			projects.*,
			resources.*,
			project_resources.*,
			asset_resources.*
		`).
		Joins("LEFT JOIN project_resources ON resources.project_resource_id IS NOT NULL AND project_resources.id = resources.project_resource_id").
		Joins("LEFT JOIN asset_resources ON resources.asset_resource_id IS NOT NULL AND asset_resources.id = resources.asset_resource_id").
		Joins("LEFT JOIN projects ON project_resources.project_id IS NOT NULL AND projects.id = project_resources.project_id")

	if err := query.Where("resources.id = ?", id).
		Scan(&detailedResource).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("resource not found")
		}
		return nil, err
	}

	return &detailedResource, nil
}

func (r *resourceRepository) DeleteProjectResourceByID(ctx context.Context, id string, filePath string) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if result := tx.Delete(&models.ProjectResource{}, id); result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if err := r.minioClient.RemoveObject(ctx, os.Getenv("MINIO_PROJECT_BUCKET"), filePath, minio.RemoveObjectOptions{}); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *resourceRepository) FindByProjectID(ctx context.Context, projectID string) ([]models.Resource, error) {
	var resources []models.Resource
	if err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Find(&resources).Error; err != nil {
		return nil, err
	}
	return resources, nil
}
