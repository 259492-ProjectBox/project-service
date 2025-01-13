package services

import (
	"github.com/project-box/repositories"
)

type Program interface {
}

type ProgramServiceImpl struct {
	programRepo repositories.ProgramRepository
}

func NewProgramService(programRepo repositories.ProgramRepository) Program {
	return &ProgramServiceImpl{
		programRepo: programRepo,
	}
}
