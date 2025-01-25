package services

import (
	"context"

	"github.com/project-box/dtos"
	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type ConfigService interface {
	GetConfigByProgramId(programId int) ([]dtos.ConfigResponse, error)
	FindConfigByNameAndProgramId(ctx context.Context, name string, programId int) (*models.Config, error)
	UpsertConfig(ctx context.Context, config *models.Config) (*models.Config, error)
}

type configServiceImpl struct {
	configRepo repositories.ConfigRepository
}

func NewConfigService(configRepo repositories.ConfigRepository) ConfigService {
	return &configServiceImpl{
		configRepo: configRepo,
	}
}

func (s *configServiceImpl) GetConfigByProgramId(programId int) ([]dtos.ConfigResponse, error) {
	configs, err := s.configRepo.GetConfigByProgramId(programId)
	var configDtos []dtos.ConfigResponse

	for _, config := range configs {
		configDtos = append(configDtos, dtos.ConfigResponse{
			ProgramID:  config.ProgramID,
			ConfigName: config.ConfigName,
			Value:      config.Value,
		})
	}
	return configDtos, err
}

func (s *configServiceImpl) FindConfigByNameAndProgramId(ctx context.Context, name string, programId int) (*models.Config, error) {
	config, err := s.configRepo.FindByNameAndProgramId(ctx, name, programId)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (s *configServiceImpl) UpsertConfig(ctx context.Context, config *models.Config) (*models.Config, error) {
	config, err := s.configRepo.Upsert(ctx, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
