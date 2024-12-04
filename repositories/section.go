package repositories

import (
	"github.com/project-box/models"
	"gorm.io/gorm"
)

type SectionRepository interface {
	repository[models.Section]
}

type sectionRepositoryImpl struct {
	db *gorm.DB
	*repositoryImpl[models.Section]
}

func NewSectionRepository(db *gorm.DB) SectionRepository {
	return &sectionRepositoryImpl{
		db:             db,
		repositoryImpl: newRepository[models.Section](db),
	}
}
