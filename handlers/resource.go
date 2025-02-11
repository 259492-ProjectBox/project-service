package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/project-box/dtos"
	"github.com/project-box/services"
)

type ResourceHandler interface {
	GetAssetResourceByProgramID(c *gin.Context)
	UploadAssetResource(c *gin.Context)
	DeleteAssetResource(c *gin.Context)
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

func (h *resourceHandler) GetAssetResourceByProgramID(c *gin.Context) {
	programId := c.Param("program_id")
	assetResources, err := h.resourceService.GetAssetResourceByProgramID(c, programId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve asset resources"})
		return
	}
	c.JSON(http.StatusOK, assetResources)
}

func (h *resourceHandler) UploadAssetResource(c *gin.Context) {
	req := &dtos.UploadAssetResource{}
	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	assetResource, err := h.resourceService.UploadAssetResource(c, req.File, req.AssetResource)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, *assetResource)
}

func (h *resourceHandler) DeleteAssetResource(c *gin.Context) {
	id := c.Param("id")
	err := h.resourceService.DeleteAssetResourceByID(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset Resource deleted successfully"})
}
