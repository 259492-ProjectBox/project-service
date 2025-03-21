package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/project-box/docs"
	"github.com/project-box/handlers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine, projectHandler handlers.ProjectHandler, resourceHandler handlers.ResourceHandler, staffHandler handlers.StaffHandler,
	configHandler handlers.ConfigHandler, projectConfigHandler handlers.ProjectConfigHandler, projectResourceConfigHandler handlers.ProjectResourceConfigHandler, projectRoleHandler handlers.ProjectRoleHandler, programHandler handlers.ProgramHandler, studentHandler handlers.StudentHandler, uploadHandler handlers.UploadHandler, keywordHandler handlers.KeywordHandler) {
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the api",
		})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router := r.Group("/api")
	SetupKeywordRouter(router, keywordHandler)
	SetupProjectRouter(router, projectHandler)
	SetupResourceRouter(router, resourceHandler)
	SetupStaffRouter(router, staffHandler)
	SetupConfigRouter(router, configHandler)
	SetupProjectConfigRouter(router, projectConfigHandler)
	SetupProjectResourceConfigRouter(router, projectResourceConfigHandler)
	SetupProjectRoleRouter(router, projectRoleHandler)
	SetUpProgramRoute(router, programHandler)
	SetupStudentRouter(router, studentHandler)
	SetupUploadRouter(router, uploadHandler)
}
