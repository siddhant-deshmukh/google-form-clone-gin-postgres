package user

import "github.com/gin-gonic/gin"

func RegisterUserAuthRoutes(userRoutesGroup *gin.RouterGroup) {
	userRoutesGroup.POST("/", registerUser)
	userRoutesGroup.POST("/login", userLogin)
}
