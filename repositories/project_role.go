package repositories

import (
	"github.com/project-box/models"
	"gorm.io/gorm"
)

type ProjectRoleRepository interface {
	repository[models.ProjectRole]
}

type projectRoleRepositoryImpl struct {
	db *gorm.DB
	*repositoryImpl[models.ProjectRole]
}

func NewProjectRoleRepository(db *gorm.DB) ProjectRoleRepository {
	return &projectRoleRepositoryImpl{
		db:             db,
		repositoryImpl: newRepository[models.ProjectRole](db),
	}
}
