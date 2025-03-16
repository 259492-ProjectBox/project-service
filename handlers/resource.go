package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/project-box/services"
)

type ResourceHandler interface {
	DeleteProjectResource(c *gin.Context)
}

type resourceHandler struct {
	minioClient     *minio.Client
	bucketName      string
	resourceService services.ResourceService
	projectService  services.ProjectService
}

func NewResourceHandler(minioClient *minio.Client, resourceService services.ResourceService, projectService services.ProjectService) ResourceHandler {
	return &resourceHandler{
		minioClient:     minioClient,
		bucketName:      os.Getenv("MINIO_BUCKET"),
		resourceService: resourceService,
		projectService:  projectService,
	}
}

// @Summary Delete a project resource
// @Description Deletes a project resource by its ID
// @Tags Resource
// @Produce json
// @Param id path string true "Resource ID"
// @Success 200 {object} map[string]interface{} "Project Resource deleted successfully"
// @Failure 404 {object} map[string]interface{} "Resource not found"
// @Failure 500 {object} map[string]interface{} "Failed to delete resource record"
// @Router /v1/projectResources/{id} [delete]
func (h *resourceHandler) DeleteProjectResource(c *gin.Context) {
	id := c.Param("id")
	detailedResource, err := h.resourceService.GetDetailedResourceByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	filePath := detailedResource.ProjectResource.Path
	if err := h.resourceService.DeleteProjectResourceByID(c, id, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete resource record"})
		return
	}

	err = h.projectService.PublishProjectMessageToElasticSearch(c, "update", detailedResource.Project.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to publish project message to ElasticSearch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project Resource deleted successfully"})
}
