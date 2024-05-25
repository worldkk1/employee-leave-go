package main

import (
	"github.com/gin-gonic/gin"
	"github.com/worldkk1/employee-leave-go/internal/app/database"
	healthCheck "github.com/worldkk1/employee-leave-go/internal/health-check"
	"github.com/worldkk1/employee-leave-go/internal/user"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	database.SetupDB()

	router := gin.Default()
	apiRouter := router.Group("/")
	healthCheck.SetupRouter(apiRouter)
	user.SetupRouter(apiRouter)

	router.Run("localhost:8080")
}
