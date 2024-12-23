package repositories

import (
	"github.com/project-box/models"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	repository[models.Employee]
	GetEmployeeByID(id int) (*models.Employee, error)
	GetEmployeeByMajorID(majorID int) ([]models.Employee, error)
	CreateEmployee(employee *models.Employee) error
	UpdateEmployee(updatedEmployee *models.Employee) (*models.Employee, error)
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

func (r *employeeRepositoryImpl) GetEmployeeByID(id int) (*models.Employee, error) {

	var employee *models.Employee
	if err := r.db.Where("id = ?", id).First(&employee).Error; err != nil {
		return nil, err
	}

	return employee, nil
}

// get employee by major id
func (r *employeeRepositoryImpl) GetEmployeeByMajorID(majorID int) ([]models.Employee, error) {
	var employees []models.Employee

	if err := r.db.Where("major_id = ?", majorID).Find(&employees).Error; err != nil {
		return nil, err
	}

	return employees, nil
}

// create employee
func (r *employeeRepositoryImpl) CreateEmployee(employee *models.Employee) error {
	if err := r.db.Create(employee).Error; err != nil {
		return err
	}
	return nil
}

// update employee
func (r *employeeRepositoryImpl) UpdateEmployee(updatedEmployee *models.Employee) (*models.Employee, error) {
	if err := r.db.Save(updatedEmployee).Error; err != nil {
		return nil, err
	}
	return updatedEmployee, nil
}
