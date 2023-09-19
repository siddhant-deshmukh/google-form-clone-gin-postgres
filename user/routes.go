package user

import "github.com/gin-gonic/gin"

func RegisterUserRoutes(userRoutesGroup *gin.RouterGroup) {
	userRoutesGroup.GET("/", getData)
	// userRoutesGroup.POST("/login", userLogin)
}
