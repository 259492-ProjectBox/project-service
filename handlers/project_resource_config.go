package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project-box/models"
	"github.com/project-box/services"
)

type ProjectResourceConfigHandler interface {
	GetProjectResourceConfigsByProgramId(c *gin.Context)
	UpsertProjectResourceConfig(c *gin.Context)
}

type projectResourceConfigHandler struct {
	projectResourceConfigService services.ProjectResourceConfigService
}

func NewProjectResourceConfigHandler(projectResourceConfigService services.ProjectResourceConfigService) ProjectResourceConfigHandler {
	return &projectResourceConfigHandler{
		projectResourceConfigService: projectResourceConfigService,
	}
}

// @Summary Get Project Resource Config by Program ID
// @Description Fetch all project resource configurations for a given program ID
// @Tags ProjectResourceConfig
// @Produce json
// @Param program_id path int true "Program ID"
// @Success 200 {array} []models.ProjectResourceConfig "Successfully fetched configurations"
// @Failure 400 {object} map[string]interface{} "Invalid program ID"
// @Failure 404 {object} map[string]interface{} "No configurations found for the given program ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/projectResourceConfigs/program/{program_id} [get]
func (h *projectResourceConfigHandler) GetProjectResourceConfigsByProgramId(c *gin.Context) {
	programId, err := strconv.Atoi(c.Param("program_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program ID"})
		return
	}

	configs, err := h.projectResourceConfigService.GetProjectResourceConfigsByProgramId(programId)
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

// @Summary Upsert Project Resource Configurations
// @Description Insert or update project resource configurations. If an ID is provided, it updates the configuration; otherwise, it inserts a new configuration.
// @Tags ProjectResourceConfig
// @Accept json
// @Produce json
// @Param configs body models.ProjectResourceConfig true "configuration to upsert"
// @Success 200 {object} map[string]interface{} "Successfully upsert configurations"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/projectResourceConfigs [put]
func (h *projectResourceConfigHandler) UpsertProjectResourceConfig(c *gin.Context) {
	var projectResourceConfig models.ProjectResourceConfig
	if err := c.ShouldBindJSON(&projectResourceConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.projectResourceConfigService.UpsertResourceProjectConfig(&projectResourceConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully upsert configurations"})
}
