package repositories

import (
	"context"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type SectionRepository interface {
	repository[models.Section]
	GetByCourseAndSemester(ctx context.Context, courseId int, semester int) (*models.Section, error)
}

type sectionRepositoryImpl struct {
	db *gorm.DB
	*repositoryImpl[models.Section]
}

func NewSectionRepository(db *gorm.DB) SectionRepository {
	return &sectionRepositoryImpl{
		db:             db,
		repositoryImpl: newRepository[models.Section](db),
	}
}

func (r *sectionRepositoryImpl) GetByCourseAndSemester(ctx context.Context, courseID int, semester int) (*models.Section, error) {
	filters := map[string]interface{}{"course_id": courseID, "semester": semester}
	var section models.Section
	if err := r.db.WithContext(ctx).Where(filters).First(&section).Error; err != nil {
		return nil, err
	}
	return &section, nil
}
