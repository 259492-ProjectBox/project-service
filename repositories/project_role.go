package repositories

import (
	"context"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type ProjectRoleRepository interface {
	GetAllByProgramId(ctx context.Context, programId int) ([]models.ProjectRole, error)
}

type projectRoleRepositoryImpl struct {
	db *gorm.DB
}

func NewProjectRoleRepository(db *gorm.DB) ProjectRoleRepository {
	return &projectRoleRepositoryImpl{
		db: db,
	}
}

func (r *projectRoleRepositoryImpl) GetAllByProgramId(ctx context.Context, programId int) ([]models.ProjectRole, error) {
	var projectRoles []models.ProjectRole
	if err := r.db.WithContext(ctx).Where("program_id = ?", programId).Preload("Program").Find(&projectRoles).Error; err != nil {
		return nil, err
	}
	return projectRoles, nil
}
