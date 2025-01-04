package services

import (
	"context"
	"errors"

	"github.com/project-box/dtos"
	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type MajorService interface {
	UpdateMajorNameService(ctx context.Context, majorID int, name string) error
	CreateMajorService(ctx context.Context, major *dtos.CreateMajorRequest) error
	GetAllMajorService(ctx context.Context) ([]models.Major, error)
}

type majorServiceImpl struct {
	majorRepo repositories.MajorRepository
}

func NewMajorService(majorRepo repositories.MajorRepository) MajorService {
	return &majorServiceImpl{
		majorRepo: majorRepo,
	}
}

// update major name
func (s *majorServiceImpl) UpdateMajorNameService(ctx context.Context, majorID int, name string) error {
	if majorID <= 0 {
		return errors.New("invalid major ID")
	}
	updateMajor := &models.Major{ID: majorID, MajorName: name}

	return s.majorRepo.UpdateMajorName(ctx, updateMajor)
}

// get all major from repository
func (s *majorServiceImpl) GetAllMajorService(ctx context.Context) ([]models.Major, error) {
	return s.majorRepo.GetAllMajor(ctx)
}

// create major from repository
func (s *majorServiceImpl) CreateMajorService(ctx context.Context, major *dtos.CreateMajorRequest) error {
	modelMajor := &models.Major{
		MajorName: major.MajorName,
	}

	return s.majorRepo.CreateMajor(ctx, modelMajor)
}
