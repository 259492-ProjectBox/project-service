package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project-box/dtos"
	"github.com/project-box/services"
)

type EmployeeHandler interface {
	GetEmployeeByIDHandler(c *gin.Context)
	GetEmployeeByMajorIDHandler(c *gin.Context)
	CreateEmployeeHandler(c *gin.Context)
	UpdateEmployeeHandler(c *gin.Context)
}

type employeeHandler struct {
	employeeService services.EmployeeService
}

func NewEmployeeHandler(employeeService services.EmployeeService) EmployeeHandler {
	return &employeeHandler{
		employeeService: employeeService,
	}
}

// @Summary Get employee by ID
// @Description Fetches an employee by their ID
// @Tags Employee
// @Produce  json
// @Param id path int true "Employee ID"
// @Success 200 {object} dtos.EmployeeResponse "Successfully retrieved employee"
// @Failure 400 {object} map[string]interface{} "Invalid employee ID"
// @Failure 404 {object} map[string]interface{} "Employee not found"
// @Router /employee/{id} [get]
func (h *employeeHandler) GetEmployeeByIDHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}
	employee, err := h.employeeService.GetEmployeeByIDService(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	c.JSON(http.StatusOK, employee)
}

// @Summary Get employees by major ID
// @Description Fetches all employees for a given major
// @Tags Employee
// @Produce  json
// @Param major_id path int true "Major ID"
// @Success 200 {object} []dtos.EmployeeResponse "Successfully retrieved employees"
// @Failure 400 {object} map[string]interface{} "Invalid major ID"
// @Failure 404 {object} map[string]interface{} "Employees not found"
// @Router /employee/GetByMajorID/{major_id} [get]
func (h *employeeHandler) GetEmployeeByMajorIDHandler(c *gin.Context) {
	majorID, err := strconv.Atoi(c.Param("major_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid major ID"})
		return
	}
	employees, err := h.employeeService.GetEmployeeByMajorIDService(c, majorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employees not found"})
		return
	}
	c.JSON(http.StatusOK, employees)
}

// @Summary Create a new employee
// @Description Creates a new employee
// @Tags Employee
// @Accept  json
// @Produce  json
// @Param employee body dtos.CreateEmployeeRequest true "Employee Data"
// @Success 201 {object} dtos.EmployeeResponse "Successfully created employee"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /employee [post]
func (h *employeeHandler) CreateEmployeeHandler(c *gin.Context) {
	req := &dtos.CreateEmployeeRequest{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	employee, err := h.employeeService.CreateEmployeeService(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, employee)
}

// @Summary Update an existing employee
// @Description Updates an employee by their ID with the provided data
// @Tags Employee
// @Accept  json
// @Produce  json
// @Param employee body dtos.UpdateEmployeeRequest true "Updated Employee Data"
// @Success 200 {object} dtos.EmployeeResponse "Successfully updated employee"
// @Failure 400 {object} map[string]interface{} "Invalid employee ID or request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /employee [put]
func (h *employeeHandler) UpdateEmployeeHandler(c *gin.Context) {
	employee := &dtos.UpdateEmployeeRequest{}
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedEmployee, err := h.employeeService.UpdateEmployeeService(c, employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedEmployee)
}
