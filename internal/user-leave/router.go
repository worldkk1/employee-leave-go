package userleave

import "github.com/gin-gonic/gin"

func SetupRouter(router *gin.RouterGroup) {
	userLeaves := router.Group("/user-leaves")

	userLeaves.POST("/", RequestLeave)
	userLeaves.GET("/:leaveRecordId", GetLeaveDetail)
}
