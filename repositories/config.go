package repositories

import (
	"github.com/project-box/models"
	"gorm.io/gorm"
)

type ConfigRepository interface {
	repository[models.Config]
}

type configRepositoryImpl struct {
	db *gorm.DB
	*repositoryImpl[models.Config]
}

func NewConfigRepository(db *gorm.DB) ConfigRepository {
	return &configRepositoryImpl{
		db:             db,
		repositoryImpl: newRepository[models.Config](db),
	}
}
