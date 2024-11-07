package repositories

import (
	"context"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type SectionRepository interface {
	repository[models.Section]
	GetByCourseAndSectionAndSemester(ctx context.Context, courseId int, sectionId *int, semester int) (*models.Section, error)
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

func (r *sectionRepositoryImpl) GetByCourseAndSectionAndSemester(ctx context.Context, courseID int, sectionID *int, semester int) (*models.Section, error) {
	filters := map[string]interface{}{"course_id": courseID, "semester": semester}
	if sectionID != nil {
		filters["id"] = *sectionID
	}

	var section models.Section
	if err := r.db.WithContext(ctx).Where(filters).First(&section).Error; err != nil {
		return nil, err
	}
	return &section, nil
}
