package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project-box/dtos"
	"github.com/project-box/services"
)

type StaffHandler interface {
	GetStaffById(c *gin.Context)
	GetStaffByProgramId(c *gin.Context)
	CreateStaff(c *gin.Context)
	UpdateStaff(c *gin.Context)
}

type staffHandler struct {
	staffService services.StaffService
}

func NewStaffHandler(staffService services.StaffService) StaffHandler {
	return &staffHandler{
		staffService: staffService,
	}
}

// @Summary Get staff by ID
// @Description Fetches an staff by their ID
// @Tags Staff
// @Produce  json
// @Param id path int true "Staff ID"
// @Success 200 {object} dtos.StaffResponse "Successfully retrieved staff"
// @Failure 400 {object} map[string]interface{} "Invalid staff ID"
// @Failure 404 {object} map[string]interface{} "Staff not found"
// @Router /v1/staffs/{id} [get]
func (h *staffHandler) GetStaffById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid staff ID"})
		return
	}
	staff, err := h.staffService.GetStaffById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}
	c.JSON(http.StatusOK, staff)
}

// @Summary Get staffs by program ID
// @Description Fetches all staffs for a given program
// @Tags Staff
// @Produce  json
// @Param program_id path int true "Program ID"
// @Success 200 {object} []dtos.StaffResponse "Successfully retrieved staffs"
// @Failure 400 {object} map[string]interface{} "Invalid program ID"
// @Failure 404 {object} map[string]interface{} "Staffs not found"
// @Router /v1/staffs/program/{program_id} [get]
func (h *staffHandler) GetStaffByProgramId(c *gin.Context) {
	programId, err := strconv.Atoi(c.Param("program_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program ID"})
		return
	}
	staffs, err := h.staffService.GetStaffByProgramId(c, programId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staffs not found"})
		return
	}
	c.JSON(http.StatusOK, staffs)
}

// @Summary Create a new staff
// @Description Creates a new staff
// @Tags Staff
// @Accept  json
// @Produce  json
// @Param staff body dtos.CreateStaffRequest true "Staff Data"
// @Success 201 {object} dtos.StaffResponse "Successfully created staff"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/staffs [post]
func (h *staffHandler) CreateStaff(c *gin.Context) {
	req := &dtos.CreateStaffRequest{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	staff, err := h.staffService.CreateStaff(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, staff)
}

// @Summary Update an existing staff
// @Description Updates an staff by their ID with the provided data
// @Tags Staff
// @Accept  json
// @Produce  json
// @Param staff body dtos.UpdateStaffRequest true "Updated Staff Data"
// @Success 200 {object} dtos.StaffResponse "Successfully updated staff"
// @Failure 400 {object} map[string]interface{} "Invalid staff ID or request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/staffs [put]
func (h *staffHandler) UpdateStaff(c *gin.Context) {
	staff := &dtos.UpdateStaffRequest{}
	if err := c.ShouldBindJSON(&staff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedStaff, err := h.staffService.UpdateStaff(c, staff)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedStaff)
}
