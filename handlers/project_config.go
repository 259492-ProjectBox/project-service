package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project-box/dtos"
	"github.com/project-box/services"
)

type ProjectConfigHandler interface {
	GetProjectConfigByProgramId(c *gin.Context)
	UpsertProjectConfig(c *gin.Context)
}

type projectconfigHandler struct {
	projectconfigService services.ProjectConfigService
}

func NewProjectConfigHandler(projectconfigService services.ProjectConfigService) ProjectConfigHandler {
	return &projectconfigHandler{
		projectconfigService: projectconfigService,
	}
}

// @Summary Get config by program ID
// @Description Fetches all config for a given program
// @Tags ProjectConfig
// @Produce json
// @Param program_id path int true "Program ID"
// @Success 200 {object} []dtos.ProjectConfigResponse "Successfully fetched config"
// @Failure 400 {object} map[string]interface{} "Invalid program ID"
// @Failure 404 {object} map[string]interface{} "Program not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/projectConfigs/program/{program_id} [get]
func (h *projectconfigHandler) GetProjectConfigByProgramId(c *gin.Context) {
	programId, err := strconv.Atoi(c.Param("program_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program ID"})
		return
	}

	configs, err := h.projectconfigService.GetProjectConfigByProgramId(programId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(configs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No configurations found for the given program ID"})
		return
	}

	c.JSON(http.StatusOK, configs)
}

// @Summary Upsert project config
// @Description Update all project config if ID is provided, otherwise insert new project config
// @Tags ProjectConfig
// @Accept json
// @Produce json
// @Param configs body []dtos.ProjectConfigUpsertRequest true "Configurations"
// @Success 200 {object} map[string]interface{} "Successfully updated config"
// @Failure 400 {object} map[string]interface{} "Invalid program ID"
// @Failure 404 {object} map[string]interface{} "Program not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/projectConfig [post]
func (h *projectconfigHandler) UpsertProjectConfig(c *gin.Context) {
	var configs []dtos.ProjectConfigUpsertRequest
	if err := c.ShouldBindJSON(&configs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.projectconfigService.UpdateProjectConfig(configs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully updated config"})
}
