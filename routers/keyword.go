package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/project-box/handlers"
)

func SetupKeywordRouter(r *gin.RouterGroup, handler handlers.KeywordHandler) {
	keywordRouteV1 := r.Group("/v1/keywords")
	{
		keywordRouteV1.GET("/all", handler.GetAllKeywords)
		keywordRouteV1.GET("", handler.GetKeywords)
		keywordRouteV1.GET("/:id", handler.GetKeyword)
		keywordRouteV1.POST("", handler.CreateKeyword)
		keywordRouteV1.PUT("", handler.UpdateKeyword)
		keywordRouteV1.DELETE("/:id", handler.DeleteKeyword)
	}
}
