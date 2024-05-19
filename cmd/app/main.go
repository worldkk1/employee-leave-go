package main

import (
	"github.com/gin-gonic/gin"
	"github.com/worldkk1/employee-leave-go/internal/app/database"
	healthCheck "github.com/worldkk1/employee-leave-go/internal/health-check"
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
	healthCheck.SetupRouter(router.Group("/"))

	router.Run("localhost:8080")
}
