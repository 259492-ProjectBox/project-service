package services

import (
	"github.com/project-box/repositories"
)

type Config interface {
}

type configServiceImpl struct {
	config repositories.ConfigRepository
}

func NewConfigService() Config {
	return &configServiceImpl{}
}
