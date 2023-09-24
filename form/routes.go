package form

import "github.com/gin-gonic/gin"

func RegisterFormRoutes(formRoutesGroup *gin.RouterGroup) {
	formRoutesGroup.GET("/:id", getFormById)
	formRoutesGroup.POST("/", createForm)
	formRoutesGroup.PUT("/:id", editForm)
	formRoutesGroup.DELETE("/:id", deleteForm)

	formRoutesGroup.GET("/q/:id", getQuestions)
	formRoutesGroup.POST("/q/:id", postQuestions)
}
