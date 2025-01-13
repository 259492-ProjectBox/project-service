package repositories

import (
	"github.com/project-box/models"
	"gorm.io/gorm"
)

type ConfigRepository interface {
	repository[models.Config]
	GetConfigByProgramId(programId int) ([]models.Config, error)
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

// get all config from program id
func (r *configRepositoryImpl) GetConfigByProgramId(programId int) ([]models.Config, error) {
	var configs []models.Config
	if err := r.db.Where("program_id = ?", programId).Find(&configs).Error; err != nil {
		return nil, err
	}
	if len(configs) == 0 {
		return []models.Config{}, nil // Return an empty slice, not nil
	}
	return configs, nil
}
