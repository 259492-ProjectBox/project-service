package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupProjectConfigRouter(r *gin.RouterGroup, handler handlers.ProjectConfigHandler) {
	projectconfigRouteV1 := r.Group("/projectConfig")
	{

		// configRouteV1.POST("", handler.CreateCalendarHandler)
		projectconfigRouteV1.GET("/GetByMajorId/:major_id", handler.GetProjectConfigByMajorIDHandler)
		projectconfigRouteV1.POST("", handler.UpsertProjectConfigHandler)
		// configRouteV1.DELETE("/:id", handler.DeleteCalendarHandler)
	}
}
