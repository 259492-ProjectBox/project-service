package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupProjectConfigRouter(r *gin.RouterGroup, handler handlers.ProjectConfigHandler) {
	projectconfigRouteV1 := r.Group("/projectConfigs")
	{

		projectconfigRouteV1.GET("/program/:program_id", handler.GetProjectConfigByProgramId)
		projectconfigRouteV1.POST("/", handler.UpsertProjectConfig)
	}
}
