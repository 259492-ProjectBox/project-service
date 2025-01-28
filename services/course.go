package services

import (
	"context"

	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type CourseService interface {
	FindCourseByCourseNo(ctx context.Context, courseNo string) (*models.Course, error)
	FindCourseByProgramID(ctx context.Context, programID int) (*models.Course, error)
}

type courseServiceImpl struct {
	courseRepo repositories.CourseRepository
}

func NewCourseService(courseRepo repositories.CourseRepository) CourseService {
	return &courseServiceImpl{
		courseRepo: courseRepo,
	}
}

func (s *courseServiceImpl) FindCourseByCourseNo(ctx context.Context, courseNo string) (*models.Course, error) {
	return s.courseRepo.FindByCourseNo(ctx, courseNo)
}

func (s *courseServiceImpl) FindCourseByProgramID(ctx context.Context, programID int) (*models.Course, error) {
	return s.courseRepo.FindByProgramID(ctx, programID)
}
