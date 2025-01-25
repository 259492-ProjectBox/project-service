package repositories

import (
	"context"
	"fmt"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type ProjectNumberCounterRepository interface {
	repository[models.ProjectNumberCounter]
	GetNextProjectNumber(ctx context.Context, tx *gorm.DB, academicYear, semester, courseID int) (int, error)
}

type projectNumberCounterRepositoryImpl struct {
	db *gorm.DB
	*repositoryImpl[models.ProjectNumberCounter]
}

func NewProjectNumberCounterRepository(db *gorm.DB) ProjectNumberCounterRepository {
	return &projectNumberCounterRepositoryImpl{
		db:             db,
		repositoryImpl: newRepository[models.ProjectNumberCounter](db),
	}
}

func (r *projectNumberCounterRepositoryImpl) GetNextProjectNumber(ctx context.Context, tx *gorm.DB, academicYear, semester, courseID int) (int, error) {
	var counter models.ProjectNumberCounter
	if tx == nil {
		tx = r.db
	}
	result := tx.WithContext(ctx).Where("academic_year = ? AND semester = ? AND course_id = ?", academicYear, semester, courseID).First(&counter)
	if result.RowsAffected == 0 {
		counter = models.ProjectNumberCounter{
			AcademicYear: academicYear,
			Semester:     semester,
			CourseID:     courseID,
			Number:       1,
		}
		if err := tx.Create(&counter).Error; err != nil {
			return 0, fmt.Errorf("could not create project number counter: %v", err)
		}
		return counter.Number, nil
	}

	counter.Number++
	if err := tx.Save(&counter).Error; err != nil {
		return 0, fmt.Errorf("could not update project number counter: %v", err)
	}
	return counter.Number, nil
}
