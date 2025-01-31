package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupStudentRouter(r *gin.RouterGroup, handler handlers.StudentHandler) {
	studentRouteV1 := r.Group("/v1/students")
	{
		studentRouteV1.GET("/:student_id", handler.GetStudentByStudentId)
		studentRouteV1.GET("/:student_id/check", handler.CheckStudentPermissionForCreateProject)
		studentRouteV1.GET("/program/:program_id/current_year", handler.GetsStudentByProgramIdOnCurrentYearAndSemester)
	}
}
