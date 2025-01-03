package services

import (
	"github.com/project-box/dtos"
	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type ProjectConfigService interface {
	GetProjectConfigByMajorIDService(majorID int) ([]dtos.ProjectConfigResponse, error)
	UpdateProjectConfigService(configs []dtos.ProjectConfigUpsertRequest) error
}

type projectconfigServiceImpl struct {
	projectconfigRepo repositories.ProjectConfigRepository
}

func NewProjectConfigService(projectconfigRepo repositories.ProjectConfigRepository) ProjectConfigService {
	return &projectconfigServiceImpl{
		projectconfigRepo: projectconfigRepo,
	}
}

func (s *projectconfigServiceImpl) GetProjectConfigByMajorIDService(majorID int) ([]dtos.ProjectConfigResponse, error) {
	configs, err := s.projectconfigRepo.GetProjectConfigByMajorID(majorID)

	var configDtos []dtos.ProjectConfigResponse

	for _, config := range configs {
		configDtos = append(configDtos, dtos.ProjectConfigResponse{
			ID:       config.ID,
			Title:    config.Title,
			MajorID:  config.MajorID,
			IsActive: config.IsActive,
		})

	}
	return configDtos, err

}

func (s *projectconfigServiceImpl) UpdateProjectConfigService(configs []dtos.ProjectConfigUpsertRequest) error {
	var updateProjectConfigs []models.ProjectConfig
	var insertProjectConfigs []models.ProjectConfig

	// Separate configs into update and insert arrays based on ID
	for _, config := range configs {
		projectConfig := models.ProjectConfig{
			ID:       config.ID,
			Title:    config.Title,
			MajorID:  config.MajorID,
			IsActive: config.IsActive,
		}

		if config.ID > 0 {
			updateProjectConfigs = append(updateProjectConfigs, projectConfig)
		} else {
			insertProjectConfigs = append(insertProjectConfigs, projectConfig)
		}
	}

	// Handle updates if there are any
	if len(updateProjectConfigs) > 0 {
		if err := s.projectconfigRepo.UpdateProjectConfig(updateProjectConfigs); err != nil {
			return err
		}
	}

	// Handle inserts if there are any
	if len(insertProjectConfigs) > 0 {
		if err := s.projectconfigRepo.InsertProjectConfig(insertProjectConfigs); err != nil {
			return err
		}
	}

	return nil
}
