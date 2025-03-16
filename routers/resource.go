package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupResourceRouter(r *gin.RouterGroup, handler handlers.ResourceHandler) {
	projectResourceRouteV1 := r.Group("/v1/projectResources")
	{
		projectResourceRouteV1.DELETE("/:id", handler.DeleteProjectResource)
	}

}
