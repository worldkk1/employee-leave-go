package database

import (
	"github.com/worldkk1/employee-leave-go/internal/app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=employee_leave port=5432"
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

	return db
}
