package services

import (
	"github.com/project-box/dtos"
	"github.com/project-box/repositories"
)

type ConfigService interface {
	GetConfigByMajorID(majorID int) ([]dtos.ConfigReponse, error)
}

type configServiceImpl struct {
	configRepo repositories.ConfigRepository
}

func NewConfigService(configRepo repositories.ConfigRepository) ConfigService {
	return &configServiceImpl{
		configRepo: configRepo,
	}
}

func (s *configServiceImpl) GetConfigByMajorID(majorID int) ([]dtos.ConfigReponse, error) {
	configs, err := s.configRepo.GetConfigByMajorID(majorID)
	// //convert to dtos.ConfigDto
	// var err error
	var configDtos []dtos.ConfigReponse

	for _, config := range configs {
		configDtos = append(configDtos, dtos.ConfigReponse{
			MajorID:    config.MajorID,
			ConfigName: config.ConfigName,
			Value:      config.Value,
		})

	}
	return configDtos, err

}
