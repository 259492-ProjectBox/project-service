package services

import "github.com/project-box/repositories"

type ProjectNumberCounterService interface {
}

type projectNumberCounterServiceImpl struct {
	projectNumberCounterRepository repositories.ProjectNumberCounterRepository
}

func NewProjectNumberCounterService(projectNumberCounterRepository repositories.ProjectNumberCounterRepository) ProjectNumberCounterService {
	return &projectNumberCounterServiceImpl{
		projectNumberCounterRepository: projectNumberCounterRepository,
	}
}
