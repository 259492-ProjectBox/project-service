package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/project-box/dtos"
	"github.com/project-box/models"
	"github.com/project-box/services"
)

type MajorHandler interface {
	GetAllMajorHandler(c *gin.Context)
	CreateMajorHandler(c *gin.Context)
	UpdateMajorNameHandler(c *gin.Context)
}

type majorHandler struct {
	majorService services.MajorService
}

func NewMajorHandler(majorService services.MajorService) MajorHandler {
	return &majorHandler{
		majorService: majorService,
	}
}

// @Summary Get all major
// @Description Fetches all major
// @Tags Major
// @Produce json
// @Success 200 {object} []models.Major "Successfully fetched major"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /major/GetAllMajor [get]
func (h *majorHandler) GetAllMajorHandler(c *gin.Context) {
	majors, err := h.majorService.GetAllMajorService(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, majors)
}

// @Summary Create major
// @Description Create a new major
// @Tags Major
// @Accept json
// @Produce json
// @Param major body dtos.CreateMajorRequest true "Major object"
// @Success 201 {object} models.Major "Successfully created major"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /major [post]
func (h *majorHandler) CreateMajorHandler(c *gin.Context) {
	var major dtos.CreateMajorRequest
	if err := c.ShouldBindJSON(&major); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.majorService.CreateMajorService(c, &major); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, major)
}

// @Summary Update major name
// @Description Update the name of a major
// @Tags Major
// @Accept json
// @Produce json
// @Param major body models.Major true "Major object containing ID and name"
// @Success 200 {object} models.Major "Successfully updated major name"
// @Failure 400 {object} map[string]interface{} "Invalid request body or parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /major/UpdateMajorName [put]
func (h *majorHandler) UpdateMajorNameHandler(c *gin.Context) {
	var major models.Major
	// Bind JSON to Major object
	if err := c.ShouldBindJSON(&major); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Validate ID and Name
	if major.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid major ID"})
		return
	}
	if major.MajorName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Major name cannot be empty"})
		return
	}

	// Call service to update the major name
	ctx := c.Request.Context()
	if err := h.majorService.UpdateMajorNameService(ctx, major.ID, major.MajorName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update major name: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully updated major name",
		"major":   major,
	})
}
