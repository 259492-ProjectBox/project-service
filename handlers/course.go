package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project-box/services"
)

type CourseHandler interface {
	GetCourseByProgramID(c *gin.Context)
	GetCourseByCourseNo(c *gin.Context)
}

type courseHandler struct {
	courseService services.CourseService
}

func NewCourseHandler(courseService services.CourseService) CourseHandler {
	return &courseHandler{
		courseService: courseService,
	}
}

// @Summary Get courses by program ID
// @Description Retrieves a list of courses for a given program ID
// @Tags Course
// @Produce json
// @Param program_id path int true "Program ID"
// @Success 200 {array} models.Course "Successfully retrieved courses"
// @Failure 400 {object} map[string]interface{} "Invalid program ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/courses/program/{program_id} [get]
func (h *courseHandler) GetCourseByProgramID(c *gin.Context) {
	programIDStr := c.Param("program_id")
	programID, err := strconv.Atoi(programIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program ID"})
		return
	}

	courses, err := h.courseService.FindCourseByProgramID(c.Request.Context(), programID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, courses)
}

// @Summary Get course by course number
// @Description Retrieves a course by its course number
// @Tags Course
// @Produce json
// @Param course_no path string true "Course Number"
// @Success 200 {object} models.Course "Successfully retrieved course"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/courses/{course_no} [get]
func (h *courseHandler) GetCourseByCourseNo(c *gin.Context) {
	courseNo := c.Param("course_no")

	course, err := h.courseService.FindCourseByCourseNo(c.Request.Context(), courseNo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}
