package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/project-box/dtos"
	"github.com/project-box/models"
	"github.com/project-box/repositories"
	"github.com/robfig/cron/v3"
)

type ConfigService interface {
	GetConfigByProgramId(programId int) ([]models.Config, error)
	GetCurrentAcademicYearAndSemester(ctx context.Context, programId int) (int, int, error)
	GetConfigByNameAndProgramId(ctx context.Context, name string, programId int) (*models.Config, error)
	GetAllAcademicYear(ctx context.Context) ([]dtos.AcademicYearResponse, error)
	UpsertConfig(ctx context.Context, config *models.Config) (*models.Config, error)
	DeleteConfig(ctx context.Context, id int) error
	IncreaseAcademicYear(ctx context.Context) error
	StartCronJob()
}

type configServiceImpl struct {
	configRepo repositories.ConfigRepository
	cron       *cron.Cron
}

func NewConfigService(configRepo repositories.ConfigRepository) ConfigService {
	service := &configServiceImpl{
		configRepo: configRepo,
		cron:       cron.New(),
	}
	service.StartCronJob()
	return service
}

func (s *configServiceImpl) GetAllAcademicYear(ctx context.Context) ([]dtos.AcademicYearResponse, error) {
	lowestAcademicYearConfig, err := s.configRepo.GetByName(ctx, "lowest_academic_year")
	if err != nil {
		return nil, err
	}
	highestAcademicYearConfig, err := s.configRepo.GetByName(ctx, "highest_academic_year")
	if err != nil {
		return nil, err
	}

	lowestAcademicYearConfigInt, err := strconv.Atoi(lowestAcademicYearConfig.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert lowest academic year: %w", err)
	}
	highestAcademicYearConfigInt, err := strconv.Atoi(highestAcademicYearConfig.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert highest academic year: %w", err)
	}

	var academicYearResponses []dtos.AcademicYearResponse
	for year := lowestAcademicYearConfigInt; year <= highestAcademicYearConfigInt; year++ {
		academicYearResponses = append(academicYearResponses, dtos.AcademicYearResponse{
			Year_AD: year - 543,
			Year_BE: year,
		})
	}

	return academicYearResponses, nil
}

func (s *configServiceImpl) IncreaseAcademicYear(ctx context.Context) error {
	highestAcademicYearConfig, err := s.configRepo.GetByName(ctx, "highest_academic_year")
	if err != nil {
		return err
	}
	highestAcademicYearConfigInt, err := strconv.Atoi(highestAcademicYearConfig.Value)
	if err != nil {
		return err
	}
	highestAcademicYearConfigInt++
	highestAcademicYearConfig.Value = strconv.Itoa(highestAcademicYearConfigInt)
	_, err = s.configRepo.Upsert(ctx, highestAcademicYearConfig)
	if err != nil {
		return err
	}
	return nil
}

func (s *configServiceImpl) GetCurrentAcademicYearAndSemester(ctx context.Context, programId int) (int, int, error) {
	academicYear, err := s.GetConfigByNameAndProgramId(ctx, "academic year", programId)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch academic year: %w", err)
	}

	academicYearInt, err := strconv.Atoi(academicYear.Value)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert academic year: %w", err)
	}

	semester, err := s.GetConfigByNameAndProgramId(ctx, "semester", programId)
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

func (s *configServiceImpl) GetConfigByNameAndProgramId(ctx context.Context, name string, programId int) (*models.Config, error) {
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

func (s *configServiceImpl) StartCronJob() {
	s.cron.AddFunc("@yearly", func() {
		ctx := context.Background()
		if err := s.IncreaseAcademicYear(ctx); err != nil {
			fmt.Printf("Failed to increase academic year: %v\n", err)
		} else {
			fmt.Println("Successfully increased academic year")
		}
	})
	s.cron.Start()
}
