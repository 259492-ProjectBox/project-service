package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupCalendarRouter(r *gin.RouterGroup, handler handlers.CalendarHandler) {
	calendarRouteV1 := r.Group("/calendar")
	{

		calendarRouteV1.POST("", handler.CreateCalendarHandler)
		calendarRouteV1.GET("/GetByMajorID/:major_id", handler.GetCalendarByMajorIDHandler)
		calendarRouteV1.PUT("", handler.UpdateCalendarHandler)
		calendarRouteV1.DELETE("/:id", handler.DeleteCalendarHandler)
	}
}
