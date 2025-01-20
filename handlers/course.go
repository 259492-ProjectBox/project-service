package handlers

import (
	"github.com/project-box/services"
)

type CourseHandler interface {
}

type courseHandler struct {
	courseService services.CourseService
}

func NewCourseHandler(courseService services.CourseService) CourseHandler {
	return &courseHandler{
		courseService: courseService,
	}
}
