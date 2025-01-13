package services

import (
	"context"

	"github.com/project-box/dtos"
	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type EmployeeService interface {
	CreateEmployee(ctx context.Context, employee *dtos.CreateEmployeeRequest) (*dtos.EmployeeResponse, error)
	GetEmployeeByID(ctx context.Context, id int) (*dtos.EmployeeResponse, error)
	UpdateEmployee(ctx context.Context, employee *dtos.UpdateEmployeeRequest) (*dtos.EmployeeResponse, error)
	DeleteEmployee(ctx context.Context, id int) error
	GetEmployeeByProgramId(ctx context.Context, programId int) ([]dtos.EmployeeResponse, error)
}

type employeeServiceImpl struct {
	employeeRepo repositories.EmployeeRepository
}

func NewEmployeeService(employeeRepo repositories.EmployeeRepository) EmployeeService {
	return &employeeServiceImpl{
		employeeRepo: employeeRepo,
	}
}

func (s *employeeServiceImpl) CreateEmployee(ctx context.Context, employee *dtos.CreateEmployeeRequest) (*dtos.EmployeeResponse, error) {

	employeeModel := &models.Employee{
		Prefix:    employee.Prefix,
		FirstName: employee.FirstName,
		LastName:  employee.LastName,
		Email:     employee.Email,
		ProgramID: employee.ProgramID,
	}

	employeeModel, err := s.employeeRepo.Create(ctx, employeeModel)
	if err != nil {
		return nil, err
	}

	return &dtos.EmployeeResponse{
		ID:        employeeModel.ID,
		Prefix:    employeeModel.Prefix,
		FirstName: employeeModel.FirstName,
		LastName:  employeeModel.LastName,
		Email:     employeeModel.Email,
		ProgramID: employeeModel.ProgramID,
	}, nil

}

func (s *employeeServiceImpl) GetEmployeeByID(ctx context.Context, id int) (*dtos.EmployeeResponse, error) {
	employee, err := s.employeeRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dtos.EmployeeResponse{
		ID:        employee.ID,
		Prefix:    employee.Prefix,
		FirstName: employee.FirstName,
		LastName:  employee.LastName,
		Email:     employee.Email,
		ProgramID: employee.ProgramID,
	}, nil
}

func (s *employeeServiceImpl) UpdateEmployee(ctx context.Context, employee *dtos.UpdateEmployeeRequest) (*dtos.EmployeeResponse, error) {
	// convert from dto to model
	updatedEmployee := &models.Employee{
		ID:        employee.ID,
		Prefix:    employee.Prefix,
		FirstName: employee.FirstName,
		LastName:  employee.LastName,
		Email:     employee.Email,
		ProgramID: employee.ProgramID,
	}
	updatedEmployee, err := s.employeeRepo.Update(ctx, employee.ID, updatedEmployee)
	if err != nil {
		return nil, err
	}

	// convert to response
	return &dtos.EmployeeResponse{
		ID:        updatedEmployee.ID,
		Prefix:    updatedEmployee.Prefix,
		FirstName: updatedEmployee.FirstName,
		LastName:  updatedEmployee.LastName,
		Email:     updatedEmployee.Email,
		ProgramID: updatedEmployee.ProgramID,
	}, nil

}

func (s *employeeServiceImpl) DeleteEmployee(ctx context.Context, id int) error {
	return s.employeeRepo.Delete(ctx, id)
}

func (s *employeeServiceImpl) GetEmployeeByProgramId(ctx context.Context, programId int) ([]dtos.EmployeeResponse, error) {
	employees, err := s.employeeRepo.GetEmployeeByProgramID(programId)
	if err != nil {
		return nil, err
	}

	// convert to response
	var employeeResponses []dtos.EmployeeResponse
	for _, employee := range employees {
		employeeResponses = append(employeeResponses, dtos.EmployeeResponse{
			ID:        employee.ID,
			Prefix:    employee.Prefix,
			FirstName: employee.FirstName,
			LastName:  employee.LastName,
			Email:     employee.Email,
			ProgramID: employee.ProgramID,
		})
	}

	return employeeResponses, nil
}
