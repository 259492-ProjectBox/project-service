package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupCalendarRouter(r *gin.RouterGroup, handler handlers.CalendarHandler) {
	calendarRouteV1 := r.Group("/calendar")
	{

		calendarRouteV1.POST("", handler.CreateCalendar)

	}
}
