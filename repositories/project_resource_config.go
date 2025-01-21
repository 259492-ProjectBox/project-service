package repositories

import (
	"github.com/project-box/models"
	"gorm.io/gorm"
)

type ProjectResourceConfigRepository interface {
	GetProjectResourceConfigsByProgramId(programID int) ([]models.ProjectResourceConfig, error)
	UpsertResourceProjectConfig(config *models.ProjectResourceConfig) error
}

type projectResourceConfigRepositoryImpl struct {
	db *gorm.DB
}

func NewProjectResourceConfigRepository(db *gorm.DB) ProjectResourceConfigRepository {
	return &projectResourceConfigRepositoryImpl{db: db}
}

func (r *projectResourceConfigRepositoryImpl) GetProjectResourceConfigsByProgramId(programID int) ([]models.ProjectResourceConfig, error) {
	var configs []models.ProjectResourceConfig
	err := r.db.Where("program_id = ?", programID).Preload("Program").Preload("ResourceType").Preload("FileExtension").Find(&configs).Error
	return configs, err
}

func (r *projectResourceConfigRepositoryImpl) UpsertResourceProjectConfig(config *models.ProjectResourceConfig) error {
	return r.db.Save(config).Error
}
