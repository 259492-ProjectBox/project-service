package repositories

import (
	"github.com/project-box/models"
	"gorm.io/gorm"
)

type ProjectConfigRepository interface {
	repository[models.ProjectConfig]
	GetProjectConfigByMajorID(majorID int) ([]models.ProjectConfig, error)
	UpdateProjectConfig(configs []models.ProjectConfig) error
	InsertProjectConfig(configs []models.ProjectConfig) error
}

type projectconfigRepositoryImpl struct {
	db *gorm.DB
	*repositoryImpl[models.ProjectConfig]
}

func NewProjectConfigRepository(db *gorm.DB) ProjectConfigRepository {
	return &projectconfigRepositoryImpl{
		db:             db,
		repositoryImpl: newRepository[models.ProjectConfig](db),
	}
}

func (r *projectconfigRepositoryImpl) GetProjectConfigByMajorID(majorID int) ([]models.ProjectConfig, error) {
	var configs []models.ProjectConfig

	if err := r.db.
		Where("major_id = ?", majorID).
		Find(&configs).
		Error; err != nil {
		return nil, err
	}

	return configs, nil
}

// update project_config with array of project_config
func (r *projectconfigRepositoryImpl) UpdateProjectConfig(configs []models.ProjectConfig) error {
	return r.db.Save(configs).Error
}

// insert project_config with array of project_config
func (r *projectconfigRepositoryImpl) InsertProjectConfig(configs []models.ProjectConfig) error {
	return r.db.Create(configs).Error
}
