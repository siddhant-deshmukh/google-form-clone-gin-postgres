package responses

import "github.com/gin-gonic/gin"

func RegisterResponseRoutes(resRouterGroup *gin.RouterGroup) {
	resRouterGroup.GET("/:fId", getResponses)
	resRouterGroup.POST("/:fId", submitResponse)
}
