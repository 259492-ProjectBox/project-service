package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupConfigRouter(r *gin.RouterGroup, handler handlers.ConfigHandler) {
	configRouteV1 := r.Group("/config")
	{

		configRouteV1.GET("/program/:program_id", handler.GetConfigByProgramId)
	}
}
