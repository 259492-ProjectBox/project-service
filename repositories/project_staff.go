package repositories

import (
	"context"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type ProjectStaffRepository interface {
	repository[models.ProjectStaff]
	CreateProjectStaff(ctx context.Context, tx *gorm.DB, projectStaff *models.ProjectStaff) error
	GetProjectStaffByProjectIdAndStaffId(ctx context.Context, projectId int, staffId int) (*models.ProjectStaff, error)
}

type projectStaffRepositoryImpl struct {
	db *gorm.DB
	*repositoryImpl[models.ProjectStaff]
}

func NewProjectStaffRepository(db *gorm.DB) ProjectStaffRepository {
	return &projectStaffRepositoryImpl{
		db:             db,
		repositoryImpl: newRepository[models.ProjectStaff](db),
	}
}

// CreateProjectStaff creates a new project staff record in the database
func (r *projectStaffRepositoryImpl) CreateProjectStaff(ctx context.Context, tx *gorm.DB, projectStaff *models.ProjectStaff) error {
	if projectStaff == nil {
		return gorm.ErrInvalidData
	}

	db := tx
	if db == nil {
		db = r.db
	}

	if err := db.WithContext(ctx).Create(projectStaff).Error; err != nil {
		return err
	}

	return nil
}

func (r *projectStaffRepositoryImpl) GetProjectStaffByProjectIdAndStaffId(ctx context.Context, projectId int, staffId int) (*models.ProjectStaff, error) {
	var projectStaff models.ProjectStaff
	err := r.db.WithContext(ctx).
		Where("project_id = ? AND staff_id = ?", projectId, staffId).
		Preload("ProjectRole.Program").
		First(&projectStaff).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &projectStaff, nil
}
