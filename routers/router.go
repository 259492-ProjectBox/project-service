package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/project-box/docs"
	"github.com/project-box/handlers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine, projectHandler handlers.ProjectHandler, resourceHandler handlers.ResourceHandler, calendarHandler handlers.CalendarHandler, employeeHandler handlers.EmployeeHandler,
	configHandler handlers.ConfigHandler, projectconfigHandler handlers.ProjectConfigHandler, programHandler handlers.ProgramHandler) {
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the api",
		})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Dynamically set Swagger URL based on the request
	// NOT WORKING YET when build go to main and change the host instead for now
	// r.GET("/swagger/*any", func(c *gin.Context) {
	// 	scheme := "http"
	// 	if c.Request.TLS != nil {
	// 		scheme = "https"
	// 	}
	// 	host := c.Request.Host
	// 	swaggerURL := fmt.Sprintf("%s://%s/swagger/doc.json", scheme, host)
	// 	ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL(swaggerURL))(c)
	// })
	router := r.Group("")
	SetupProjectRouter(router, projectHandler)
	SetupResourceRouter(router, resourceHandler)
	SetupCalendarRouter(router, calendarHandler)
	SetupEmployeeRouter(router, employeeHandler)
	SetupConfigRouter(router, configHandler)
	SetupProjectConfigRouter(router, projectconfigHandler)
	SetUpProgramRoute(router, programHandler)
}
