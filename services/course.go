package services

import (
	"github.com/project-box/repositories"
)

type CourseService interface {
}

type courseServiceImpl struct {
	courseRepo repositories.CourseRepository
}

func NewCourseService(courseRepo repositories.CourseRepository) CourseService {
	return &courseServiceImpl{
		courseRepo: courseRepo,
	}
}
