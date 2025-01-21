package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/project-box/dtos"
	"github.com/project-box/models"
	"github.com/project-box/services"
)

type ProgramHandler interface {
	GetPrograms(c *gin.Context)
	CreateProgram(c *gin.Context)
	UpdateProgram(c *gin.Context)
}

type programHandler struct {
	programService services.ProgramService
}

func NewProgramHandler(programService services.ProgramService) ProgramHandler {
	return &programHandler{
		programService: programService,
	}
}

// @Summary Get All Programs
// @Description Retrieves all programs from the database
// @Tags Program
// @Produce json
// @Success 200 {array} models.Program "Successfully fetched programs"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/programs [get]
func (h *programHandler) GetPrograms(c *gin.Context) {
	programs, err := h.programService.GetPrograms(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, programs)
}

// @Summary Create a New Program
// @Description Creates a new program in the database
// @Tags Program
// @Accept json
// @Produce json
// @Param program body dtos.CreateProgramRequest true "Program creation details"
// @Success 201 {object} models.Program "Successfully created program"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/programs [post]
func (h *programHandler) CreateProgram(c *gin.Context) {
	var program dtos.CreateProgramRequest
	if err := c.ShouldBindJSON(&program); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.programService.CreateProgram(c, &program); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, program)
}

// @Summary Update Program
// @Description update program
// @Tags Program
// @Accept json
// @Produce json
// @Param program body models.Program true "Program details"
// @Success 200 {object} map[string]interface{} "Successfully updated program"
// @Failure 400 {object} map[string]interface{} "Invalid request body or parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/programs [put]
func (h *programHandler) UpdateProgram(c *gin.Context) {
	var program models.Program
	if err := c.ShouldBindJSON(&program); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if program.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program ID"})
		return
	}
	if program.ProgramNameTH == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Program name th cannot be empty"})
		return
	}

	ctx := c.Request.Context()
	updatedProgram, err := h.programService.UpdateProgram(ctx, &program)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update program: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully updated program",
		"program": updatedProgram,
	})
}
