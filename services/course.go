package services

import (
	"context"

	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type CourseService interface {
	CreateCourse(ctx context.Context, course *models.Course) (*models.Course, error)
	UpdateCourse(ctx context.Context, course *models.Course) (*models.Course, error)
	DeleteCourse(ctx context.Context, courseId int) error
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

func (s *courseServiceImpl) CreateCourse(ctx context.Context, course *models.Course) (*models.Course, error) {
	return s.courseRepo.Create(ctx, course)
}

func (s *courseServiceImpl) UpdateCourse(ctx context.Context, course *models.Course) (*models.Course, error) {
	return s.courseRepo.Update(ctx, course.ID, course)
}

func (s *courseServiceImpl) DeleteCourse(ctx context.Context, courseId int) error {
	return s.courseRepo.Delete(ctx, courseId)
}

func (s *courseServiceImpl) FindCourseByCourseNo(ctx context.Context, courseNo string) (*models.Course, error) {
	return s.courseRepo.FindByCourseNo(ctx, courseNo)
}

func (s *courseServiceImpl) FindCourseByProgramID(ctx context.Context, programID int) (*models.Course, error) {
	return s.courseRepo.FindByProgramID(ctx, programID)
}
