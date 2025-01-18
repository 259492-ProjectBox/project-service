package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupProjectRouter(r *gin.RouterGroup, handler handlers.ProjectHandler) {
	projectRouteV1 := r.Group("/v1/projects")
	{
		projectRouteV1.GET("/:id", handler.GetProjectById)
		projectRouteV1.GET("/student/:student_id", handler.GetProjectsByStudentId)
		projectRouteV1.POST("/", handler.CreateProject)
		projectRouteV1.PUT("/", handler.UpdateProject)
		projectRouteV1.DELETE("/:id", handler.DeleteProject)
	}
}
