package user

import "github.com/gin-gonic/gin"

func SetupRouter(router *gin.RouterGroup) {
	users := router.Group("/users")

	users.POST("", CreateUser)
}
