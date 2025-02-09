package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project-box/services"
)

type UploadHandler interface {
	UploadStudentEnrollmentFile(c *gin.Context)
	UploadCreateProjectFile(c *gin.Context)
}

type uploadHandler struct {
	uploadService services.UploadService
}

func NewUploadHandler(uploadService services.UploadService) UploadHandler {
	return &uploadHandler{
		uploadService: uploadService,
	}
}

// @Summary Upload student enrollment file
// @Description Uploads and processes a student enrollment file for a given program ID
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param program_id path int true "Program ID"
// @Param file formData file true "Student Enrollment File"
// @Success 200 {object} map[string]interface{} "file processed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid program ID or failed to retrieve the file"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/uploads/program/{program_id}/student-enrollment [post]
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

	err = h.uploadService.ProcessStudentEnrollmentFile(c.Request.Context(), programId, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file processed successfully"})
}

// @Summary Upload create project file
// @Description Uploads and processes a create project file for a given program ID
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param program_id path int true "Program ID"
// @Param file formData file true "Create Project File"
// @Success 200 {object} map[string]interface{} "file processed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid program ID or failed to retrieve the file"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/uploads/program/{program_id}/create-project [post]
func (h *uploadHandler) UploadCreateProjectFile(c *gin.Context) {
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

	err = h.uploadService.ProcessCreateProjectFile(c.Request.Context(), programId, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file processed successfully"})
}
