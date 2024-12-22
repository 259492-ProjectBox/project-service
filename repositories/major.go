package repositories

import (
	"context"
	"errors"

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
	var major models.Major

	// Use the correct field name for filtering (ID column corresponds to "id")
	if err := r.db.WithContext(ctx).
		Where("id = ?", majorID).
		First(&major).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("major ID does not exist")
		}
		return nil, err
	}

	return &major, nil
}
