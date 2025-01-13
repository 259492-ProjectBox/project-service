package services

import (
	"context"
	"errors"

	"github.com/project-box/dtos"
	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type ProgramService interface {
	UpdateProgramName(ctx context.Context, programId int, name string) error
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

// get all program from repository
func (s *programServiceImpl) GetPrograms(ctx context.Context) ([]models.Program, error) {
	return s.programRepo.GetPrograms(ctx)
}

// update program name
func (s *programServiceImpl) UpdateProgramName(ctx context.Context, programId int, programName string) error {
	if programId <= 0 {
		return errors.New("invalid program ID")
	}
	updatedProgram := &models.Program{ID: programId, ProgramName: programName}

	return s.programRepo.UpdateProgramName(ctx, updatedProgram)
}

// create program from repository
func (s *programServiceImpl) CreateProgram(ctx context.Context, programBody *dtos.CreateProgramRequest) error {
	program := &models.Program{
		ProgramName: programBody.ProgramName,
	}

	return s.programRepo.CreateProgram(ctx, program)
}
