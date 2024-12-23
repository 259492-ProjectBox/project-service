package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupEmployeeRouter(r *gin.RouterGroup, handler handlers.EmployeeHandler) {
	employeeRouteV1 := r.Group("/employee")
	{
		employeeRouteV1.GET("/:id", handler.GetEmployeeByIDHandler)
		employeeRouteV1.GET("/GetByMajorID/:major_id", handler.GetEmployeeByMajorIDHandler)
		employeeRouteV1.POST("", handler.CreateEmployeeHandler)
		employeeRouteV1.PUT("", handler.UpdateEmployeeHandler)
	}
}
