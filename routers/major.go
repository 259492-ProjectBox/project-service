package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupMajorRouter(r *gin.RouterGroup, handler handlers.MajorHandler) {
	majorRouteV1 := r.Group("/major")
	{

		majorRouteV1.POST("", handler.CreateMajorHandler)
		majorRouteV1.GET("/GetAllMajor", handler.GetAllMajorHandler)
		majorRouteV1.PUT("/UpdateMajorName", handler.UpdateMajorNameHandler)
		// majorRouteV1.DELETE("/:id", handler.DeleteCalendarHandler)
	}
}
