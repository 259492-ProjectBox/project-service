package services

import (
	"github.com/project-box/dtos"
	"github.com/project-box/repositories"
)

type ConfigService interface {
	GetConfigByProgramId(programId int) ([]dtos.ConfigReponse, error)
}

type configServiceImpl struct {
	configRepo repositories.ConfigRepository
}

func NewConfigService(configRepo repositories.ConfigRepository) ConfigService {
	return &configServiceImpl{
		configRepo: configRepo,
	}
}

func (s *configServiceImpl) GetConfigByProgramId(programId int) ([]dtos.ConfigReponse, error) {
	configs, err := s.configRepo.GetConfigByProgramId(programId)
	var configDtos []dtos.ConfigReponse

	for _, config := range configs {
		configDtos = append(configDtos, dtos.ConfigReponse{
			ProgramID:  config.ProgramID,
			ConfigName: config.ConfigName,
			Value:      config.Value,
		})

	}
	return configDtos, err

}
