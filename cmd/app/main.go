package main

import (
	"github.com/gin-gonic/gin"
	healthcheck "github.com/worldkk1/employee-leave-go/internal/health-check"
)

func main() {
	router := gin.Default()
	healthcheck.SetupRouter(router.Group("/"))

	router.Run("localhost:8080")
}
