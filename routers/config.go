package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupConfigRouter(r *gin.RouterGroup, handler handlers.ConfigHandler) {
	configRouteV1 := r.Group("/v1/configs")
	{
		configRouteV1.GET("/program/:program_id", handler.GetConfigByProgramId)
		configRouteV1.PUT("", handler.UpsertConfig)
	}
}
