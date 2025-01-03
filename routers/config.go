package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupConfigRouter(r *gin.RouterGroup, handler handlers.ConfigHandler) {
	configRouteV1 := r.Group("/config")
	{

		// configRouteV1.POST("", handler.CreateCalendarHandler)
		configRouteV1.GET("/GetByMajorId/:major_id", handler.GetConfigByMajorIDHandler)
		// configRouteV1.PUT("", handler.UpdateCalendarHandler)
		// configRouteV1.DELETE("/:id", handler.DeleteCalendarHandler)
	}
}
