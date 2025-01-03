package repositories

import (
	"github.com/project-box/models"
	"gorm.io/gorm"
)

type ConfigRepository interface {
	repository[models.Config]
	GetConfigByMajorID(majorID int) ([]models.Config, error)
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

// get all config from major id
func (r *configRepositoryImpl) GetConfigByMajorID(majorID int) ([]models.Config, error) {
	var configs []models.Config
	if err := r.db.Where("major_id = ?", majorID).Find(&configs).Error; err != nil {
		return nil, err
	}
	if len(configs) == 0 {
		return []models.Config{}, nil // Return an empty slice, not nil
	}
	return configs, nil
}
