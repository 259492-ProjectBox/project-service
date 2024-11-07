package repositories

import (
	"github.com/project-box/models"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	repository[models.Employee]
}

type employeeRepositoryImpl struct {
	db *gorm.DB
	*repositoryImpl[models.Employee]
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepositoryImpl{
		db:             db,
		repositoryImpl: newRepository[models.Employee](db),
	}
}
