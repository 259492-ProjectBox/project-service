package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type ResourceRepository interface {
	CreateResource(ctx context.Context, resource *models.Resource) (*models.Resource, error)
	FindDetailedResourceByID(ctx context.Context, id string) (*models.DetailedResource, error)
	DeleteResourceByID(ctx context.Context, id string) error
	FindByProjectID(ctx context.Context, projectID string) ([]models.Resource, error)
}

type resourceRepository struct {
	db *gorm.DB
}

func NewResourceRepository(db *gorm.DB) ResourceRepository {
	return &resourceRepository{db: db}
}

func (r *resourceRepository) CreateResource(ctx context.Context, resource *models.Resource) (*models.Resource, error) {
	if err := r.db.WithContext(ctx).Create(resource).Error; err != nil {
		return nil, err
	}
	return resource, nil
}
func (r *resourceRepository) FindDetailedResourceByID(ctx context.Context, id string) (*models.DetailedResource, error) {
	var detailedResource models.DetailedResource
	if err := r.db.WithContext(ctx).
		Table("resources").
		Select("projects.*,resources.*,project_resources.*,asset_resources.*").
		Joins("LEFT JOIN project_resources ON project_resources.id = resources.project_resource_id").
		Joins("LEFT JOIN asset_resources ON asset_resources.id = resources.asset_resource_id").
		Joins("LEFT JOIN projects ON projects.id = project_resources.project_id").
		Where("resources.id = ?", id).
		Scan(&detailedResource).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("resource not found")
		}
		return nil, err
	}
	fmt.Println(detailedResource)
	return &detailedResource, nil
}

func (r *resourceRepository) DeleteResourceByID(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Resource{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("resource not found")
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
