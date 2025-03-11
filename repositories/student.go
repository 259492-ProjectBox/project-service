package repositories

import (
	"context"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type StudentRepository interface {
	repository[models.Student]
	GetStudentByStudentId(ctx context.Context, studentId string) (*models.Student, error)
	GetStudentByStudentIdAndProgramIdOnCurrentYearAndSemester(ctx context.Context, studentId string, programId, academicYear, semester int) (*models.Student, error)
	GetStudentByProgramIdOnAcademicYearAndSemester(ctx context.Context, programId, academicYear, semester int) ([]models.Student, error)
	CheckStudentDuplicateProjectOnCurrentYearAndSemester(ctx context.Context, studentId string, academicYear, semester int) (bool, error)
}

type studentRepositoryImpl struct {
	db         *gorm.DB
	configRepo ConfigRepository
	*repositoryImpl[models.Student]
}

func NewStudentRepository(db *gorm.DB, configRepo ConfigRepository) StudentRepository {
	return &studentRepositoryImpl{
		db:             db,
		configRepo:     configRepo,
		repositoryImpl: newRepository[models.Student](db),
	}
}

func (r *studentRepositoryImpl) CheckStudentDuplicateProjectOnCurrentYearAndSemester(ctx context.Context, studentId string, academicYear, semester int) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Project{}).
		Joins("JOIN project_students ON project_students.project_id = projects.id").
		Joins("JOIN students ON students.id = project_students.student_id").
		Where("projects.academic_year = ? AND projects.semester = ?", academicYear, semester).
		Where("students.student_id = ?", studentId).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *studentRepositoryImpl) GetStudentByStudentId(ctx context.Context, studentId string) (*models.Student, error) {
	var student models.Student
	if err := r.db.WithContext(ctx).Where("student_id = ?", studentId).Order("created_at DESC").Preload("Course.Program").Preload("Program").First(&student).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *studentRepositoryImpl) GetStudentByStudentIdAndProgramIdOnCurrentYearAndSemester(ctx context.Context, studentId string, programId int, academicYear, semester int) (*models.Student, error) {
	var student models.Student
	if err := r.db.WithContext(ctx).
		Where("student_id = ? AND program_id = ? AND academic_year = ? AND semester = ? ", studentId, programId, academicYear, semester).
		Preload("Course.Program").
		Preload("Program").
		First(&student).Error; err != nil {
		return nil, err
	}

	return &student, nil
}

func (r *studentRepositoryImpl) GetStudentByProgramIdOnAcademicYearAndSemester(ctx context.Context, programId, academicYear, semester int) ([]models.Student, error) {
	var students []models.Student
	if err := r.db.WithContext(ctx).Where("program_id = ? AND academic_year = ? AND semester = ? ", programId, academicYear, semester).Preload("Course.Program").Preload("Program").Find(&students).Error; err != nil {
		return nil, err
	}
	return students, nil
}
