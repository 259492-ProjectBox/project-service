package repositories

import (
	"github.com/project-box/models"
	"gorm.io/gorm"
)

type MajorRepository interface {
	repository[models.Major]
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
