package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupStaffRouter(r *gin.RouterGroup, handler handlers.StaffHandler) {
	staffRouteV1 := r.Group("/staff")
	{
		staffRouteV1.GET("/:id", handler.GetStaffById)
		staffRouteV1.GET("/program/:program_id", handler.GetStaffByProgramId)
		staffRouteV1.POST("/", handler.CreateStaff)
		staffRouteV1.PUT("", handler.UpdateStaff)
	}
}
