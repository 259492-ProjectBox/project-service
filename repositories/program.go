package repositories

import (
	"context"
	"errors"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type ProgramRepository interface {
	repository[models.Program]
	GetProgramById(ctx context.Context, programId int) (*models.Program, error)
	GetPrograms(ctx context.Context) ([]models.Program, error)
	CreateProgram(ctx context.Context, program *models.Program) error
}

type programRepositoryImpl struct {
	db *gorm.DB
	*repositoryImpl[models.Program]
}

func NewProgramRepository(db *gorm.DB) ProgramRepository {
	return &programRepositoryImpl{
		db:             db,
		repositoryImpl: newRepository[models.Program](db),
	}
}

// get all program
func (r *programRepositoryImpl) GetPrograms(ctx context.Context) ([]models.Program, error) {
	var programs []models.Program

	if err := r.db.WithContext(ctx).Find(&programs).Error; err != nil {
		return nil, err

	}
	return programs, nil
}

// create program
func (r *programRepositoryImpl) CreateProgram(ctx context.Context, program *models.Program) error {
	return r.db.WithContext(ctx).Create(program).Error
}

func (r *programRepositoryImpl) GetProgramById(ctx context.Context, programId int) (*models.Program, error) {
	var program models.Program

	if err := r.db.WithContext(ctx).
		Where("id = ?", programId).
		First(&program).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("program ID does not exist")
		}
		return nil, err
	}

	return &program, nil
}
