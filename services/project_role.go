package services

import (
	"context"

	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type ProjectRoleService interface {
	GetAllProjectRolesByProgramId(ctx context.Context, programId int) ([]models.ProjectRole, error)
}

type projectRoleServiceImpl struct {
	projectRoleRepo repositories.ProjectRoleRepository
}

func NewProjectRoleService(projectRoleRepo repositories.ProjectRoleRepository) ProjectRoleService {
	return &projectRoleServiceImpl{
		projectRoleRepo: projectRoleRepo,
	}
}

func (s *projectRoleServiceImpl) GetAllProjectRolesByProgramId(ctx context.Context, programId int) ([]models.ProjectRole, error) {
	return s.projectRoleRepo.GetAllByProgramId(ctx, programId)
}
