package question

import "github.com/gin-gonic/gin"

func RegisterQuestionRoutes(questionRoutesGroup *gin.RouterGroup) {
	questionRoutesGroup.POST("/", newQuestion)
	questionRoutesGroup.GET("/:qId", getQuestions)
	questionRoutesGroup.PUT("/:qId", editQuestion)
	questionRoutesGroup.DELETE("/:qId", deleteQuestion)
}
