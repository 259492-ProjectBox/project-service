package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupProjectRoleRouter(r *gin.RouterGroup, handler handlers.ProjectRoleHandler) {
	projectRoleRouteV1 := r.Group("/v1/projectRoles")
	{
		projectRoleRouteV1.GET("/program/:program_id", handler.GetByProgram)
	}
}
