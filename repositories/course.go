package repositories

import (
	"context"
	"fmt"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type CourseRepository interface {
	repository[models.Course]
	GetByCourseAndSemester(ctx context.Context, courseId int, semester int) (*models.Course, error)
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
	fmt.Println("courseID", courseID)
	fmt.Println("semester", semester)
	filters := map[string]interface{}{"id": courseID, "semester": semester}
	var course models.Course
	if err := r.db.WithContext(ctx).Where(filters).First(&course).Error; err != nil {
		return nil, err
	}
	return &course, nil
}
