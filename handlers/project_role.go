package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project-box/services"
)

type ProjectRoleHandler interface {
	GetByProgram(c *gin.Context)
}

type projectRoleHandler struct {
	projectRoleService services.ProjectRoleService
}

func NewProjectRoleHandler(projectRoleService services.ProjectRoleService) ProjectRoleHandler {
	return &projectRoleHandler{
		projectRoleService: projectRoleService,
	}
}

// @Summary Get all project roles by program ID
// @Description Retrieves all project roles for a given program ID
// @Tags ProjectRole
// @Produce json
// @Param program_id path int true "Program ID"
// @Success 200 {array} models.ProjectRole "Successfully retrieved project roles"
// @Failure 400 {object} map[string]interface{} "Invalid program ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/projectRoles/program/{program_id} [get]
func (h *projectRoleHandler) GetByProgram(c *gin.Context) {
	programIdStr := c.Param("program_id")
	programId, err := strconv.Atoi(programIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program ID"})
		return
	}

	projectRoles, err := h.projectRoleService.GetAllProjectRolesByProgramId(c.Request.Context(), programId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projectRoles)
}
