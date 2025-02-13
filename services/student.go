package services

import (
	"context"
	"strconv"
	"sync"

	"github.com/project-box/models"
	"github.com/project-box/repositories"
	"gorm.io/gorm"
)

type StudentService interface {
	CreateStudents(ctx context.Context, students []models.Student) error
	UpsertStudents(ctx context.Context, students []models.Student) error
	GetStudentByStudentId(ctx context.Context, studentId string) (*models.Student, error)
	GetStudentByProgramIdOnCurrentYearAndSemester(ctx context.Context, programId int) ([]models.Student, error)
	GetStudentByStudentIdAndProgramIdOnCurrentYearAndSemester(ctx context.Context, studentId string, programId int) (*models.Student, error)
}

type studentServiceImpl struct {
	studentRepo repositories.StudentRepository
	configRepo  repositories.ConfigRepository
	db          *gorm.DB
}

func NewStudentService(configRepo repositories.ConfigRepository, studentRepo repositories.StudentRepository, db *gorm.DB) StudentService {
	return &studentServiceImpl{
		studentRepo: studentRepo,
		configRepo:  configRepo,
		db:          db,
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

func (s *studentServiceImpl) UpsertStudents(ctx context.Context, students []models.Student) error {
	if len(students) == 0 {
		return nil
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	errChan := make(chan error, len(students))
	var wg sync.WaitGroup

	for _, student := range students {
		wg.Add(1)
		go func(student models.Student) {
			defer wg.Done()
			if err := s.upsertStudent(ctx, tx, student); err != nil {
				errChan <- err
			}
		}(student)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		tx.Rollback()
		return <-errChan
	}

	return tx.Commit().Error
}

func (s *studentServiceImpl) upsertStudent(ctx context.Context, tx *gorm.DB, student models.Student) error {
	existingStudent, err := s.GetStudentByStudentIdAndProgramIdOnCurrentYearAndSemester(ctx, student.StudentID, student.ProgramID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if existingStudent != nil {
		existingStudent.StudentID = student.StudentID
		existingStudent.SecLab = student.SecLab
		existingStudent.FirstName = student.FirstName
		existingStudent.LastName = student.LastName
		existingStudent.Semester = student.Semester
		existingStudent.AcademicYear = student.AcademicYear
		existingStudent.CourseID = student.CourseID
		existingStudent.ProgramID = student.ProgramID
		return tx.Save(existingStudent).Error
	}

	return tx.Create(&student).Error
}

func (s *studentServiceImpl) GetStudentByStudentIdAndProgramIdOnCurrentYearAndSemester(ctx context.Context, studentId string, programId int) (*models.Student, error) {
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

	return s.studentRepo.GetStudentByStudentIdAndProgramIdOnCurrentYearAndSemester(ctx, studentId, programId, academicYear, semester)
}
