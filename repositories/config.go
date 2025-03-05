package repositories

import (
	"context"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type ConfigRepository interface {
	repository[models.Config]
	GetConfigByProgramId(programId int) ([]models.Config, error)
	GetByNameAndProgramId(ctx context.Context, name string, programId int) (*models.Config, error)
	GetByName(ctx context.Context, name string) (*models.Config, error)
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

func (r *configRepositoryImpl) GetConfigByProgramId(programId int) ([]models.Config, error) {
	var configs []models.Config
	if err := r.db.Where("program_id = ?", programId).Preload("Program").Find(&configs).Error; err != nil {
		return nil, err
	}
	if len(configs) == 0 {
		return []models.Config{}, nil
	}
	return configs, nil
}

func (r *configRepositoryImpl) GetByNameAndProgramId(ctx context.Context, name string, programId int) (*models.Config, error) {
	var config models.Config
	err := r.db.WithContext(ctx).Where("config_name = ? AND program_id = ?", name, programId).Preload("Program").First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *configRepositoryImpl) GetByName(ctx context.Context, name string) (*models.Config, error) {
	var config models.Config
	err := r.db.WithContext(ctx).Where("config_name = ? ", name).Preload("Program").First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}
