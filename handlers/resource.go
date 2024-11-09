// handlers/resource_handler.go
package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/project-box/models"
	"github.com/project-box/services"
)

type ResourceHandler interface {
	UploadResource(c *gin.Context)
	GetResourceByID(c *gin.Context)
	DeleteResource(c *gin.Context)
	GetResourcesByProjectID(c *gin.Context)
}

type resourceHandler struct {
	minioClient     *minio.Client
	bucketName      string
	resourceService services.ResourceService
}

func NewResourceHandler(minioClient *minio.Client, resourceService services.ResourceService) ResourceHandler {
	return &resourceHandler{
		minioClient:     minioClient,
		bucketName:      os.Getenv("MINIO_BUCKET"),
		resourceService: resourceService,
	}
}

// UploadResource godoc
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
func (h *resourceHandler) UploadResource(c *gin.Context) {
	// Parse form data
	projectID, err := strconv.Atoi(c.PostForm("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project_id"})
		return
	}

	resourceTypeID, err := strconv.Atoi(c.PostForm("resource_type_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource_type_id"})
		return
	}

	// Get the file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	// Generate unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)

	// Upload to MinIO
	contentType := "application/pdf" // You might want to make this dynamic based on file type
	_, err = h.minioClient.PutObject(
		context.Background(),
		h.bucketName,
		filename,
		file,
		header.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	// Generate presigned URL
	url, err := h.minioClient.PresignedGetObject(context.Background(), h.bucketName, filename, time.Hour*24*7, nil)
	if err != nil {
		// Cleanup the uploaded file if URL generation fails
		h.minioClient.RemoveObject(context.Background(), h.bucketName, filename, minio.RemoveObjectOptions{})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate URL"})
		return
	}

	// Create resource in database
	resource, err := h.resourceService.CreateResource(c, &models.Resource{
		Title:          &filename,
		ProjectID:      projectID,
		ResourceTypeID: resourceTypeID,
		URL:            url.String(),
	})

	if err != nil {
		// Cleanup the uploaded file if database save fails
		h.minioClient.RemoveObject(context.Background(), h.bucketName, filename, minio.RemoveObjectOptions{})
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save resource: %v", err)})
		return
	}
	response := models.UploadResourceResponse{
		ID:             resource.ID,
		Title:          &filename,
		ProjectID:      projectID,
		ResourceTypeID: resourceTypeID,
		URL:            url.String(),
	}

	c.JSON(http.StatusCreated, response)
}

// GetResourceByID godoc
// @Summary Get a resource by ID
// @Description Get a resource by its ID
// @Tags Resource
// @Produce json
// @Param id path int true "Resource ID"
// @Success 200 {object} models.UploadResourceResponse
// @Failure 404 {object} map[string]string
// @Router /resource/{id} [get]
func (h *resourceHandler) GetResourceByID(c *gin.Context) {
	id := c.Param("id")

	resource, err := h.resourceService.GetResourceByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	// fmt.Println(*resource.Title)
	response := models.UploadResourceResponse{
		ID:             resource.ID,
		Title:          resource.Title,
		ProjectID:      resource.ProjectID,
		ResourceTypeID: resource.ResourceTypeID,
		URL:            resource.URL,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteResource godoc
// @Summary Delete a resource
// @Description Delete a resource and its file
// @Tags Resource
// @Param id path int true "Resource ID"
// @Success 200 {object} models.ResponseMessage
// @Failure 500 {object} models.ResponseMessage
// @Router /resource/{id} [delete]
func (h *resourceHandler) DeleteResource(c *gin.Context) {
	id := c.Param("id")

	// Get resource to get the file URL
	resource, err := h.resourceService.GetResourceByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	// Extract filename from URL
	// filename := filepath.Base(resource.URL)
	filename := *resource.Title

	// Delete from MinIO first
	err = h.minioClient.RemoveObject(context.Background(), h.bucketName, filename, minio.RemoveObjectOptions{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseMessage{Message: "Failed to delete file"})
		return
	}

	// Then delete from database
	if err := h.resourceService.DeleteResource(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseMessage{Message: "Failed to delete resource record"})
		return
	}

	c.JSON(http.StatusOK, models.ResponseMessage{Message: "Resource deleted successfully"})

	c.JSON(http.StatusOK, gin.H{"message": "Resource deleted successfully"})
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
func (h *resourceHandler) GetResourcesByProjectID(c *gin.Context) {
	projectID := c.Param("project_id")

	resources, err := h.resourceService.GetResourcesByProjectID(c, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resources"})
		return
	}
	var response []models.UploadResourceResponse

	// Iterate over the resources and map them to UploadResourceResponse
	for _, resource := range resources {
		response = append(response, models.UploadResourceResponse{
			ID:             resource.ID,
			Title:          resource.Title, // Assuming resource.Title is a *string
			ProjectID:      resource.ProjectID,
			ResourceTypeID: resource.ResourceTypeID,
			URL:            resource.URL,
		})
	}

	c.JSON(http.StatusOK, response)
}
