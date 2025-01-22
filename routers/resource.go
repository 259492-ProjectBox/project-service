package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupResourceRouter(r *gin.RouterGroup, handler handlers.ResourceHandler) {
	assetResourceRouteV1 := r.Group("/v1/assetResources")
	{
		assetResourceRouteV1.POST("", handler.UploadAssetResource)
		assetResourceRouteV1.GET("/program/:program_id", handler.GetAssetResourceByProgramID)
	}

	projectResourceRouteV1 := r.Group("/v1/projectResources")
	{
		projectResourceRouteV1.DELETE("/:id", handler.DeleteProjectResource)
	}

}
