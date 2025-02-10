package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupProjectResourceConfigRouter(r *gin.RouterGroup, handler handlers.ProjectResourceConfigHandler) {
	projectResourceConfigRouteV1 := r.Group("/v1/projectResourceConfigs")
	{
		projectResourceConfigRouteV1.GET("/program/:program_id", handler.GetProjectResourceConfigsByProgramId)
		projectResourceConfigRouteV1.PUT("", handler.UpsertProjectResourceConfig)
	}

	projectResourceConfigRouteV2 := r.Group("/v2/projectResourceConfigs")
	{
		projectResourceConfigRouteV2.PUT("", handler.UpsertProjectResourceConfigV2)
	}

}
