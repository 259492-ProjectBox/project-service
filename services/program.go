package services

import (
	"context"
	"errors"

	"github.com/project-box/dtos"
	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type ProgramService interface {
	GetProgramByID(ctx context.Context, id int) (*models.Program, error)
	UpdateProgram(ctx context.Context, program *models.Program) (*models.Program, error)
	CreateProgram(ctx context.Context, program *dtos.CreateProgramRequest) error
	GetPrograms(ctx context.Context) ([]models.Program, error)
}

type programServiceImpl struct {
	programRepo repositories.ProgramRepository
}

func NewProgramService(programRepo repositories.ProgramRepository) ProgramService {
	return &programServiceImpl{
		programRepo: programRepo,
	}
}

func (s *programServiceImpl) GetProgramByID(ctx context.Context, id int) (*models.Program, error) {
	return s.programRepo.Get(ctx, id)
}

func (s *programServiceImpl) GetPrograms(ctx context.Context) ([]models.Program, error) {
	return s.programRepo.GetPrograms(ctx)
}

func (s *programServiceImpl) UpdateProgram(ctx context.Context, program *models.Program) (*models.Program, error) {
	if program.ID <= 0 {
		return nil, errors.New("invalid program ID")
	}
	return s.programRepo.Update(ctx, program.ID, program)
}

func (s *programServiceImpl) CreateProgram(ctx context.Context, programBody *dtos.CreateProgramRequest) error {
	program := &models.Program{
		ProgramNameTH: programBody.ProgramNameTH,
		ProgramNameEN: programBody.ProgramNameEN,
	}

	return s.programRepo.CreateProgram(ctx, program)
}
