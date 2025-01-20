package services

import (
	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type ProjectResourceConfigService interface {
	GetProjectResourceConfigsByProgramId(programID int) ([]models.ProjectResourceConfig, error)
	UpsertResourceProjectConfig(config *models.ProjectResourceConfig) error
}

type projectResourceConfigServiceImpl struct {
	projectResourceConfigRepo repositories.ProjectResourceConfigRepository
}

func NewProjectResourceConfigService(repo repositories.ProjectResourceConfigRepository) ProjectResourceConfigService {
	return &projectResourceConfigServiceImpl{projectResourceConfigRepo: repo}
}

func (s *projectResourceConfigServiceImpl) GetProjectResourceConfigsByProgramId(programID int) ([]models.ProjectResourceConfig, error) {
	return s.projectResourceConfigRepo.GetProjectResourceConfigsByProgramId(programID)
}

func (s *projectResourceConfigServiceImpl) UpsertResourceProjectConfig(config *models.ProjectResourceConfig) error {
	return s.projectResourceConfigRepo.UpsertResourceProjectConfig(config)
}
