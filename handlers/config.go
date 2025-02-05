package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project-box/models"
	"github.com/project-box/services"
	"gorm.io/gorm"
)

type ConfigHandler interface {
	GetConfigByProgramId(c *gin.Context)
	UpsertConfig(c *gin.Context)
	DeleteConfig(c *gin.Context)
}

type configHandler struct {
	configService services.ConfigService
}

func NewConfigHandler(configService services.ConfigService) ConfigHandler {
	return &configHandler{
		configService: configService,
	}
}

// @Summary Get config by program ID
// @Description Fetches all config for a given program
// @Tags Config
// @Produce json
// @Param program_id path int true "Program ID"
// @Success 200 {object} []models.Config "Successfully fetched config"
// @Failure 400 {object} map[string]interface{} "Invalid program ID"
// @Failure 404 {object} map[string]interface{} "Program not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/configs/program/{program_id} [get]
func (h *configHandler) GetConfigByProgramId(c *gin.Context) {
	programId, err := strconv.Atoi(c.Param("program_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program ID"})
		return
	}

	configs, err := h.configService.GetConfigByProgramId(programId)
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

// @Summary Upsert config for a program
// @Description Creates a new config or updates an existing config for the given program
// @Tags Config
// @Produce json
// @Param config body models.Config true "Config details"
// @Success 200 {object} models.Config "Successfully upserted config"
// @Failure 400 {object} map[string]interface{} "Invalid program ID or config data"
// @Failure 404 {object} map[string]interface{} "Program not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/configs [put]
func (h *configHandler) UpsertConfig(c *gin.Context) {
	config := &models.Config{}
	if err := c.ShouldBindJSON(config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid config data"})
		return
	}

	config, err := h.configService.UpsertConfig(c, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, *config)
}

// @Summary Delete config by ID
// @Description Deletes a configuration by its ID
// @Tags Config
// @Produce json
// @Param id path int true "Config ID"
// @Success 200 {object} map[string]interface{} "Successfully deleted config"
// @Failure 400 {object} map[string]interface{} "Invalid config ID"
// @Failure 404 {object} map[string]interface{} "Config not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/configs/{id} [delete]
func (h *configHandler) DeleteConfig(c *gin.Context) {
	configIdStr := c.Param("id")
	configId, err := strconv.Atoi(configIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid config ID"})
		return
	}

	err = h.configService.DeleteConfig(c.Request.Context(), configId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Config not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted config"})
}
