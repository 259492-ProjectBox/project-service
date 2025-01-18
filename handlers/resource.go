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

// DeleteProjectResource godoc
// @Summary Delete a project resource
// @Description Delete a project resource and its associated file
// @Tags Resource
// @Param id path int true "Resource ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projectResource/{id} [delete]
func (h *resourceHandler) DeleteProjectResource(c *gin.Context) {
	id := c.Param("id")
	detailedResource, err := h.resourceService.GetDetailedResourceByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	filePath := detailedResource.Resource.Path
	if err := h.resourceService.DeleteProjectResourceByID(c, id, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete resource record"})
		return
	}

	h.projectService.PublishProjectMessageToElasticSearch(c, "update", detailedResource.Project.ID)

	c.JSON(http.StatusOK, gin.H{"message": "Project Resource deleted successfully"})
}

// GetAssetResourceByProgramID godoc
// @Summary Get asset resources by program ID
// @Description Retrieve all asset resources associated with a specific program
// @Tags Resource
// @Param program_id path string true "Program ID"
// @Success 200 {array} dtos.AssetResource
// @Failure 500 {object} map[string]string
// @Router /assetResources/{program_id} [get]
func (h *resourceHandler) GetAssetResourceByProgramID(c *gin.Context) {
	programId := c.Param("program_id")
	assetResources, err := h.resourceService.GetAssetResourceByProgramID(c, programId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve asset resources"})
		return
	}
	c.JSON(http.StatusOK, assetResources)
}

// UploadAssetResource godoc
// @Summary Upload an asset resource
// @Description Upload a new asset resource for a specific program
// @Tags Resource
// @Param body body dtos.UploadAssetResource true "Upload Asset Resource Request"
// @Success 200 {object} dtos.AssetResource
// @Failure 400 {object} map[string]string
// @Router /assetResources [post]
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

// DeleteAssetResource godoc
// @Summary Delete an asset resource
// @Description Delete an asset resource by its ID
// @Tags Resource
// @Param id path int true "Asset Resource ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /assetResources/{id} [delete]
func (h *resourceHandler) DeleteAssetResource(c *gin.Context) {
	id := c.Param("id")
	err := h.resourceService.DeleteAssetResourceByID(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset Resource deleted successfully"})
}
