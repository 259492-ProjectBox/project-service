package repositories

import (
	"fmt"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type ProjectNumberCounterRepository interface {
	repository[models.ProjectNumberCounter]
	GetNextProjectNumber(academicYear, semester, courseID int) (int, error)
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

func (r *projectNumberCounterRepositoryImpl) GetNextProjectNumber(academicYear, semester, courseID int) (int, error) {
	var counter models.ProjectNumberCounter

	err := r.db.Where("academic_year = ? AND semester = ? AND course_id = ?", academicYear, semester, courseID).First(&counter).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, fmt.Errorf("could not query project number counter: %v", err)
	}

	if err == gorm.ErrRecordNotFound {
		counter = models.ProjectNumberCounter{
			AcademicYear: academicYear,
			Semester:     semester,
			CourseID:     courseID,
			Number:       1,
		}
		if err := r.db.Create(&counter).Error; err != nil {
			return 0, fmt.Errorf("could not create project number counter: %v", err)
		}
		return counter.Number, nil
	}

	// If counter exists, increment the number and save it back
	counter.Number++
	if err := r.db.Save(&counter).Error; err != nil {
		return 0, fmt.Errorf("could not update project number counter: %v", err)
	}

	return counter.Number, nil
}
