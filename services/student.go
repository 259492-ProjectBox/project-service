package services

import (
	"context"

	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type StudentService interface {
	CreateStudents(ctx context.Context, students []models.Student) error
}

type studentServiceImpl struct {
	studentRepo repositories.StudentRepository
}

func NewStudentService(studentRepo repositories.StudentRepository) StudentService {
	return &studentServiceImpl{
		studentRepo: studentRepo,
	}
}

// CreateStudents creates multiple student records in the database.
func (s *studentServiceImpl) CreateStudents(ctx context.Context, students []models.Student) error {
	// Validate input to ensure students slice is not empty
	if len(students) == 0 {
		return nil // No students to create
	}

	if err := s.studentRepo.CreateMany(ctx, students); err != nil {
		return err
	}

	return nil
}
