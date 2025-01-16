package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupCalendarRouter(r *gin.RouterGroup, handler handlers.CalendarHandler) {
	calendarRouteV1 := r.Group("/v1/calendars")
	{

		calendarRouteV1.POST("/", handler.CreateCalendar)
		calendarRouteV1.GET("/program/:program_id", handler.GetCalendarByProgramId)
		calendarRouteV1.PUT("/", handler.UpdateCalendar)
		calendarRouteV1.DELETE("/:id", handler.DeleteCalendar)
	}
}
