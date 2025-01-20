package services

import (
	"github.com/project-box/repositories"
)

type StudentService interface{}

type studentServiceImpl struct {
	studentRepo repositories.CourseRepository
}

func NewStudentService(studentRepo repositories.CourseRepository) StudentService {
	return &studentServiceImpl{
		studentRepo: studentRepo,
	}
}
