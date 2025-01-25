package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project-box/services"
)

type UploadHandler interface {
	UploadStudentEnrollmentFile(c *gin.Context)
}

type uploadHandler struct {
	uploadService services.UploadService
}

func NewUploadHandler(uploadService services.UploadService) UploadHandler {
	return &uploadHandler{
		uploadService: uploadService,
	}
}

func (h *uploadHandler) UploadStudentEnrollmentFile(c *gin.Context) {
	programId, err := strconv.Atoi(c.Param("program_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program ID"})
		return
	}
	
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to retrieve the file"})
		return
	}

	err = h.uploadService.ProcessStudentEnrollmentFile(c, programId, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file processed successfully"})
}
