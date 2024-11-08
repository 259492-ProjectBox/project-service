// handlers/resource_handler.go
package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
// @Param title formData string false "Resource Title"
// @Success 201 {object} models.Resource
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

	title := c.PostForm("title")

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
		Title:          &title,
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

	c.JSON(http.StatusCreated, resource)
}

// GetResourceByID godoc
// @Summary Get a resource by ID
// @Description Get a resource by its ID
// @Tags Resource
// @Produce json
// @Param id path int true "Resource ID"
// @Success 200 {object} models.Resource
// @Failure 404 {object} map[string]string
// @Router /resource/{id} [get]
func (h *resourceHandler) GetResourceByID(c *gin.Context) {
	id := c.Param("id")

	resource, err := h.resourceService.GetResourceByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	c.JSON(http.StatusOK, resource)
}

// DeleteResource godoc
// @Summary Delete a resource
// @Description Delete a resource and its file
// @Tags Resource
// @Param id path int true "Resource ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
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
	filename := filepath.Base(resource.URL)

	// Delete from MinIO first
	err = h.minioClient.RemoveObject(context.Background(), h.bucketName, filename, minio.RemoveObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	// Then delete from database
	if err := h.resourceService.DeleteResource(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete resource record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resource deleted successfully"})
}

// GetResourcesByProjectID godoc
// @Summary Get resources by project ID
// @Description Get all resources associated with a project
// @Tags Resource
// @Produce json
// @Param project_id path int true "Project ID"
// @Success 200 {array} models.Resource
// @Failure 404 {object} map[string]string
// @Router /resource/project/{project_id} [get]
func (h *resourceHandler) GetResourcesByProjectID(c *gin.Context) {
	projectID := c.Param("project_id")

	resources, err := h.resourceService.GetResourcesByProjectID(c, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resources"})
		return
	}

	c.JSON(http.StatusOK, resources)
}
