package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project-box/dtos"
	"github.com/project-box/services"
)

type ProjectHandler interface {
	CreateProject(c *gin.Context)
	UpdateProject(c *gin.Context)
	DeleteProject(c *gin.Context)
}

type projectHandler struct {
	projectService services.ProjectService
}

func NewProjectHandler(projectService services.ProjectService) ProjectHandler {
	return &projectHandler{
		projectService: projectService,
	}
}

// @Summary Create a new project
// @Description Creates a new project with the provided data
// @Tags Project
// @Accept  multipart/form-data
// @Produce  json
// @Param urls formData string false "List of URLs"
// @Param titles formData string false "List of Titles"
// @Param files formData file false "Upload Files"
// @Param project formData string true "Project Data"
// @Success 201 {object} models.Project "Successfully created project"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/projects [post]
func (h *projectHandler) CreateProject(c *gin.Context) {
	req := &dtos.CreateProjectRequest{}
	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := h.projectService.CreateProjectWithFiles(c, req.Project, req.ProjectResources, req.Files)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, *project)
}

// @Summary Update an existing project
// @Description Updates a project by its ID with the provided data
// @Tags Project
// @Accept  multipart/form-data
// @Produce  json
// @Param urls formData string false "List of URLs"
// @Param titles formData string false "List of Titles"
// @Param files formData file false "Upload Files"
// @Param project formData string true "Project Data"
// @Success 200 {object} models.Project "Successfully updated project"
// @Failure 400 {object} map[string]interface{} "Invalid project ID or request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/projects/{id} [put]
func (h *projectHandler) UpdateProject(c *gin.Context) {
	req := &dtos.UpdateProjectRequest{}
	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := h.projectService.UpdateProjectWithFiles(c, req.Project, req.ProjectResources, req.Files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

// DeleteProject deletes a project by its ID
// @Summary Delete a project by ID
// @Description Deletes the specified project using its ID
// @Tags Project
// @Produce  json
// @Param id path int true "Project ID"
// @Success 200 {object} map[string]interface{} "Project deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid project ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/projects/{id} [delete]
func (h *projectHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	if err := h.projectService.DeleteProject(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}
