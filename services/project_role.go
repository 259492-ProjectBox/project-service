package services

import (
	"context"
	"errors"

	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type ProjectRoleService interface {
	GetAllProjectRolesByProgramId(ctx context.Context, programId int) ([]models.ProjectRole, error)
	GetProjectRoleByRoleName(ctx context.Context, roleName string) (*models.ProjectRole, error)
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

func (s *projectRoleServiceImpl) GetProjectRoleByRoleName(ctx context.Context, roleName string) (*models.ProjectRole, error) {
	projectRoles, err := s.projectRoleRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, role := range projectRoles {
		if role.RoleNameTH == roleName || role.RoleNameEN == roleName {
			return &role, nil
		}
	}

	return nil, errors.New("project role not found")
}
