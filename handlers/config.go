package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project-box/services"
)

type ConfigHandler interface {
	GetConfigByMajorIDHandler(c *gin.Context)
}

type configHandler struct {
	configService services.ConfigService
}

func NewConfigHandler(configService services.ConfigService) ConfigHandler {
	return &configHandler{
		configService: configService,
	}
}

// @Summary Get config by major ID
// @Description Fetches all config for a given major
// @Tags Config
// @Produce json
// @Param major_id path int true "Major ID"
// @Success 200 {object} []dtos.ConfigReponse "Successfully fetched config"
// @Failure 400 {object} map[string]interface{} "Invalid major ID"
// @Failure 404 {object} map[string]interface{} "Major not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /config/GetByMajorId/{major_id} [get]
func (h *configHandler) GetConfigByMajorIDHandler(c *gin.Context) {
	majorID, err := strconv.Atoi(c.Param("major_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid major ID"})
		return
	}

	configs, err := h.configService.GetConfigByMajorID(majorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(configs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No configurations found for the given major ID"})
		return
	}

	c.JSON(http.StatusOK, configs)
}
