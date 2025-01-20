package handlers

import (
	"github.com/project-box/services"
)

type StudentHandler interface {
}

type studentHandler struct {
	studentService services.StudentService
}

func NewStudentHandler(studentService services.StudentService) StudentHandler {
	return &studentHandler{
		studentService: studentService,
	}
}
