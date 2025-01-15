package repositories

import (
	"context"
	"errors"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type ResourceTypeRepository interface {
	GetResourceTypeByName(ctx context.Context, tx *gorm.DB, typeName string) (*models.ResourceType, error)
}

type resourceTypeRepository struct {
	db *gorm.DB
}

func NewResourceTypeRepository(db *gorm.DB) ResourceTypeRepository {
	return &resourceTypeRepository{
		db: db,
	}
}

func (r *resourceTypeRepository) GetResourceTypeByName(ctx context.Context, tx *gorm.DB, typeName string) (*models.ResourceType, error) {
	var resourceType models.ResourceType

	if tx == nil {
		tx = r.db
	}

	err := tx.WithContext(ctx).Where("type_name = ?", typeName).First(&resourceType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("resource type not found")
		}
		return nil, err
	}

	return &resourceType, nil
}
