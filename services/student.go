package services

import (
	"context"
	"strconv"

	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type StudentService interface {
	CreateStudents(ctx context.Context, students []models.Student) error
	GetStudentByStudentId(ctx context.Context, studentId string) (*models.Student, error)
	GetStudentByProgramIdOnCurrentYearAndSemester(ctx context.Context, programId int) ([]models.Student, error)
}

type studentServiceImpl struct {
	studentRepo repositories.StudentRepository
	configRepo  repositories.ConfigRepository
}

func NewStudentService(configRepo repositories.ConfigRepository, studentRepo repositories.StudentRepository) StudentService {
	return &studentServiceImpl{
		studentRepo: studentRepo,
		configRepo:  configRepo,
	}
}

func (s *studentServiceImpl) CreateStudents(ctx context.Context, students []models.Student) error {
	if len(students) == 0 {
		return nil
	}

	return s.studentRepo.CreateMany(ctx, students)
}

func (s *studentServiceImpl) GetStudentByStudentId(ctx context.Context, studentId string) (*models.Student, error) {
	return s.studentRepo.GetStudentByStudentId(ctx, studentId)
}

func (s *studentServiceImpl) GetStudentByProgramIdOnCurrentYearAndSemester(ctx context.Context, programId int) ([]models.Student, error) {
	academicYearConfig, err := s.configRepo.GetByNameAndProgramId(ctx, "academic year", programId)
	if err != nil {
		return nil, err
	}

	semesterConfig, err := s.configRepo.GetByNameAndProgramId(ctx, "semester", programId)
	if err != nil {
		return nil, err
	}

	academicYear, err := strconv.Atoi(academicYearConfig.Value)
	if err != nil {
		return nil, err
	}

	semester, err := strconv.Atoi(semesterConfig.Value)
	if err != nil {
		return nil, err
	}

	return s.studentRepo.GetStudentByProgramIdOnCurrentYearAndSemester(ctx, programId, academicYear, semester)
}
