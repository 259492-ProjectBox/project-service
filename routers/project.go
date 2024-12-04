package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupProjectRouter(r *gin.RouterGroup, handler handlers.ProjectHandler) {
	projectRoute := r.Group("/v1/projects")
	{
		projectRoute.GET("/:id", handler.GetProjectById)
		projectRoute.GET("/student/:student_id", handler.GetProjectsByStudentId)
		projectRoute.POST("", handler.CreateProject)
		projectRoute.PUT("/:id", handler.UpdateProject)
		projectRoute.DELETE("/:id", handler.DeleteProject)
	}
}
