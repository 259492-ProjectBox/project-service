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
	// UploadResource(c *gin.Context)
	GetAssetResourceByProgramID(c *gin.Context)
	UploadAssetResource(c *gin.Context)
	DeleteAssetResource(c *gin.Context)
	// GetResourceByID(c *gin.Context)
	DeleteProjectResource(c *gin.Context)
	// GetResourcesByProjectID(c *gin.Context)
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

// @Summary Upload a file and create a resource
// @Description Upload a file to MinIO and save its information as a resource
// @Tags Resource
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param project_id formData int true "Project ID"
// @Param resource_type_id formData int true "Resource Type ID"
// @Success 201 {object} models.UploadResourceResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /resource [post]
// func (h *resourceHandler) UploadResource(c *gin.Context) {
// 	projectID, err := strconv.Atoi(c.PostForm("project_id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project_id"})
// 		return
// 	}

// 	resourceTypeID, err := strconv.Atoi(c.PostForm("resource_type_id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource_type_id"})
// 		return
// 	}

// 	file, header, err := c.Request.FormFile("file")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
// 		return
// 	}
// 	defer file.Close()

// 	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)

// 	contentType := "application/pdf" // You might want to make this dynamic based on file type
// 	_, err = h.minioClient.PutObject(
// 		context.Background(),
// 		h.bucketName,
// 		filename,
// 		file,
// 		header.Size,
// 		minio.PutObjectOptions{ContentType: contentType},
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
// 		return
// 	}

// 	// Generate presigned URL
// 	url, err := h.minioClient.PresignedGetObject(context.Background(), h.bucketName, filename, time.Hour*24*7, nil)
// 	if err != nil {
// 		// Cleanup the uploaded file if URL generation fails
// 		h.minioClient.RemoveObject(context.Background(), h.bucketName, filename, minio.RemoveObjectOptions{})
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate URL"})
// 		return
// 	}

// 	// Create resource in database
// 	resource, err := h.resourceService.CreateResource(c, &models.Resource{
// 		Title:          filename,
// 		ProjectID:      projectID,
// 		ResourceTypeID: resourceTypeID,
// 		URL:            url.String(),
// 	})

// 	if err != nil {
// 		// Cleanup the uploaded file if database save fails
// 		h.minioClient.RemoveObject(context.Background(), h.bucketName, filename, minio.RemoveObjectOptions{})
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save resource: %v", err)})
// 		return
// 	}
// 	response := models.UploadResourceResponse{
// 		ID:             resource.ID,
// 		Title:          filename,
// 		ProjectID:      projectID,
// 		ResourceTypeID: resourceTypeID,
// 		URL:            url.String(),
// 	}

// 	c.JSON(http.StatusCreated, response)
// }

// GetResourceByID godoc
// @Summary Get a resource by ID
// @Description Get a resource by its ID
// @Tags Resource
// @Produce json
// @Param id path int true "Resource ID"
// @Success 200 {object} models.UploadResourceResponse
// @Failure 404 {object} map[string]string
// @Router /resource/{id} [get]
// func (h *resourceHandler) GetResourceByID(c *gin.Context) {
// 	id := c.Param("id")
// 	resource, err := h.resourceService.GetResourceByID(c, id)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
// 		return
// 	}

// 	response := models.UploadResourceResponse{
// 		ID:             resource.ID,
// 		Title:          resource.Title,
// 		ProjectID:      resource.ProjectID,
// 		ResourceTypeID: resource.ResourceTypeID,
// 		URL:            resource.URL,
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// DeleteProjectResource godoc
// @Summary Delete a resource
// @Description Delete a resource and its file
// @Tags Resource
// @Param id path int true "Resource ID"
// @Success 200 {object} models.ResponseMessage
// @Failure 500 {object} models.ResponseMessage
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

	h.projectService.PublishProjectMessageToElasticSearch(c, "update", &detailedResource.Project)

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

	assetResource, err := h.resourceService.UploadAssetResource(c, &req.AssetResource, req.File, req.Title)
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

// GetResourcesByProjectID godoc
// @Summary Get resources by project ID
// @Description Get all resources associated with a project
// @Tags Resource
// @Produce json
// @Param project_id path int true "Project ID"
// @Success 200 {array} models.UploadResourceResponse
// @Failure 404 {object} map[string]string
// @Router /resource/project/{project_id} [get]
// func (h *resourceHandler) GetResourcesByProjectID(c *gin.Context) {
// 	projectID := c.Param("project_id")

// 	resources, err := h.resourceService.GetResourcesByProjectID(c, projectID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resources"})
// 		return
// 	}
// 	var response []models.UploadResourceResponse

// 	// Iterate over the resources and map them to UploadResourceResponse
// 	for _, resource := range resources {
// 		response = append(response, models.UploadResourceResponse{
// 			ID:             resource.ID,
// 			Title:          resource.Title, // Assuming resource.Title is a *string
// 			ProjectID:      resource.ProjectID,
// 			ResourceTypeID: resource.ResourceTypeID,
// 			URL:            resource.URL,
// 		})
// 	}

// 	c.JSON(http.StatusOK, response)
// }
