package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project-box/services"
)

type ConfigHandler interface {
	GetConfigByProgramId(c *gin.Context)
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
// @Success 200 {object} []dtos.ConfigReponse "Successfully fetched config"
// @Failure 400 {object} map[string]interface{} "Invalid program ID"
// @Failure 404 {object} map[string]interface{} "Program not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/config/program/{program_id} [get]
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
