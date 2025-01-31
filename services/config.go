package services

import (
	"context"

	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type ConfigService interface {
	GetConfigByProgramId(programId int) ([]models.Config, error)
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

func (s *configServiceImpl) GetConfigByProgramId(programId int) ([]models.Config, error) {
	configs, err := s.configRepo.GetConfigByProgramId(programId)
	if err != nil {
		return nil, err
	}
	return configs, err
}

func (s *configServiceImpl) FindConfigByNameAndProgramId(ctx context.Context, name string, programId int) (*models.Config, error) {
	config, err := s.configRepo.GetByNameAndProgramId(ctx, name, programId)
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
