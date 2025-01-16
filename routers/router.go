package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/project-box/docs"
	"github.com/project-box/handlers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine, projectHandler handlers.ProjectHandler, resourceHandler handlers.ResourceHandler, calendarHandler handlers.CalendarHandler, staffHandler handlers.StaffHandler,
	configHandler handlers.ConfigHandler, projectconfigHandler handlers.ProjectConfigHandler, programHandler handlers.ProgramHandler) {
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the api",
		})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router := r.Group("api")
	SetupProjectRouter(router, projectHandler)
	SetupResourceRouter(router, resourceHandler)
	SetupCalendarRouter(router, calendarHandler)
	SetupStaffRouter(router, staffHandler)
	SetupConfigRouter(router, configHandler)
	SetupProjectConfigRouter(router, projectconfigHandler)
	SetUpProgramRoute(router, programHandler)
}
