package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type ConfigService interface {
	GetConfigByProgramId(programId int) ([]models.Config, error)
	GetCurrentAcademicYearAndSemester(ctx context.Context, programId int) (int, int, error)
	DeleteConfig(ctx context.Context, id int) error
	FindConfigByNameAndProgramId(ctx context.Context, name string, programId int) (*models.Config, error)
	UpsertConfig(ctx context.Context, config *models.Config) (*models.Config, error)
}

type configServiceImpl struct {
	configRepo repositories.ConfigRepository
}

func NewConfigService(configRepo repositories.ConfigRepository) ConfigService {
	return &configServiceImpl{
		configRepo: configRepo,
	}
}

func (s *configServiceImpl) GetCurrentAcademicYearAndSemester(ctx context.Context, programId int) (int, int, error) {
	academicYear, err := s.FindConfigByNameAndProgramId(ctx, "academic year", programId)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch academic year: %w", err)
	}

	academicYearInt, err := strconv.Atoi(academicYear.Value)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert academic year: %w", err)
	}

	semester, err := s.FindConfigByNameAndProgramId(ctx, "semester", programId)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch semester: %w", err)
	}

	semesterInt, err := strconv.Atoi(semester.Value)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert semester: %w", err)
	}

	return academicYearInt, semesterInt, nil
}

func (s *configServiceImpl) GetConfigByProgramId(programId int) ([]models.Config, error) {
	configs, err := s.configRepo.GetConfigByProgramId(programId)
	if err != nil {
		return nil, err
	}
	return configs, err
}

func (s *configServiceImpl) FindConfigByNameAndProgramId(ctx context.Context, name string, programId int) (*models.Config, error) {
	config, err := s.configRepo.GetByNameAndProgramId(ctx, name, programId)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (s *configServiceImpl) UpsertConfig(ctx context.Context, config *models.Config) (*models.Config, error) {
	config, err := s.configRepo.Upsert(ctx, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (s *configServiceImpl) DeleteConfig(ctx context.Context, id int) error {
	return s.configRepo.Delete(ctx, id)
}
