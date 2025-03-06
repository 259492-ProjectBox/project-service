package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupCourseRouter(r *gin.RouterGroup, handler handlers.CourseHandler) {
	courseRouteV1 := r.Group("/v1/courses")
	{
		courseRouteV1.GET("/program/:program_id", handler.GetCourseByProgramID)
		courseRouteV1.GET("/:course_no", handler.GetCourseByCourseNo)
		courseRouteV1.POST("", handler.CreateCourse)
		courseRouteV1.PUT("", handler.UpdateCourse)
		courseRouteV1.DELETE("/:course_id", handler.DeleteCourse)
	}
}
