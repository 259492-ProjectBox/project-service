package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupResourceRouter(r *gin.RouterGroup, handler handlers.ResourceHandler) {
	// resourceRouteV1 := r.Group("/v1/projects/resource")
	// {
	// resourceRouteV1.POST("/resource", handler.UploadResource)
	// resourceRouteV1.GET("/resource/:id", handler.GetResourceByID)

	// resourceRouteV1.GET("/resource/project/:project_id", handler.GetResourcesByProjectID)
	// }

	projectResourceRouteV1 := r.Group("/v1/projects/projectResource")
	{
		projectResourceRouteV1.DELETE("/:id", handler.DeleteProjectResource)
	}

}
