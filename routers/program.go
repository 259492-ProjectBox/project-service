package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetUpProgramRoute(r *gin.RouterGroup, handler handlers.ProgramHandler) {
	programRouteV1 := r.Group("/v1/programs")
	{

		programRouteV1.POST("/", handler.CreateProgram)
		programRouteV1.GET("/", handler.GetPrograms)
		programRouteV1.PUT("/update-name", handler.UpdateProgramName)
	}
}
