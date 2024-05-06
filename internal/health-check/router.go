package healthcheck

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.RouterGroup) {
	router.GET("/health", HealthCheck)
}
