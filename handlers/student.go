package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project-box/services"
	"gorm.io/gorm"
)

type StudentHandler interface {
	GetsStudentByProgramIdOnCurrentYearAndSemester(c *gin.Context)
	GetStudentByStudentId(c *gin.Context)
	CheckStudentPermissionForCreateProject(c *gin.Context)
}

type studentHandler struct {
	studentService services.StudentService
}

func NewStudentHandler(studentService services.StudentService) StudentHandler {
	return &studentHandler{
		studentService: studentService,
	}
}

// @Summary Get students by student ID
// @Description Retrieves a list of students for a given student ID
// @Tags Student
// @Produce json
// @Param student_id path int true "Student ID"
// @Success 200 {array} models.Student "Successfully retrieved students"
// @Failure 400 {object} map[string]interface{} "Invalid program ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/students/{student_id}/program/{program_id} [get]
func (h *studentHandler) GetStudentByStudentId(c *gin.Context) {
	studentId := c.Param("student_id")
	programIdStr := c.Param("program_id")
	programId, err := strconv.Atoi(programIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program ID"})
		return
	}
	students, err := h.studentService.GetStudentByStudentIdAndProgramIdOnCurrentYearAndSemester(c.Request.Context(), studentId, programId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}

// @Summary Check student permission for creating a project
// @Description Checks if a student has permission to create a project based on their student ID
// @Tags Student
// @Produce json
// @Param student_id path int true "Student ID"
// @Success 200 {object} map[string]interface{} "Successfully checked permission"
// @Failure 400 {object} map[string]interface{} "Invalid student ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/students/{student_id}/program/{program_id}/check [get]
func (h *studentHandler) CheckStudentPermissionForCreateProject(c *gin.Context) {
	studentId := c.Param("student_id")
	programIdStr := c.Param("program_id")
	programId, err := strconv.Atoi(programIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program ID"})
		return
	}
	student, err := h.studentService.GetStudentByStudentIdAndProgramIdOnCurrentYearAndSemester(c.Request.Context(), studentId, programId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"has_permission": false})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	var hasPermission bool
	if student != nil {
		hasPermission = true
	}

	c.JSON(http.StatusOK, gin.H{"has_permission": hasPermission})
}

// @Summary Get students by program ID, academic year, and semester
// @Description Retrieves a list of students for a given program ID, academic year, and semester
// @Tags Student
// @Produce json
// @Param program_id path int true "Program ID"
// @Success 200 {array} models.Student "Successfully retrieved students"
// @Failure 400 {object} map[string]interface{} "Invalid parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/students/program/{program_id}/current_year [get]
func (h *studentHandler) GetsStudentByProgramIdOnCurrentYearAndSemester(c *gin.Context) {
	programIdStr := c.Param("program_id")
	programId, err := strconv.Atoi(programIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program ID"})
		return
	}

	students, err := h.studentService.GetStudentByProgramIdOnCurrentYearAndSemester(c.Request.Context(), programId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}
