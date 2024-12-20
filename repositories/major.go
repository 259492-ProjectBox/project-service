package repositories

import (
	"context"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type MajorRepository interface {
	repository[models.Major]
	GetByMajorID(ctx context.Context, majorID int) (*models.Major, error)
}

type majorRepositoryImpl struct {
	db *gorm.DB
	*repositoryImpl[models.Major]
}

func NewMajorRepository(db *gorm.DB) MajorRepository {
	return &majorRepositoryImpl{
		db:             db,
		repositoryImpl: newRepository[models.Major](db),
	}
}

func (r *majorRepositoryImpl) GetByMajorID(ctx context.Context, majorID int) (*models.Major, error) {
	filters := map[string]interface{}{"major_id": majorID}
	var major models.Major
	if err := r.db.WithContext(ctx).Where(filters).First(&major).Error; err != nil {
		return nil, err
	}
	return &major, nil
}
