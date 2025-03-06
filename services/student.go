package services

import (
	"context"
	"fmt"

	"github.com/project-box/models"
	"github.com/project-box/repositories"
	"gorm.io/gorm"
)

type StudentService interface {
	CreateStudents(ctx context.Context, students []models.Student) error
	UpsertStudents(ctx context.Context, students []models.Student, programId int) ([]models.Student, error)
	GetStudentByStudentId(ctx context.Context, studentId string) (*models.Student, error)
	GetStudentByProgramIdOnCurrentYearAndSemester(ctx context.Context, programId int) ([]models.Student, error)
	GetStudentByProgramIdOnAcademicYearAndSemester(ctx context.Context, programId, academicYear, semester int) ([]models.Student, error)
	GetStudentByStudentIdAndProgramIdOnCurrentYearAndSemester(ctx context.Context, studentId string, programId int) (*models.Student, error)
}

type studentServiceImpl struct {
	studentRepo   repositories.StudentRepository
	configService ConfigService
	db            *gorm.DB
}

func NewStudentService(configService ConfigService, studentRepo repositories.StudentRepository, db *gorm.DB) StudentService {
	return &studentServiceImpl{
		studentRepo:   studentRepo,
		configService: configService,
		db:            db,
	}
}

func (s *studentServiceImpl) CreateStudents(ctx context.Context, students []models.Student) error {
	if len(students) == 0 {
		return nil
	}
	return s.studentRepo.CreateMany(ctx, students)
}

func (s *studentServiceImpl) GetStudentByProgramIdOnCurrentYearAndSemester(ctx context.Context, programId int) ([]models.Student, error) {
	academicYear, semester, err := s.configService.GetCurrentAcademicYearAndSemester(ctx, programId)
	if err != nil {
		return nil, err
	}
	return s.studentRepo.GetStudentByProgramIdOnAcademicYearAndSemester(ctx, programId, academicYear, semester)
}

func (s *studentServiceImpl) GetStudentByProgramIdOnAcademicYearAndSemester(ctx context.Context, programId int, academicYear, semester int) ([]models.Student, error) {
	return s.studentRepo.GetStudentByProgramIdOnAcademicYearAndSemester(ctx, programId, academicYear, semester)
}

func (s *studentServiceImpl) UpsertStudents(ctx context.Context, students []models.Student, programId int) ([]models.Student, error) {
	if len(students) == 0 {
		return nil, nil
	}

	academicYear, semester, err := s.configService.GetCurrentAcademicYearAndSemester(ctx, programId)
	if err != nil {
		return nil, err
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var upsertedStudents []models.Student
	for _, student := range students {
		upsertedStudent, err := s.upsertStudent(ctx, tx, &student, semester, academicYear)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		upsertedStudents = append(upsertedStudents, *upsertedStudent)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return upsertedStudents, nil
}

func (s *studentServiceImpl) upsertStudent(ctx context.Context, tx *gorm.DB, student *models.Student, semester, academicYear int) (*models.Student, error) {
	existingStudent, err := s.GetStudentByStudentIdAndProgramIdOnCurrentYearAndSemester(ctx, student.StudentID, student.ProgramID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		fmt.Printf("Student %s not found in program %d for year %d semester %d, creating new record...\n",
			student.StudentID, student.ProgramID, semester, academicYear)
		if err := tx.Create(&student).Error; err != nil {
			return nil, err
		}
		return student, nil
	}

	existingStudent = updateStudentFields(existingStudent, student)
	if err := tx.Save(&existingStudent).Error; err != nil {
		return nil, err
	}

	return existingStudent, nil
}

func updateStudentFields(existing, new *models.Student) *models.Student {
	existing.SecLab = new.SecLab
	existing.FirstName = new.FirstName
	existing.LastName = new.LastName
	if new.Email != nil {
		existing.Email = new.Email
	}
	existing.Semester = new.Semester
	existing.AcademicYear = new.AcademicYear
	existing.CourseID = new.CourseID
	existing.ProgramID = new.ProgramID
	return existing
}

func (s *studentServiceImpl) GetStudentByStudentId(ctx context.Context, studentId string) (*models.Student, error) {
	student, err := s.studentRepo.GetStudentByStudentId(ctx, studentId)
	if err != nil {
		return nil, err
	}
	return student, nil
}

func (s *studentServiceImpl) GetStudentByStudentIdAndProgramIdOnCurrentYearAndSemester(ctx context.Context, studentId string, programId int) (*models.Student, error) {
	academicYear, semester, err := s.configService.GetCurrentAcademicYearAndSemester(ctx, programId)
	if err != nil {
		return nil, err
	}
	return s.studentRepo.GetStudentByStudentIdAndProgramIdOnCurrentYearAndSemester(ctx, studentId, programId, academicYear, semester)
}
