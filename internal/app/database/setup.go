package database

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/worldkk1/employee-leave-go/internal/app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDB() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("failed to load .env file")
	}

	port, err := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	if err != nil {
		panic("failed to get db config")
	}
	var (
		host     = os.Getenv("DATABASE_HOST")
		user     = os.Getenv("DATABASE_USERNAME")
		password = os.Getenv("DATABASE_PASSWORD")
		dbName   = os.Getenv("DATABASE_NAME")
	)

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(
		&models.User{},
		&models.LeaveType{},
		&models.UserLeave{},
		&models.UserLeaveRecord{},
	)

	DB = db
}
