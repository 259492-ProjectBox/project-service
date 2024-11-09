package repositories

import (
	"context"
	"errors"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type ResourceRepository interface {
	CreateResource(ctx context.Context, resource *models.Resource) (*models.Resource, error)
	FindResourceByID(ctx context.Context, id string) (*models.Resource, error)
	DeleteResource(ctx context.Context, id string) error
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

func (r *resourceRepository) FindResourceByID(ctx context.Context, id string) (*models.Resource, error) {
	var resource models.Resource
	if err := r.db.WithContext(ctx).First(&resource, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("resource not found")
		}
		return nil, err
	}
	return &resource, nil
}

func (r *resourceRepository) DeleteResource(ctx context.Context, id string) error {
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
