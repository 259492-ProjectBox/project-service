package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupProjectRouter(r *gin.RouterGroup, handler handlers.ProjectHandler) {
	r.GET("/projects/:id", handler.GetProjectById)
	r.GET("/projects/student/:student_id", handler.GetProjectsByStudentId)
	r.POST("/projects", handler.CreateProject)
	r.PUT("/projects/:id", handler.UpdateProject)
	r.DELETE("/projects/:id", handler.DeleteProject)
}
