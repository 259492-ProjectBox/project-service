package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupEmployeeRouter(r *gin.RouterGroup, handler handlers.EmployeeHandler) {
	employeeRouteV1 := r.Group("/employee")
	{
		employeeRouteV1.GET("/:id", handler.GetEmployeeById)
		employeeRouteV1.GET("/program/:program_id", handler.GetEmployeeByProgramId)
		employeeRouteV1.POST("", handler.CreateEmployee)
		employeeRouteV1.PUT("", handler.UpdateEmployee)
	}
}
