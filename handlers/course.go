package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project-box/models"
	"github.com/project-box/services"
)

type CourseHandler interface {
	GetCourseByProgramID(c *gin.Context)
	GetCourseByCourseNo(c *gin.Context)
	CreateCourse(c *gin.Context)
	UpdateCourse(c *gin.Context)
	DeleteCourse(c *gin.Context)
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

// @Summary Create a new course
// @Description Creates a new course
// @Tags Course
// @Accept json
// @Produce json
// @Param course body models.Course true "Course data"
// @Success 201 {object} models.Course "Successfully created course"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/courses [post]
func (h *courseHandler) CreateCourse(c *gin.Context) {
	var req models.Course
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	course, err := h.courseService.CreateCourse(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, course)
}

// @Summary Update an existing course
// @Description Updates an existing course
// @Tags Course
// @Accept json
// @Produce json
// @Param course body models.Course true "Course data"
// @Success 200 {object} models.Course "Successfully updated course"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/courses [put]
func (h *courseHandler) UpdateCourse(c *gin.Context) {
	var req models.Course
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	course, err := h.courseService.UpdateCourse(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

// @Summary Delete a course
// @Description Deletes a course by its ID
// @Tags Course
// @Produce json
// @Param course_id path int true "Course ID"
// @Success 204 "Successfully deleted course"
// @Failure 400 {object} map[string]interface{} "Invalid course ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /v1/courses/{course_id} [delete]
func (h *courseHandler) DeleteCourse(c *gin.Context) {
	courseIDStr := c.Param("course_id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	if err := h.courseService.DeleteCourse(c.Request.Context(), courseID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
