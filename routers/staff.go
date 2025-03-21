package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupStaffRouter(r *gin.RouterGroup, handler handlers.StaffHandler) {
	staffRouteV1 := r.Group("/v1/staffs")
	{
		staffRouteV1.GET("/:id", handler.GetStaffById)
		staffRouteV1.GET("/program/:program_id", handler.GetStaffByProgramId)
		staffRouteV1.GET("/email/:email", handler.GetStaffByEmail)
		staffRouteV1.POST("", handler.CreateStaff)
		staffRouteV1.PUT("", handler.UpdateStaff)
		staffRouteV1.GET("/GetAllStaffs", handler.GetAllStaff)
	}

	staffRouteV2 := r.Group("/v2/staffs")
	{
		staffRouteV2.GET("/program/:program_id", handler.GetStaffByProgramIdV2)
	}
}
