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
	UpdateMajorName(ctx context.Context, major *models.Major) error
	GetAllMajor(ctx context.Context) ([]models.Major, error)
	CreateMajor(ctx context.Context, major *models.Major) error
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

// get all major
func (r *majorRepositoryImpl) GetAllMajor(ctx context.Context) ([]models.Major, error) {
	var majors []models.Major

	if err := r.db.WithContext(ctx).Find(&majors).Error; err != nil {
		return nil, err

	}
	return majors, nil
}

// create major
func (r *majorRepositoryImpl) CreateMajor(ctx context.Context, major *models.Major) error {
	return r.db.WithContext(ctx).Create(major).Error
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

// update major name with major model
func (r *majorRepositoryImpl) UpdateMajorName(ctx context.Context, major *models.Major) error {
	return r.db.WithContext(ctx).
		Model(&models.Major{}).
		Where("id = ?", major.ID).
		Update("major_name", major.MajorName).
		Error
}
