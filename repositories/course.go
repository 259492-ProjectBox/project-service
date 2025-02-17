package repositories

import (
	"context"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type CourseRepository interface {
	repository[models.Course]
	GetByCourseAndSemester(ctx context.Context, courseId int, semester int) (*models.Course, error)
	FindByCourseNo(ctx context.Context, courseNo string) (*models.Course, error)
	FindByProgramID(ctx context.Context, programID int) (*models.Course, error)
}

type courseRepositoryImpl struct {
	db *gorm.DB
	*repositoryImpl[models.Course]
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepositoryImpl{
		db:             db,
		repositoryImpl: newRepository[models.Course](db),
	}
}

func (r *courseRepositoryImpl) GetByCourseAndSemester(ctx context.Context, courseID int, semester int) (*models.Course, error) {
	filters := map[string]interface{}{"id": courseID, "semester": semester}
	var course models.Course
	if err := r.db.WithContext(ctx).Where(filters).First(&course).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *courseRepositoryImpl) FindByCourseNo(ctx context.Context, courseNo string) (*models.Course, error) {
	var course models.Course
	if err := r.db.WithContext(ctx).Where("course_no = ?", courseNo).First(&course).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *courseRepositoryImpl) FindByProgramID(ctx context.Context, programID int) (*models.Course, error) {
	var courses *models.Course
	if err := r.db.WithContext(ctx).Where("program_id = ?", programID).First(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}
