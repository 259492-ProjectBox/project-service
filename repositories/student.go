package repositories

import (
	"github.com/project-box/models"
	"gorm.io/gorm"
)

type StudentRepository interface {
	repository[models.Student]
}

type studentRepositoryImpl struct {
	db *gorm.DB
	*repositoryImpl[models.Student]
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepositoryImpl{
		db:             db,
		repositoryImpl: newRepository[models.Student](db),
	}
}
